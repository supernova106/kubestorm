package routers

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	klog "k8s.io/klog"
)

// APIStatus godoc
type APIStatus struct {
	Status string
}

// GetStatus godoc
// @Summary API Healthcheck
// @Description get string by GET
// @Produce  json
// @Success 200 {object} APIStatus
// @Router /status [get]
func GetStatus(c *gin.Context) {
	response := &APIStatus{
		Status: "ok",
	}
	c.JSON(http.StatusOK, response)
}

// homeDir godoc
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// isFileExist godoc
func isFileExist(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {
		klog.V(1).Infof("File %v exists", fileName)
		return true
	}

	klog.V(1).Infof("File %v does not exist\n", fileName)
	return false
}

// hash godoc
func hash(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

// catchError godoc
func catchError(e error) {
	if e != nil {
		panic(e)
	}
}

// writeToFie godoc
func writeToFile(data []byte, fileName string) {
	err := ioutil.WriteFile(fileName, data, 0644)
	catchError(err)
}
