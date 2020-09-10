package runtime

import (
	"github.com/muka/go-bluetooth/api/service"
	"github.com/muka/pi-wifi/ble"
	"github.com/muka/pi-wifi/server"
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

	instance.HTTP = server.NewHTTPServer(
		// connect
		func(connstr string) (string, error) {
			return instance.connect(connstr)
		},
		// status
		func() (string, error) {
			return instance.getStatus()
		},
		// list APs
		func() (aps []server.AccessPoint, err error) {

			list, err := instance.listAP()
			if err != nil {
				return nil, err
			}

			for _, ap1 := range list {
				aps = append(aps, server.AccessPoint{
					SSID:     string(ap1.SSID),
					Strength: int(ap1.Strength),
				})
			}

			return aps, nil
		},
	)
	if err != nil {
		return instance, err
	}

	return instance, nil
}

// Runtime handle instances of
type Runtime struct {
	Wifi            *wifi.Manager
	Ble             *service.App
	HTTP            server.HTTPServer
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

	// seems 0 is not working well
	// use uint32 max value

	cancel, err := r.Ble.Advertise(4294967295)

	r.CancelAdvertise = cancel
	if err != nil {
		return err
	}

	err = r.HTTP.Serve()
	if err != nil {
		return err
	}

	select {}
}
