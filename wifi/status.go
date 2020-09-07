package wifi

import (
	"context"
)

//GetConnectivity return an enum for the connectivity
func (m *Manager) GetConnectivity() (uint32, error) {
	return m.networkManager.GetConnectivity(context.Background())
}
