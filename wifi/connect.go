package wifi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/godbus/dbus/v5"
)

//ConnectionNamePrefix the prefix for a connection id to identify managed connections
var ConnectionNamePrefix = "piwifi__"

// Connect to a wifi network
func (m *Manager) Connect(connectionParams ConnectionParams) error {

	if connectionParams.AuthType != AuthTypeWPA {
		return errors.New("Only WPA authentication type is supported at this time")
	}

	// enable wifi
	err := m.EnableWifi()
	if err != nil {
		return fmt.Errorf("EnableWifi: %s", err)
	}

	log.Debugf("Connecting..")

	// get wifi devices
	devices, err := m.GetWifiDevices()

	if len(devices) == 0 {
		// todo: check if a device can be activated or unblocked
		return fmt.Errorf("No WIFI devices available")
	}

	device := devices[0]

	// create wifi connection if not exists
	connection, err := m.CreateWifiConnection(connectionParams)
	if err != nil {
		return fmt.Errorf("CreateWifiConnection: %s", err)
	}

	// try to connect
	_, err = m.activateConnection(device.Path, connection.ObjectPath)
	if err != nil {
		return fmt.Errorf("activateConnection: %s", err)
	}

	log.Debugf("Connection initiated")

	return nil
}

// CreateWifiConnection create a unique connection or update if already exists
func (m *Manager) CreateWifiConnection(connectionParams ConnectionParams) (conn Connection, err error) {

	connections, err := m.GetConnections()
	if err != nil {
		return conn, fmt.Errorf("GetConnections: %s", err)
	}

	var connection Connection
	for _, conn := range connections {
		if conn.SSID == connectionParams.SSID {
			connection = conn
		}
	}

	if connection.ID == "" {

		// create wifi connection if not exists
		uuid, err := uuid.NewUUID()
		if err != nil {
			return conn, fmt.Errorf("uuid: %s", err)
		}

		connectionConfig := createWifiConnection(uuid.String(), connectionParams)

		_, err = m.settings.AddConnection(context.Background(), connectionConfig)
		if err != nil {
			return conn, err
		}

		connection, err = m.GetConnectionBySSID(connection.SSID)
		if err != nil {
			if !strings.Contains(err.Error(), "No connection by ssid") {
				return conn, fmt.Errorf("GetConnectionBySSID: %s", err)
			}
		}

		log.Tracef("Created connection %s", connection.ID)
	} else {
		log.Tracef("Connection %s exists", connection.ID)
	}

	return connection, nil
}

func (m *Manager) activateConnection(devicePath, connectionPath dbus.ObjectPath) (dbus.ObjectPath, error) {

	log.Tracef("%s %s", connectionPath, devicePath)

	activeConnections, err := m.GetActiveConnections()
	if err != nil {
		return dbus.ObjectPath(""), fmt.Errorf("GetActiveConnections: %s", err)
	}

	for _, ac := range activeConnections {
		if ac.ObjectPath == connectionPath {
			log.Tracef("Connection is already active")
			return ac.ActiveConnectionPath, nil
		}
	}

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
	// log.Printf("Connection activated: %s", activeConn)

	return activeConn, nil
}

func createWifiConnection(uuid string, connectionParams ConnectionParams) map[string]map[string]dbus.Variant {

	// uuid.FromBytes([]byte(ssid))

	wifi := map[string]dbus.Variant{
		"ssid": dbus.MakeVariant([]byte(connectionParams.SSID)),
		"mode": dbus.MakeVariant("infrastructure"),
	}

	conn := map[string]dbus.Variant{
		"type": dbus.MakeVariant("802-11-wireless"),
		"uuid": dbus.MakeVariant(uuid),
		"id":   dbus.MakeVariant(fmt.Sprintf("%s%s", ConnectionNamePrefix, connectionParams.SSID)),
	}

	ip4 := map[string]dbus.Variant{
		"method": dbus.MakeVariant("auto"),
	}
	ip6 := map[string]dbus.Variant{
		"method": dbus.MakeVariant("ignore"),
	}

	// todo: handle other connection types
	keyMgm := ""
	if connectionParams.AuthType == AuthTypeWPA {
		keyMgm = "wpa-psk"
	}

	wsec := map[string]dbus.Variant{
		"key-mgmt": dbus.MakeVariant(keyMgm),
		"auth-alg": dbus.MakeVariant("open"),
		"psk":      dbus.MakeVariant(connectionParams.Password),
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
