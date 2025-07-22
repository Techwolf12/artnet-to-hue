package hue

import (
	"bytes"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/pion/dtls/v2"
	"github.com/techwolf12/artnet-to-hue/pkg/config"
	"net"
	"net/http"
	"time"
)

type EntertainmentConfigResponse struct {
	Errors []interface{} `json:"errors"`
	Data   []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"data"`
}

type EntertainmentLightState struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

type Streamer struct {
	conn *dtls.Conn
}

func StartEntertainmentArea(config config.Config) error {
	url := fmt.Sprintf("https://%s/clip/v2/resource/entertainment_configuration/%s", config.HueBridgeIP, config.EntertainmentZone)
	body := []byte(`{"action":"start"}`)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("hue-application-key", config.Username)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start entertainment area: %s", resp.Status)
	}
	return nil
}

func BuildHueStreamPacket(configID string, states []EntertainmentLightState) []byte {
	protocolName := []byte("HueStream")
	header := []byte{0x02, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}
	configIDBuf := make([]byte, 36)
	copy(configIDBuf, configID)
	var channels []byte
	for i, state := range states {
		ch := []byte{byte(i)}
		ch = append(ch, byte(state.Red), byte(state.Red), byte(state.Green), byte(state.Green), byte(state.Blue), byte(state.Blue))
		channels = append(channels, ch...)
	}
	return bytes.Join([][]byte{protocolName, header, configIDBuf, channels}, nil)
}

func (hs *Streamer) Connect(config config.Config, hueAppId string) error {
	if hs.conn != nil {
		return nil // Already connected
	}
	clientKey, err := hex.DecodeString(config.ClientKey)
	if err != nil {
		return fmt.Errorf("failed to decode client key: %v", err)
	}
	cfg := &dtls.Config{
		PSK: func(hint []byte) ([]byte, error) {
			return clientKey, nil
		},
		PSKIdentityHint:    []byte(hueAppId),
		CipherSuites:       []dtls.CipherSuiteID{dtls.TLS_PSK_WITH_AES_128_GCM_SHA256},
		InsecureSkipVerify: true,
		FlightInterval:     500 * time.Millisecond,
	}
	conn, err := dtls.Dial("udp", &net.UDPAddr{IP: config.HueBridgeIP, Port: 2100}, cfg)
	if err != nil {
		return fmt.Errorf("DTLS dial failed: %w", err)
	}
	hs.conn = conn
	return nil
}

func (hs *Streamer) StreamToHue(config config.Config, states []EntertainmentLightState) error {
	if hs.conn == nil {
		return fmt.Errorf("DTLS connection not established")
	}
	packet := BuildHueStreamPacket(config.EntertainmentZone, states)
	_, err := hs.conn.Write(packet)
	return err
}
