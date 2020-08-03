package ble

import (
	"os"
	"time"

	"github.com/muka/go-bluetooth/hw"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/muka/go-bluetooth/api/service"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

//Serve start a GATT server to expose credentials
func Serve(adapterID string) error {

	log.SetLevel(log.TraceLevel)

	btmgmt := hw.NewBtMgmt(adapterID)
	if len(os.Getenv("DOCKER")) > 0 {
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
		UUIDSuffix: "-0000-1000-8000-00805F9B34FB",
		UUID:       "1234",
	}

	a, err := service.NewApp(options)
	if err != nil {
		return err
	}
	defer a.Close()

	a.SetName(viper.GetString("machine_name"))

	log.Infof("HW address %s", a.Adapter().Properties.Address)

	if !a.Adapter().Properties.Powered {
		err = a.Adapter().SetPowered(true)
		if err != nil {
			log.Fatalf("Failed to power the adapter: %s", err)
		}
	}

	service1, err := a.NewService("2233")
	if err != nil {
		return err
	}

	err = a.AddService(service1)
	if err != nil {
		return err
	}

	char1, err := service1.NewChar("3344")
	if err != nil {
		return err
	}

	char1.Properties.Flags = []string{
		gatt.FlagCharacteristicRead,
		gatt.FlagCharacteristicWrite,
	}

	char1.OnRead(service.CharReadCallback(func(c *service.Char, options map[string]interface{}) ([]byte, error) {
		log.Warnf("GOT READ REQUEST")
		return []byte{42}, nil
	}))

	char1.OnWrite(service.CharWriteCallback(func(c *service.Char, value []byte) ([]byte, error) {
		log.Warnf("GOT WRITE REQUEST")
		return value, nil
	}))

	err = service1.AddChar(char1)
	if err != nil {
		return err
	}

	descr1, err := char1.NewDescr("4455")
	if err != nil {
		return err
	}

	descr1.Properties.Flags = []string{
		gatt.FlagDescriptorRead,
		gatt.FlagDescriptorWrite,
	}

	descr1.OnRead(service.DescrReadCallback(func(c *service.Descr, options map[string]interface{}) ([]byte, error) {
		log.Warnf("GOT READ REQUEST")
		return []byte{42}, nil
	}))
	descr1.OnWrite(service.DescrWriteCallback(func(d *service.Descr, value []byte) ([]byte, error) {
		log.Warnf("GOT WRITE REQUEST")
		return value, nil
	}))

	err = char1.AddDescr(descr1)
	if err != nil {
		return err
	}

	err = a.Run()
	if err != nil {
		return err
	}

	log.Infof("Exposed service %s", service1.Properties.UUID)

	timeout := uint32(6 * 3600) // 6h
	log.Infof("Advertising for %ds...", timeout)
	cancel, err := a.Advertise(timeout)
	if err != nil {
		return err
	}

	defer cancel()

	wait := make(chan bool)
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		wait <- true
	}()

	<-wait

	return nil
}
