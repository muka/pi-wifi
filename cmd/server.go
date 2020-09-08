package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/muka/pi-wifi/runtime"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a GATT server to enable wifi connection",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		instance, err := runtime.NewRuntime()
		if err != nil {
			log.Fatal(err)
		}

		defer instance.Stop()

		err = instance.Start()
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
