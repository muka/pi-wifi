package wifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseConnection(t *testing.T, info ConnectionParams) {

	res, err := ParseConnection(info.String())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.AuthType, info.AuthType)
	assert.Equal(t, res.SSID, info.SSID)
	assert.Equal(t, res.Password, info.Password)
	assert.Equal(t, res.Hidden, info.Hidden)

}

func TestParseConnection1(t *testing.T) {

	info := ConnectionParams{
		AuthType: AuthTypeWPA,
		SSID:     "mynetwork",
		Password: "mypass",
	}

	parseConnection(t, info)
}

func TestParseConnection2(t *testing.T) {

	info := ConnectionParams{
		AuthType: AuthTypeWPAEAP,
		SSID:     `mynetwork\:\;`,
		Password: `my\:pass\;`,
	}

	parseConnection(t, info)
}
