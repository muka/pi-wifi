package runtime

import (
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
			switch char.UUID[4:8] {
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
	Wifi            *wifi.Manager
	Ble             *service.App
	CancelAdvertise func()
}

// Stop stop the runtime
func (r *Runtime) Stop() {
	r.CancelAdvertise()
	r.Ble.Close()
}

// Start initialize the runtime
func (r *Runtime) Start() error {

	// check if wifi is connects
	// subscrbe to connectivity changes
	// start BLE advertising if not

	err := r.Ble.Run()
	if err != nil {
		return err
	}

	cancel, err := r.Ble.Advertise(0)

	r.CancelAdvertise = cancel
	if err != nil {
		return err
	}

	select {}
}
