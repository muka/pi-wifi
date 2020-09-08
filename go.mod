module github.com/muka/pi-wifi

go 1.14

replace github.com/muka/go-bluetooth => ../go-bluetooth

require (
	github.com/godbus/dbus/v5 v5.0.3
	github.com/google/uuid v1.1.2
	github.com/muka/go-bluetooth v0.0.0
	github.com/muka/network_manager v0.0.0-20200903202308-ae5ede816e07
	github.com/prometheus/common v0.4.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
)
