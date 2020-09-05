package wifi

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/godbus/dbus/v5"
	"github.com/muka/network_manager"
)

// AccessPoint wrap AP information
type AccessPoint struct {
	SSID     []byte
	Strength byte
}

func (m *Manager) scanAccessPoints(devicePath dbus.ObjectPath) error {

	wireless := network_manager.NewNetworkManager_Device_Wireless(m.conn.Object(nmNs, devicePath))

	// todo: avoid request new scan if scanning is already in progress

	err := wireless.RequestScan(context.Background(), map[string]dbus.Variant{})
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) getAccessPoints(devicePath dbus.ObjectPath) ([]AccessPoint, error) {

	wireless := network_manager.NewNetworkManager_Device_Wireless(m.conn.Object(nmNs, devicePath))

	list := []AccessPoint{}

	accessPoints, err := wireless.GetAccessPoints(context.Background())
	if err != nil {
		return list, err
	}

	for _, accessPointPath := range accessPoints {

		accessPoint := network_manager.NewNetworkManager_AccessPoint(m.conn.Object(nmNs, accessPointPath))

		ssid, err := accessPoint.GetSsid(context.Background())
		if err != nil {
			log.Errorf("Error on GetSsid: %s", err)
			continue
		}

		strength, err := accessPoint.GetStrength(context.Background())
		if err != nil {
			log.Errorf("Error on GetStrength: %s", err)
			continue
		}

		list = append(list, AccessPoint{
			SSID:     ssid,
			Strength: strength,
		})
	}

	return list, err
}
