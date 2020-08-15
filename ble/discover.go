package ble

import (
	"os"
	"os/signal"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/muka/go-bluetooth/hw"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Discover a device by name
func Discover(adapterID string) error {

	log.SetLevel(log.TraceLevel)

	log.Info("Starting discovery")

	btmgmt := hw.NewBtMgmt(adapterID)
	if len(os.Getenv("DOCKER")) > 0 {
		btmgmt.BinPath = viper.GetString("btmgmt_bin")
	}

	// set LE mode
	btmgmt.SetPowered(false)
	btmgmt.SetLe(true)
	btmgmt.SetBredr(false)
	btmgmt.SetDiscoverable(false)
	btmgmt.SetPowered(true)

	//clean up connection on exit
	defer api.Exit()

	a, err := adapter.GetAdapter(adapterID)
	if err != nil {
		return err
	}

	log.Debug("Flush cached devices")
	err = a.FlushDevices()
	if err != nil {
		return err
	}

	log.Debug("Start discovery")
	discovery, cancel, err := api.Discover(a, nil)
	if err != nil {
		return err
	}
	defer cancel()

	go func() {
		devices := map[string]*device.Device1{}

		for ev := range discovery {

			if ev.Type == adapter.DeviceRemoved {
				continue
			}

			dev, err := device.NewDevice1(ev.Path)
			if err != nil {
				log.Errorf("%s: %s", ev.Path, err)
				continue
			}

			if dev == nil {
				log.Errorf("%s: not found", ev.Path)
				continue
			}

			log.Infof("Found name=%s addr=%s rssi=%d", dev.Properties.Name, dev.Properties.Address, dev.Properties.RSSI)

			if _, ok := devices[dev.Properties.Address]; ok {
				log.Warnf("Skip duplicated address %s", dev.Properties.Address)
				continue
			}

			devices[dev.Properties.Address] = dev

			if dev.Properties.Name == viper.GetString("service_name") {
				// stop discovery
				cancel()
				err = Client(adapterID, dev)
			}

		}

	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill) // get notified of all OS signals

	sig := <-ch
	log.Infof("Received signal [%v]; shutting down...\n", sig)
	return nil
}
