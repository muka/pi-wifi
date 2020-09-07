package ble

import (
	"github.com/muka/go-bluetooth/hw"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/muka/go-bluetooth/api/service"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

// NewService start a GATT server to expose credentials
func NewService() (*service.App, error) {

	adapterID := viper.GetString("ble_adapter")
	if adapterID == "" {
		adapterID = "hci0"
	}

	btmgmt := hw.NewBtMgmt(adapterID)
	if viper.GetString("btmgmt_bin") == "" {
		btmgmt.BinPath = viper.GetString("btmgmt_bin")
	}

	// set LE mode
	btmgmt.SetPowered(false)
	btmgmt.SetLe(true)
	btmgmt.SetBredr(false)
	btmgmt.SetPowered(true)

	options := service.AppOptions{
		AdapterID:  adapterID,
		AgentCaps:  agent.CapNoInputNoOutput,
		UUIDSuffix: viper.GetString("ble_uuid_suffix"),
		UUID:       viper.GetString("ble_uuid_id"),
	}

	a, err := service.NewApp(options)
	if err != nil {
		return nil, err
	}
	defer a.Close()

	a.SetName(viper.GetString("service_name"))

	log.Debugf("HW address %s", a.Adapter().Properties.Address)

	if !a.Adapter().Properties.Powered {
		err = a.Adapter().SetPowered(true)
		if err != nil {
			log.Fatalf("Failed to power the adapter: %s", err)
		}
	}

	service1, err := a.NewService(viper.GetString("ble_service_id"))
	if err != nil {
		return nil, err
	}

	err = a.AddService(service1)
	if err != nil {
		return nil, err
	}

	// char1 - write connection config and read status
	char1, err := service1.NewChar(viper.GetString("ble_char_id_wifi"))
	if err != nil {
		return nil, err
	}

	char1.Properties.Flags = []string{
		gatt.FlagCharacteristicRead,
		gatt.FlagCharacteristicWrite,
		gatt.FlagCharacteristicNotify,
	}

	err = service1.AddChar(char1)
	if err != nil {
		return nil, err
	}

	// char2 - list wifi connections

	char2, err := service1.NewChar(viper.GetString("ble_char_id_ap"))
	if err != nil {
		return nil, err
	}

	char2.Properties.Flags = []string{
		gatt.FlagCharacteristicRead,
		gatt.FlagCharacteristicNotify,
	}

	err = service1.AddChar(char2)
	if err != nil {
		return nil, err
	}

	err = a.Run()
	if err != nil {
		return nil, err
	}

	log.Infof("Exposed service %s", service1.Properties.UUID)

	// timeout := uint32(6 * 3600) // 6h
	// log.Infof("Advertising for %ds...", timeout)
	// cancel, err := a.Advertise(timeout)
	// if err != nil {
	// 	return nil, err
	// }

	// defer cancel()

	// wait := make(chan bool)
	// go func() {
	// 	time.Sleep(time.Duration(timeout) * time.Second)
	// 	wait <- true
	// }()

	// <-wait

	return a, nil
}
