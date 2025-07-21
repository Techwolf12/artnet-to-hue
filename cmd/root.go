/*
Copyright © 2025 Christiaan de Die le Clercq <contact@techwolf12.nl>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "artnet-to-hue",
	Short: "artnet-to-hue is a bridge between Art-Net and Philips Hue.",
	Long:  `artnet-to-hue is a bridge between Art-Net and Philips Hue. It allows you to control Philips Hue lights in an entertainment zone using Art-Net packets.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
