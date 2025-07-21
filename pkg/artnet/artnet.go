package artnet

import (
	"encoding/binary"
	"errors"
	"github.com/techwolf12/artnet-to-hue/pkg/config"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/ipv4"
)

const (
	artnetPort      = 6454
	artnetHeader    = "Art-Net\x00"
	maxLights       = 10
	dmxPacketLength = 512
	opDmx           = 0x5000
	opPoll          = 0x2000
	opPollReply     = 0x2100
)

type Listener struct {
	conn         *net.UDPConn
	universe     uint16
	startAddress int
	numLights    int
	cb           func([]byte)
	mu           sync.Mutex
	config       config.Config
}

func (l *Listener) advertiseNode(shortName, longName string) {
	buf := make([]byte, 1024)
	for {
		n, addr, err := l.conn.ReadFromUDP(buf)
		if err != nil || n < 10 {
			continue
		}
		if !strings.HasPrefix(string(buf[:8]), artnetHeader) {
			continue
		}
		op := binary.LittleEndian.Uint16(buf[8:10])
		if op == opPoll {
			if l.config.Debug {
				log.Printf("Received ArtPoll from %s", addr)
			}
			reply := buildArtPollReply(l.conn.LocalAddr().(*net.UDPAddr), shortName, longName, l.universe)
			_, err := l.conn.WriteToUDP(reply, addr)
			if err != nil {
				log.Printf("Error writing to UDP: %v", err)
				return
			}
		}
	}
}

func buildArtPollReply(localAddr *net.UDPAddr, shortName, longName string, universe uint16) []byte {
	b := make([]byte, 239)
	copy(b[0:], artnetHeader)
	binary.LittleEndian.PutUint16(b[8:], opPollReply)
	ip := localAddr.IP.To4()
	if ip == nil {
		ip = net.IPv4(127, 0, 0, 1)
	}
	copy(b[10:], ip)
	binary.BigEndian.PutUint16(b[14:], artnetPort)
	b[16] = 0x00 // Version info
	b[17] = 0x01
	b[18] = 0x00            // NetSwitch
	b[19] = 0x00            // SubSwitch
	b[20] = 0x01            // OemHi
	b[21] = 0x23            // OemLo
	b[22] = 0x00            // Ubea version
	b[23] = 0x00            // Status1
	b[24] = 0x00            // EstaManLo
	b[25] = 0x00            // EstaManHi
	copy(b[26:], shortName) // Short name (max 18 bytes)
	copy(b[44:], longName)  // Long name (max 64 bytes)
	copy(b[108:], "#0001 [OK]")
	b[172] = 0x01 // Num ports
	b[190] = byte(universe & 0xFF)
	b[191] = byte((universe >> 8) & 0xFF)
	b[173] = 0x00
	return b
}

func multicastAddrForUniverse(universe uint16) string {
	hi := (universe >> 8) & 0xff
	lo := universe & 0xff
	return net.JoinHostPort(
		net.IPv4(239, 255, byte(hi), byte(lo)).String(),
		strconv.Itoa(artnetPort),
	)
}

func NewListener(config config.Config) (*Listener, error) {
	if config.NumLights < 1 || config.NumLights > maxLights {
		return nil, errors.New("numLights must be between 1 and 10")
	}
	if config.ArtNetStartAddress < 1 {
		return nil, errors.New("startAddress out of DMX range, must be between 1 and 512")
	}
	if config.ArtNetStartAddress+(4*config.NumLights)-1 > dmxPacketLength {
		return nil, errors.New("exceeding DMX packet length, (startAddress + 3 * lights) must be <= 512")
	}

	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.Itoa(artnetPort))
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	maddr, _ := net.ResolveUDPAddr("udp", multicastAddrForUniverse(config.ArtNetUniverse))
	p := ipv4.NewPacketConn(conn)
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		// Only join on interfaces that are up and support multicast
		if (iface.Flags&net.FlagUp) != 0 && (iface.Flags&net.FlagMulticast) != 0 {
			_ = p.JoinGroup(&iface, maddr)
		}
	}

	l := &Listener{
		conn:         conn,
		universe:     config.ArtNetUniverse,
		startAddress: config.ArtNetStartAddress,
		numLights:    config.NumLights * 3, // Each light uses 3 channels (RGB)
		config:       config,
	}
	go l.listen()
	go l.advertiseNode("artnet-to-hue", "Artnet to Hue Bridge")
	return l, nil
}

func (l *Listener) OnUpdate(cb func([]byte)) {
	l.mu.Lock()
	l.cb = cb
	l.mu.Unlock()
}

func (l *Listener) listen() {
	buf := make([]byte, 1024)
	for {
		n, _, err := l.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		if n < 18 { // Minimum ArtDMX packet size
			continue
		}
		if !strings.HasPrefix(string(buf[:8]), artnetHeader) {
			continue
		}
		op := binary.LittleEndian.Uint16(buf[8:10])
		if op != opDmx {
			continue
		}
		universe := binary.LittleEndian.Uint16(buf[14:16])
		if universe != l.universe {
			continue
		}
		length := int(binary.BigEndian.Uint16(buf[16:18]))
		if length < l.startAddress-1+l.numLights {
			continue
		}
		dmx := buf[18 : 18+length]
		start := l.startAddress - 1
		end := start + l.numLights
		if end > len(dmx) {
			continue
		}
		values := make([]byte, l.numLights)
		copy(values, dmx[start:end])
		l.mu.Lock()
		if l.cb != nil {
			go l.cb(values)
		}
		l.mu.Unlock()
	}
}
