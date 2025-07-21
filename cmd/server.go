/*
Copyright © 2025 Christiaan de Die le Clercq <contact@techwolf12.nl>

*/
package cmd

import (
	"fmt"
	"github.com/techwolf12/artnet-to-hue/pkg/artnet"
	artnetHueConfig "github.com/techwolf12/artnet-to-hue/pkg/config"
	"github.com/techwolf12/artnet-to-hue/pkg/hue"
	"log"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Art-Net to Hue bridge server",
	Long:  `Start the Art-Net to Hue bridge server to listen for Art-Net packets and control Philips Hue lights.`,
	Run:   serverRun,
}

func serverRun(cmd *cobra.Command, args []string) {
	hueBridgeIP, _ := cmd.Flags().GetIP("hue-bridge-ip")
	if hueBridgeIP.IsUnspecified() {
		fmt.Println("Error: Hue bridge IP address is required")
		return
	}
	username, _ := cmd.Flags().GetString("username")
	if username == "" {
		fmt.Println("Error: Username for the Hue bridge is required")
		return
	}
	entertainmentZone, _ := cmd.Flags().GetInt("entertainment-zone")
	if entertainmentZone < 0 {
		fmt.Println("Error: Entertainment zone ID must be a non-negative integer")
		return
	}
	numLights, _ := cmd.Flags().GetInt("lights")
	if numLights <= 0 {
		fmt.Println("Error: Number of lights in the entertainment zone must be a positive integer")
		return
	}
	if numLights > 10 {
		fmt.Println("Error: Number of lights in the entertainment zone cannot exceed 10")
		return
	}
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		fmt.Println("Debug mode is enabled")
	}
	artnetUniverse, _ := cmd.Flags().GetUint16("artnet-universe")
	if artnetUniverse < 0 {
		fmt.Println("Error: Art-Net universe must be a non-negative integer")
		return
	}
	artnetDMXStart, _ := cmd.Flags().GetInt("artnet-dmx-start")
	if artnetDMXStart < 0 {
		fmt.Println("Error: Art-Net DMX start channel must be a non-negative integer")
		return
	}
	config := artnetHueConfig.Config{
		HueBridgeIP:        hueBridgeIP,
		Username:           username,
		EntertainmentZone:  entertainmentZone,
		NumLights:          numLights,
		ArtNetUniverse:     artnetUniverse,
		ArtNetStartAddress: artnetDMXStart,
		Debug:              debug,
	}
	fmt.Printf("Starting server with Hue Bridge IP: %s, Username: %s, Entertainment Zone: %d, Art-Net Universe: %d, Art-Net DMX Start: %d\n",
		hueBridgeIP, username, entertainmentZone, artnetUniverse, artnetDMXStart)

	listener, err := artnet.NewListener(config)
	if err != nil {
		panic(err)
	}

	listener.OnUpdate(func(values []byte) {
		if config.Debug {
			log.Printf("DMX values: %v\n", values)
		}
		if len(values) < config.NumLights*3 {
			log.Printf("Received fewer values than expected for %d lights: %d values received", config.NumLights, len(values))
			return
		}
		states := make([]hue.EntertainmentLightState, config.NumLights)
		for i := 0; i < config.NumLights; i++ {
			startIndex := i * 3
			states[i] = hue.EntertainmentLightState{
				On:  values[startIndex] > 0,
				Bri: int(values[startIndex+1]),
				XY:  [2]float64{float64(values[startIndex+2]) / 255.0, float64(values[startIndex+2]) / 255.0},
			}
		}
		hue.EntertainmentSend(config, states)
	})

	select {}
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IPP("hue-bridge-ip", "i", nil, "IP address of the hue bridge")
	serverCmd.Flags().StringP("username", "u", "", "Username for the hue bridge")
	serverCmd.Flags().IntP("entertainment-zone", "e", 0, "Entertainment zone ID for the hue bridge")
	serverCmd.Flags().IntP("lights", "l", 10, "Number of lights in the entertainment zone (default: 10)")
	serverCmd.Flags().Uint16P("artnet-universe", "n", 0, "Art-Net universe to listen on")
	serverCmd.Flags().IntP("artnet-dmx-start", "a", 1, "Art-Net DMX start channel")
	serverCmd.Flags().BoolP("debug", "d", false, "Debug mode (default: false)")
}
