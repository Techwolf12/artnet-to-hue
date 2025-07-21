/*
Copyright © 2025 Christiaan de Die le Clercq <contact@techwolf12.nl>

*/
package cmd

import (
	"fmt"
	artnetHueConfig "github.com/techwolf12/artnet-to-hue/pkg/config"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Art-Net to Hue bridge server",
	Long:  `Start the Art-Net to Hue bridge server to listen for Art-Net packets and control Philips Hue lights.`,
	Run:   serverRun,
}

func serverRun(cmd *cobra.Command, args []string) {
	fmt.Println("server called")
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
	artnetUniverse, _ := cmd.Flags().GetInt("artnet-universe")
	if artnetUniverse < 0 {
		fmt.Println("Error: Art-Net universe must be a non-negative integer")
		return
	}
	artnetDMXStart, _ := cmd.Flags().GetInt("artnet-dmx-start")
	if artnetDMXStart < 0 {
		fmt.Println("Error: Art-Net DMX start channel must be a non-negative integer")
		return
	}
	_ = artnetHueConfig.Config{
		HueBridgeIP:        hueBridgeIP,
		Username:           username,
		EntertainmentZone:  entertainmentZone,
		ArtNetUniverse:     artnetUniverse,
		ArtNetStartAddress: artnetDMXStart,
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().IPP("hue-bridge-ip", "i", nil, "IP address of the hue bridge")
	serverCmd.Flags().StringP("username", "u", "", "Username for the hue bridge")
	serverCmd.Flags().IntP("entertainment-zone", "e", 0, "Entertainment zone ID for the hue bridge")
	serverCmd.Flags().IntP("artnet-universe", "n", 0, "Art-Net universe to listen on")
	serverCmd.Flags().IntP("artnet-dmx-start", "a", 0, "Art-Net DMX start channel")
}
