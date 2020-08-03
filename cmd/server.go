package cmd

import (
	"log"

	"github.com/muka/pi-wifi/ble"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a GATT server to enable wifi connection",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if adapterID == "" {
			adapterID = "hci0"
		}

		if err := ble.Serve(adapterID); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
