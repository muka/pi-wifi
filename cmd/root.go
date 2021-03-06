package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var adapterID string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pi-wifi",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./pi-wifi.yaml)")
	rootCmd.PersistentFlags().StringVar(&adapterID, "adapter", "", "the bluetooth adapter to use (default is hci0)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".pi-wifi" (without extension).
		viper.AddConfigPath("./")
		viper.AddConfigPath("./config")
		viper.SetConfigName("pi-wifi")
	}

	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("log_level", "info")
	viper.SetDefault("service_name", "pi-wifi")

	viper.SetDefault("http_port", 9099)
	viper.SetDefault("http_public", true)
	viper.SetDefault("http_public_dir", "./public/")

	// viper.SetDefault("btmgmt_bin", "/usr/bin/btmgmt")
	viper.SetDefault("ble_adapter", "hci0")
	viper.SetDefault("ble_uuid_suffix", "-0000-1000-8000-00805f9b34fb")
	viper.SetDefault("ble_uuid_id", "1234")
	viper.SetDefault("ble_service_id", "2233")
	viper.SetDefault("ble_char_id_wifi", "3344")
	viper.SetDefault("ble_char_id_ap", "4455")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file %s", viper.ConfigFileUsed())
	}

	lvl, err := log.ParseLevel(viper.GetString("log_level"))
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(lvl)

}
