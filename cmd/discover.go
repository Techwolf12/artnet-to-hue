/*
Copyright © 2025 Christiaan de Die le Clercq <contact@techwolf12.nl>

*/
package cmd

import (
	"fmt"
	"github.com/techwolf12/artnet-to-hue/pkg/hue"
	"log"

	"github.com/spf13/cobra"
)

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Discover Hue bridges on the network",
	Long:  `Discover Hue bridges on the local network and list their IP addresses and bridge ID.`,
	Run:   runDiscover,
}

func runDiscover(cmd *cobra.Command, args []string) {
	fmt.Println("Running discovery of Hue bridges...")
	hueBridges, err := hue.DiscoverBridges()
	if err != nil {
		log.Printf("Error discovering Hue bridges: %v\n", err)
		return
	}
	for _, bridge := range hueBridges {
		fmt.Printf("Found Hue bridge (%s), to use this run: artnet-to-hue pair -i %s\n", bridge.ID, bridge.InternalIPAddress)
	}
}
func init() {
	rootCmd.AddCommand(discoverCmd)
}
