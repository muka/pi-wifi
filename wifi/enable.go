package wifi

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// EnableWifi unlock wifi via network manager
func (m *Manager) EnableWifi() error {

	enabled, err := m.networkManager.GetWirelessEnabled(context.Background())
	if err != nil {
		return err
	}

	if !enabled {
		log.Debug("Enabling WIFI")
		err := m.networkManager.SetWirelessEnabled(context.Background(), true)
		if err != nil {
			return err
		}
	}

	log.Debug("WIFI is enabled")

	return nil
}
