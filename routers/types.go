package routers

import v1 "k8s.io/api/core/v1"

// AuthConfig godoc
type AuthConfig struct {
	ServerName         string `json:"serverName" binding:"required"`
	Server             string `json:"server" binding:"required"`
	Token              string `json:"token" binding:"required"`
	ServerCADataString string `json:"serverCADataString" binding:"required"`
}

// AuthError godoc
type AuthError struct {
	Error   error  `json: "error"`
	Code    int    `json: "code"`
	Message string `json: "message"`
}

// PodList godoc
type PodList v1.PodList
