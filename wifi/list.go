package wifi

import (
	"context"

	"github.com/godbus/dbus/v5"
	"github.com/muka/network_manager"
	log "github.com/sirupsen/logrus"
)

// Device wrap a NetworkManager_Device instance
type Device struct {
	Interface string
	Path      dbus.ObjectPath
	Device    *network_manager.NetworkManager_Device
}

// GetWifiDevices enumerate WIFI devices
func (m *Manager) GetWifiDevices() ([]Device, error) {

	list := []Device{}

	devices, err := m.networkManager.GetAllDevices(context.Background())
	if err != nil {
		return list, err
	}

	for _, devicePath := range devices {
		device := network_manager.NewNetworkManager_Device(m.conn.Object(nmNs, devicePath))

		deviceType, err := device.GetDeviceType(context.Background())
		if err != nil {
			log.Warnf("Error reading device type %s: %s", devicePath, err)
			continue
		}

		deviceInterface, err := device.GetInterface(context.Background())
		if err != nil {
			log.Warnf("Error reading device interface %s: %s", devicePath, err)
			continue
		}

		if network_manager.NM_DEVICE_TYPE_WIFI == deviceType {
			list = append(list, Device{
				Path:      devicePath,
				Device:    device,
				Interface: deviceInterface,
			})
		}
	}

	return list, nil
}
