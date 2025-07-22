/*
Copyright © 2025 Christiaan de Die le Clercq <contact@techwolf12.nl>

*/
package cmd

import (
	"fmt"
	"github.com/techwolf12/artnet-to-hue/pkg/config"
	"github.com/techwolf12/artnet-to-hue/pkg/hue"

	"github.com/spf13/cobra"
)

// bridgeInfoCmd represents the bridgeInfo command
var bridgeInfoCmd = &cobra.Command{
	Use:   "bridgeInfo",
	Short: "Get information about entertainment zones on the bridge",
	Long:  `Get information about entertainment zones on the bridge`,
	Run:   bridgeInfoHandler,
}

func bridgeInfoHandler(cmd *cobra.Command, args []string) {
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

	artnetConfig := config.Config{
		HueBridgeIP: hueBridgeIP,
		Username:    username,
	}

	entertainmentZones, err := hue.GetEntertainmentConfigID(artnetConfig)
	if err != nil {
		fmt.Printf("Error fetching entertainment zones: %v\n", err)
		return
	}
	for _, zone := range entertainmentZones {
		fmt.Printf("ID: %s, Name: %s, Status: %s\n", zone.ID, zone.Name, zone.Status)
	}
}

func init() {
	rootCmd.AddCommand(bridgeInfoCmd)

	bridgeInfoCmd.Flags().IPP("hue-bridge-ip", "i", nil, "IP address of the hue bridge")
	bridgeInfoCmd.Flags().StringP("username", "u", "", "Username for the hue bridge")
}
