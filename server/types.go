package server

type connectRequest struct {
	Payload string `json:"payload"`
}

type statusResponse struct {
	Status string `json:"status"`
}

// AccessPoint json response
type AccessPoint struct {
	SSID     string `json:"ssid"`
	Strength int    `json:"strength"`
}

type listAPResponse struct {
	AccessPoints []AccessPoint `json:"accessPoints"`
}
