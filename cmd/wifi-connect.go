package cmd

import (
	"log"

	"github.com/muka/pi-wifi/wifi"
	"github.com/spf13/cobra"
)

// wifiConnect represents the server command
var wifiConnect = &cobra.Command{
	Use:   "wifi-connect",
	Short: "Connect to wifi",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 2 {
			log.Fatal("Please provide ssid and password")
		}

		manager, err := wifi.NewManager()
		if err != nil {
			log.Fatal(err)
		}

		err = manager.Connect(args[0], args[1], "")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(wifiConnect)
}
