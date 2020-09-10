package runtime

import (
	"fmt"
	"strings"

	"github.com/muka/go-bluetooth/api/service"
	log "github.com/sirupsen/logrus"
)

func (r *Runtime) onConnectionRead(char *service.Char) {
	char.OnRead(service.CharReadCallback(func(c *service.Char, options map[string]interface{}) ([]byte, error) {
		log.Debug("onConnectionRead callback")
		res, err := r.getStatus()
		return []byte(res), err
	}))
}

func (r *Runtime) onConnectionWrite(char *service.Char) {
	char.OnWrite(service.CharWriteCallback(func(c *service.Char, value []byte) ([]byte, error) {
		log.Debugf("onConnectionWrite callback: %s", value)
		res, err := r.connect(string(value))
		return []byte(res), err
	}))

}

func (r *Runtime) onAPList(char *service.Char) {
	char.OnRead(service.CharReadCallback(func(c *service.Char, options map[string]interface{}) ([]byte, error) {
		log.Debug("onAPList callback")

		aps, err := r.listAP()
		if err != nil {
			log.Errorf("Failed to list wifi devices: %s", err)
			return []byte{}, err
		}

		list := []string{}
		for _, ap := range aps {
			list = append(list, fmt.Sprintf("%s;%b\n", ap.SSID, ap.Strength))
		}

		return []byte(strings.Join(list, "\n") + "\n"), nil
	}))

}
