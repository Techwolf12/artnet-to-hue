/*
Copyright © 2025 Christiaan de Die le Clercq <contact@techwolf12.nl>

*/
package cmd

import (
	"github.com/techwolf12/artnet-to-hue/pkg/hue"
	"log"

	"github.com/spf13/cobra"
)

// pairCmd represents the pair command
var pairCmd = &cobra.Command{
	Use:   "pair",
	Short: "Create a new user on the hue bridge.",
	Long:  `Create a new user on the Hue bridge. You can then use this user to run the server.`,
	Run:   pairRun,
}

func pairRun(cmd *cobra.Command, args []string) {
	hueBridgeIP, _ := cmd.Flags().GetIP("hue-bridge-ip")
	if hueBridgeIP.IsUnspecified() {
		log.Println("Error: Hue bridge IP address is required")
		return
	}
	deviceType := "artnet-to-hue#artnet-to-hue"
	username, clientKey, err := hue.GetHueUsername(hueBridgeIP, deviceType)
	if err != nil {
		log.Printf("Error getting Hue username: %v\n", err)
		return
	}
	log.Println("Be sure to read the help for server. You can now use this username to run the server with the following command:")
	log.Printf("artnet-to-hue server -i %s -u %s -c %s\n", hueBridgeIP, username, clientKey)
}

func init() {
	rootCmd.AddCommand(pairCmd)

	pairCmd.Flags().IPP("hue-bridge-ip", "i", nil, "IP address of the hue bridge")
}
