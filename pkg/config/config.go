package config

import "net"

type Config struct {
	HueBridgeIP        net.IP
	Username           string
	ClientKey          string
	EntertainmentZone  string
	NumLights          int
	ArtNetUniverse     uint16
	ArtNetStartAddress int
	Debug              bool
}
