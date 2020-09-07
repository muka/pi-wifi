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

		if len(args) < 1 {
			log.Fatal("Please provide a connection string eg. WIFI:T:WPA;S:mynetwork;P:mypass;;")
		}

		connstr := args[0]
		connParams, err := wifi.ParseConnection(connstr)
		if err != nil {
			log.Fatalf("Failed to parse connection string: %s", err)
		}

		manager, err := wifi.NewManager()
		if err != nil {
			log.Fatal(err)
		}

		err = manager.Connect(connParams)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(wifiConnect)
}
