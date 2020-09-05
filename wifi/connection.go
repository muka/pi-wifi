package wifi

import (
	"context"
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/muka/network_manager"
	log "github.com/sirupsen/logrus"
)

// Connection wrap a managed connection
type Connection struct {
	ID         string
	Type       string
	SSID       string
	ObjectPath dbus.ObjectPath
}

// GetConnectionBySSID retrieve a manged connection by ssid
func (m *Manager) GetConnectionBySSID(ssid string) (connection Connection, err error) {

	connections, err := m.GetConnections()
	if err != nil {
		return connection, fmt.Errorf("GetConnections: %s", err)
	}

	for _, conn := range connections {
		if conn.SSID == ssid {
			return conn, nil
		}
	}

	return connection, fmt.Errorf("Connection configuration for %s not found", ssid)
}

// GetConnections retrieve all WIFI connections
func (m *Manager) GetConnections() (connections []Connection, err error) {

	connectionsPath, err := m.settings.GetConnections(context.Background())
	if err != nil {
		return connections, err
	}

	for _, connectionPath := range connectionsPath {

		connection := network_manager.NewNetworkManager_Settings_Connection(
			m.conn.Object(nmNs, connectionPath),
		)

		config, err := connection.GetSettings(context.Background())
		if err != nil {
			log.Errorf("Failed read settings for %s", connectionPath)
			continue
		}

		if connectionInfo, ok := config["connection"]; ok {
			if _, ok := connectionInfo["type"]; ok {
				if _, ok := connectionInfo["id"]; ok {

					connectionType := connectionInfo["type"].Value().(string)
					connectionID := connectionInfo["id"].Value().(string)

					if connectionType == "802-11-wireless" {
						if strings.HasPrefix(connectionID, ConnectionNamePrefix) {

							connectionType := connectionInfo["type"].Value().(string)
							connectionID := connectionInfo["id"].Value().(string)
							connectionSSID := string(config["802-11-wireless"]["ssid"].Value().([]byte))

							log.Debugf(
								"Found connection id=%s type=%s ssid=%s",
								connectionType,
								connectionID,
								connectionSSID,
							)

							connections = append(connections, Connection{
								Type:       connectionType,
								ID:         connectionID,
								SSID:       connectionSSID,
								ObjectPath: connectionPath,
							})
						}
					}
				}
			}
		}
	}

	return connections, nil
}
