package runtime

import (
	"strings"

	"github.com/muka/go-bluetooth/api/service"
	"github.com/muka/pi-wifi/ble"
	"github.com/muka/pi-wifi/wifi"
	"github.com/spf13/viper"
)

//NewRuntime create a new Runtime instance
func NewRuntime() (instance Runtime, err error) {

	instance.Wifi, err = wifi.NewManager()
	if err != nil {
		return instance, err
	}

	instance.Ble, err = ble.NewService()
	if err != nil {
		return instance, err
	}

	services := instance.Ble.GetServices()
	for _, service := range services {
		for _, char := range service.GetChars() {
			switch strings.Split(char.UUID, "-")[0] {
			// connection
			case viper.GetString("ble_char_id_wifi"):
				instance.onConnectionRead(char)
				instance.onConnectionWrite(char)
				break
			// ap list
			case viper.GetString("ble_char_id_ap"):
				instance.onAPList(char)
				break
			}
		}
	}

	return instance, nil
}

// Runtime handle instances of
type Runtime struct {
	Wifi *wifi.Manager
	Ble  *service.App
}

// Start initialize the runtime
func (r *Runtime) Start() {

	// check if wifi is connects
	// subscrbe to connectivity changes
	// start BLE advertising if not

}
