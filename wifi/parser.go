package wifi

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// AuthType Wifi authentication mechanism
type AuthType string

const (
	// AuthTypeWEP WEP authentication
	AuthTypeWEP = "WEP"
	// AuthTypeWPA WPA authentication
	AuthTypeWPA = "WPA"
	// AuthTypeWPAEAP WPA2-EAP authentication
	AuthTypeWPAEAP = "WPA2-EAP"
	// AuthTypeNopass no password
	AuthTypeNopass = "nopass"
)

// ConnectionParams wrap connection information from QRCode like format
type ConnectionParams struct {
	// T	WPA	Authentication type; can be WEP or WPA or WPA2-EAP, or nopass for no password. Or, omit for no password.
	AuthType AuthType
	// S	mynetwork	Network SSID. Required. Enclose in double quotes if it is an ASCII name, but could be interpreted as hex (i.e. "ABCD")
	SSID string
	// P	mypass	Password, ignored if T is nopass (in which case it may be omitted). Enclose in double quotes if it is an ASCII name, but could be interpreted as hex (i.e. "ABCD")
	Password string
	// H	true	Optional. True if the network SSID is hidden. Note this was mistakenly also used to specify phase 2 method in releases up to 4.7.8 / Barcode Scanner 3.4.0. If not a boolean, it will be interpreted as phase 2 method (see below) for backwards-compatibility
	Hidden bool
	// E	TTLS	(WPA2-EAP only) EAP method, like TTLS or PWD
	TTLS string
	// A	anon	(WPA2-EAP only) Anonymous identity
	Anon string
	// I	myidentity	(WPA2-EAP only) Identity
	Identity string
	// Phase2 PH2	MSCHAPV2	(WPA2-EAP only) Phase 2 method, like MSCHAPV2
	Phase2 string
}

func (c *ConnectionParams) String() string {

	params := fmt.Sprintf(
		"WIFI:T:%s;S:%s;P:%s;H:%t;",
		c.AuthType,
		c.SSID,
		c.Password,
		c.Hidden,
	)

	if c.AuthType == AuthTypeWPAEAP {
		params += fmt.Sprintf(
			"E:%s;A:%s;I:%s;PH2:%s;",
			c.TTLS,
			c.Anon,
			c.Identity,
			c.Phase2,
		)
	}

	return params + ";"
}

// ParseConnection parse config from string
// eg. WIFI:T:WPA;S:mynetwork;P:mypass;;
func ParseConnection(raw string) (params ConnectionParams, err error) {

	if len(raw) < 10 {
		return params, errors.New("Connections string is empty or too short")
	}

	if strings.ToUpper(raw[0:4]) != "WIFI" {
		return params, errors.New("Missing WIFI: prefix")
	}

	raw = raw[5:]

	if raw[len(raw)-2:] == ";;" {
		raw = raw[:len(raw)-2]
	}

	sectionRegex := regexp.MustCompile(`([^\\];{1})`)
	sections := sectionRegex.FindAllStringIndex(raw, -1)

	lastPos := 0
	for _, section := range sections {

		part := raw[lastPos : section[1]-1]
		lastPos = section[1]

		if len(part) == 0 {
			continue
		}

		if len(part) > 3 && part[:3] == "PH2" {
			params.Phase2 = part[3:]
			continue
		}

		sectionName := part[:1]
		sectionValue := part[2:]

		switch sectionName {
		case "T":
			var authType AuthType
			switch sectionValue {
			case "WPA2-EAP":
				authType = AuthTypeWPAEAP
				break
			case "WEP":
				authType = AuthTypeWEP
				break
			case "nopass":
				authType = AuthTypeNopass
				break
			case "WPA":
				authType = AuthTypeWPA
				break
			}
			params.AuthType = authType
			break
		case "S":
			params.SSID = sectionValue
			break
		case "P":
			params.Password = sectionValue
			break
		case "H":
			var hidden bool
			if sectionValue == "true" {
				hidden = true
			}
			params.Hidden = hidden
			break
		case "E":
			params.TTLS = sectionValue
			break
		case "A":
			params.Anon = sectionValue
			break
		case "I":
			params.Identity = sectionValue
			break
		}
	}

	return params, nil
}
