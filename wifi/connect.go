package wifi

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/godbus/dbus/v5"
)

//ConnectionNamePrefix the prefix for a connection id to identify managed connections
var ConnectionNamePrefix = "piwifi__"

// Connect to a wifi network
func (m *Manager) Connect(ssid, password, encryption string) error {

	// enable wifi
	err := m.EnableWifi()
	if err != nil {
		return fmt.Errorf("EnableWifi: %s", err)
	}

	// get wifi devices
	devices, err := m.GetWifiDevices()

	if len(devices) == 0 {
		// todo: check if a device can be activated or unblocked
		return fmt.Errorf("No WIFI devices available")
	}

	device := devices[0]

	// create wifi connection if not exists
	connection, err := m.CreateWifiConnection(ssid, password, encryption)
	if err != nil {
		return fmt.Errorf("CreateWifiConnection: %s", err)
	}

	// try to connect
	activeConnectionPath, err := m.activateConnection(device.Path, connection.ObjectPath)
	if err != nil {
		return fmt.Errorf("activateConnection: %s", err)
	}

	log.Tracef("Activated connection path %s", activeConnectionPath)

	return nil
}

// CreateWifiConnection create a unique connection or update if already exists
func (m *Manager) CreateWifiConnection(ssid, password, encryption string) (conn Connection, err error) {

	connections, err := m.GetConnections()
	if err != nil {
		return conn, fmt.Errorf("GetConnections: %s", err)
	}

	var connection Connection
	for _, conn := range connections {
		if conn.SSID == ssid {
			connection = conn
		}
	}

	if connection.ID == "" {

		// create wifi connection if not exists
		uuid, err := uuid.NewUUID()
		if err != nil {
			return conn, fmt.Errorf("uuid: %s", err)
		}

		connectionConfig := createWifiConnection(uuid.String(), ssid, password, encryption)

		_, err = m.settings.AddConnection(context.Background(), connectionConfig)
		if err != nil {
			return conn, err
		}

		connection, err = m.GetConnectionBySSID(ssid)
		if err != nil {
			return conn, err
		}

		log.Debugf("Created connection %s", connection.ID)
	} else {
		log.Debugf("Connection %s exists", connection.ID)
	}

	return connection, nil
}

func (m *Manager) activateConnection(devicePath, connectionPath dbus.ObjectPath) (dbus.ObjectPath, error) {

	log.Tracef("%s %s", connectionPath, devicePath)

	activeConn, err := m.networkManager.ActivateConnection(
		context.Background(),
		connectionPath,
		devicePath,
		dbus.ObjectPath("/"), // select AP automatically
	)
	if err != nil {
		return dbus.ObjectPath(""), err
	}

	//todo: check if active

	log.Printf("Connection activated: %s", activeConn)

	return activeConn, nil
}

func createWifiConnection(uuid, ssid, password, encryption string) map[string]map[string]dbus.Variant {

	// uuid.FromBytes([]byte(ssid))

	wifi := map[string]dbus.Variant{
		"ssid": dbus.MakeVariant([]byte(ssid)),
		"mode": dbus.MakeVariant("infrastructure"),
	}

	conn := map[string]dbus.Variant{
		"type": dbus.MakeVariant("802-11-wireless"),
		"uuid": dbus.MakeVariant(uuid),
		"id":   dbus.MakeVariant(fmt.Sprintf("%s%s", ConnectionNamePrefix, ssid)),
	}

	ip4 := map[string]dbus.Variant{
		"method": dbus.MakeVariant("auto"),
	}
	ip6 := map[string]dbus.Variant{
		"method": dbus.MakeVariant("ignore"),
	}

	wsec := map[string]dbus.Variant{
		"key-mgmt": dbus.MakeVariant("wpa-psk"),
		"auth-alg": dbus.MakeVariant("open"),
		"psk":      dbus.MakeVariant(password),
	}

	con := map[string]map[string]dbus.Variant{
		"connection":               conn,
		"802-11-wireless":          wifi,
		"802-11-wireless-security": wsec,
		"ipv4":                     ip4,
		"ipv6":                     ip6,
	}

	return con
}
