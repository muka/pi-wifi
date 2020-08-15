package cmd

import (
	"log"

	"github.com/muka/pi-wifi/ble"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start a client to connect with a server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if adapterID == "" {
			adapterID = "hci0"
		}

		err := ble.Discover(adapterID)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
