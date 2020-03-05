package routers

// StormCluster godoc
type StormCluster struct {
	ServerName         string `json:"serverName"`
	Server             string `json:"server"`
	Token              string `json:"token"`
	ServerCADataString string `json:"serverCADataString"`
}

// AuthFriendlyErr godoc
type AuthFriendlyErr struct {
	HTTPStatus string `json: "httpStatus"`
	Message    string `json: "message"`
}
