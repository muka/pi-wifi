package wifi

import (
	"context"
	"errors"
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

// ActiveConnection wrap an active connection
type ActiveConnection struct {
	Connection
	ActiveConnectionPath dbus.ObjectPath
}

// GetConnection retrieve a manged connection by callback filter
func (m *Manager) GetConnection(fn func(conn Connection) bool) (connection Connection, err error) {

	connections, err := m.GetConnections()
	if err != nil {
		return connection, fmt.Errorf("GetConnections: %s", err)
	}

	for _, conn := range connections {
		if fn(conn) {
			return conn, nil
		}
	}

	return connection, errors.New("Connection not found")
}

// GetConnectionBySSID retrieve a manged connection by ssid
func (m *Manager) GetConnectionBySSID(ssid string) (conn Connection, err error) {
	conn, err = m.GetConnection(func(conn Connection) bool {
		return conn.SSID == ssid
	})
	if err != nil {
		err = fmt.Errorf("No connection by ssid=%s", ssid)
	}
	return conn, err
}

// GetConnectionByID retrieve a manged connection by ID
func (m *Manager) GetConnectionByID(id string) (conn Connection, err error) {
	conn, err = m.GetConnection(func(conn Connection) bool {
		return conn.ID == id
	})
	if err != nil {
		err = fmt.Errorf("No connection by id=%s", id)
	}
	return conn, err
}

// GetConnections retrieve all WIFI connections managed by the application
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

							log.Tracef(
								"Found connection id=%s type=%s ssid=%s",
								connectionID,
								connectionType,
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

// GetActiveConnections retrieve active WIFI connections managed by the application
func (m *Manager) GetActiveConnections() (connections []ActiveConnection, err error) {

	connectionsPath, err := m.networkManager.GetActiveConnections(context.Background())
	if err != nil {
		return connections, err
	}

	for _, connectionPath := range connectionsPath {

		connection := network_manager.NewNetworkManager_Connection_Active(
			m.conn.Object(nmNs, connectionPath),
		)

		connectionID, err := connection.GetId(context.Background())
		if err != nil {
			log.Errorf("Failed read ID for %s", connectionPath)
			continue
		}

		if strings.HasPrefix(connectionID, ConnectionNamePrefix) {

			conn, err := m.GetConnectionByID(connectionID)
			if err != nil {
				log.Errorf("Failed to load connection with ID=%s", err)
				continue
			}

			log.Tracef(
				"Found active connection id=%s type=%s ssid=%s",
				conn.ID,
				conn.Type,
				conn.SSID,
			)

			connections = append(connections, ActiveConnection{
				Connection:           conn,
				ActiveConnectionPath: connectionPath,
			})
		}
	}

	return connections, nil
}
