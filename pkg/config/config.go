package config

import "net"

type Config struct {
	HueBridgeIP        net.IP
	Username           string
	EntertainmentZone  int
	ArtNetUniverse     int
	ArtNetStartAddress int
}
