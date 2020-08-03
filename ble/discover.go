package ble

import (
	"os"
	"os/signal"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

// Discover a device by name
func Discover(adapterID string) error {

	log.Info("Starting discovery")

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

			log.Infof("name=%s addr=%s rssi=%d", dev.Properties.Name, dev.Properties.Address, dev.Properties.RSSI)

		}

	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill) // get notified of all OS signals

	sig := <-ch
	log.Infof("Received signal [%v]; shutting down...\n", sig)
	return nil
}
