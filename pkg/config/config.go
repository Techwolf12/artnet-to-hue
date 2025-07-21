package config

import "net"

type Config struct {
	HueBridgeIP        net.IP
	Username           string
	EntertainmentZone  int
	NumLights          int
	ArtNetUniverse     uint16
	ArtNetStartAddress int
	Debug              bool
}
