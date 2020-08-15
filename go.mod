module github.com/muka/pi-wifi

go 1.14

replace github.com/muka/go-bluetooth => ../go-bluetooth

require (
	github.com/godbus/dbus/v5 v5.0.3
	github.com/muka/go-bluetooth v0.0.0
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
)
