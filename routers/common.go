package routers

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
		return true
	}

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
		panic(e.Error())
	}
}

// writeToFie godoc
func writeToFile(data []byte, fileName string) error {
	return ioutil.WriteFile(fileName, data, 0644)
}

// Error godoc
func Error(c *gin.Context, authError *AuthError) bool {
	if authError.Error != nil {
		c.Error(authError.Error)
		c.AbortWithStatusJSON(authError.Code, gin.H{"status": false, "error": authError.Error.Error(), "message": authError.Message})
		return true // signal that there was an error and the caller should return
	}
	return false // no error, can continue
}

// ErrorK8sClient godoc
func ErrorK8sClient(c *gin.Context, err error) bool {
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{"status": false, "error": err.Error(), "message": "something went wrong"})
		return true // signal that there was an error and the caller should return
	}
	return false // no error, can continue
}
