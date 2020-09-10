package runtime

import (
	"github.com/muka/network_manager"
	"github.com/muka/pi-wifi/wifi"
	log "github.com/sirupsen/logrus"
)

func (r *Runtime) listAP() (list []wifi.AccessPoint, err error) {

	devices, err := r.Wifi.GetWifiDevices()
	if err != nil {
		log.Errorf("Failed to list wifi devices: %s", err)
		return list, err
	}

	for _, device := range devices {
		aps, err := r.Wifi.GetAccessPoints(device.Path)
		if err != nil {
			log.Warnf("Error getting access points list for %s: %s", device.Interface, err)
			continue
		}
		for _, ap := range aps {
			list = append(list, ap)
		}
	}

	return list, nil
}

func (r *Runtime) connect(connstr string) (string, error) {

	params, err := wifi.ParseConnection(connstr)
	if err != nil {
		log.Errorf("Error parsing connection parameters: %s", err)
		return "parse_failure", err
	}

	err = r.Wifi.Connect(params)
	if err != nil {
		log.Errorf("Error connecting: %s", err)
		return "conn_failure", err
	}

	return "ok", nil
}

func (r *Runtime) getStatus() (string, error) {

	conn, err := r.Wifi.GetConnectivity()
	if err != nil {
		log.Errorf("Failed to get connectivity: %s", err)
		return "", err
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

	return res, nil
}
