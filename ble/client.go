package ble

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

// Client creates a client
func Client(adapterID string, dev *device.Device1) (err error) {

	if adapterID == "" {
		return fmt.Errorf("Adapter name not provided")
	}

	if dev == nil {
		return fmt.Errorf("Device not provided")
	}

	//Connect DBus System bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	// do not reuse agent0 from service
	agent.NextAgentPath()

	ag := agent.NewSimpleAgent()
	err = agent.ExposeAgent(conn, ag, agent.CapNoInputNoOutput, true)
	if err != nil {
		return fmt.Errorf("SimpleAgent: %s", err)
	}

	// a, err := adapter.GetAdapterFromDevicePath(dev.Path())
	// if err != nil {
	// 	return err
	// }

	changes, err := dev.WatchProperties()
	go func() {
		for {
			select {
			case ev := <-changes:
				log.Infof("updated %s=%v", ev.Name, ev.Value)

				if !dev.Properties.Connected {
					err = connect(dev, ag, adapterID)
					if err != nil {
						log.Errorf("connect err: %s", err)
					}
				}

				break
			}
		}
	}()

	err = connect(dev, ag, adapterID)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	// retrieveServices(a, dev)

	select {}
	// return nil
}

func connect(dev *device.Device1, ag *agent.SimpleAgent, adapterID string) error {

	props, err := dev.GetProperties()
	if err != nil {
		return fmt.Errorf("Failed to load props: %s", err)
	}

	log.Infof("Found device name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)

	if props.Connected {
		log.Trace("Device is connected")
		return nil
	}

	if !props.Paired || !props.Trusted {
		log.Trace("Pairing device")

		err := dev.Pair()
		if err != nil {
			return fmt.Errorf("Pair failed: %s", err)
		}

		log.Info("Pair succeed")
		agent.SetTrusted(adapterID, dev.Path())
	}

	if !props.Connected {
		log.Info("Connecting device")
		err = dev.Connect()
		if err != nil {
			// if !strings.Contains(err.Error(), "Connection refused") {
			return fmt.Errorf("Connect failed: %s", err)
			// }
		}
	}

	return nil
}

func retrieveServices(a *adapter.Adapter1, dev *device.Device1) error {

	log.Info("Listing exposed services")

	list, err := dev.GetAllServicesAndUUID()
	if err != nil {
		return err
	}

	if len(list) == 0 {
		time.Sleep(time.Second * 2)
		return retrieveServices(a, dev)
	}

	for _, servicePath := range list {
		log.Debugf("%s", servicePath)
	}

	return nil
}
