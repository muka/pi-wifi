package runtime

import (
	"fmt"
	"strings"

	"github.com/muka/go-bluetooth/api/service"
	"github.com/muka/network_manager"
	"github.com/muka/pi-wifi/wifi"
	log "github.com/sirupsen/logrus"
)

func (r *Runtime) onConnectionRead(char *service.Char) {
	char.OnRead(service.CharReadCallback(func(c *service.Char, options map[string]interface{}) ([]byte, error) {

		log.Debug("onConnectionRead callback")

		conn, err := r.Wifi.GetConnectivity()
		if err != nil {
			log.Errorf("Failed to get connectivity: %s", err)
			return []byte{}, err
		}

		res := "disconnected"
		switch conn {
		case network_manager.NM_CONNECTIVITY_FULL:
			res = "connected"
			break
		case network_manager.NM_CONNECTIVITY_LIMITED:
			res = "limited"
			break
		case network_manager.NM_CONNECTIVITY_UNKNOWN:
			res = "unknown"
			break
		}

		return []byte(res), nil
	}))
}

func (r *Runtime) onConnectionWrite(char *service.Char) {
	char.OnWrite(service.CharWriteCallback(func(c *service.Char, value []byte) ([]byte, error) {
		log.Debugf("onConnectionWrite callback: %s", value)

		params, err := wifi.ParseConnection(string(value))
		if err != nil {
			log.Errorf("Error parsing connection parameters: %s", err)
			return []byte{}, err
		}

		err = r.Wifi.Connect(params)
		if err != nil {
			log.Errorf("Error connecting: %s", err)
			return []byte("Internal error"), err
		}

		return []byte("ok"), nil
	}))

}

func (r *Runtime) onAPList(char *service.Char) {
	char.OnRead(service.CharReadCallback(func(c *service.Char, options map[string]interface{}) ([]byte, error) {
		log.Debug("onAPList callback")

		list := []string{}

		devices, err := r.Wifi.GetWifiDevices()
		if err != nil {
			log.Errorf("Failed to list wifi devices: %s", err)
			return []byte{}, err
		}

		for _, device := range devices {
			aps, err := r.Wifi.GetAccessPoints(device.Path)
			if err != nil {
				log.Warnf("Error getting access points list for %s: %s", device.Interface, err)
				continue
			}
			for _, ap := range aps {
				list = append(list, fmt.Sprintf("%s;%s;%b\n", device.Interface, ap.SSID, ap.Strength))
			}
		}

		return []byte(strings.Join(list, "\n") + "\n"), nil
	}))

}
