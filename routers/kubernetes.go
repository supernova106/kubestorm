package routers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetResources godoc
func GetResources(c *gin.Context) {
	cluster := c.Query("cluster")
	resourcesType := c.Query("type")

	clientset, authError := getClientSet(cluster)
	if !Error(c, authError) {
		switch resourcesType {
		case "services":
			data, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
			if !ErrorK8sClient(c, err) {
				c.JSON(http.StatusOK, data)
			}
		case "serviceaccounts":
			data, err := clientset.CoreV1().ServiceAccounts("").List(metav1.ListOptions{})
			if !ErrorK8sClient(c, err) {
				c.JSON(http.StatusOK, data)
			}
		case "limitranges":
			data, err := clientset.CoreV1().LimitRanges("").List(metav1.ListOptions{})
			if !ErrorK8sClient(c, err) {
				c.JSON(http.StatusOK, data)
			}
		case "pods":
			data, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
			if !ErrorK8sClient(c, err) {
				c.JSON(http.StatusOK, data)
			}
		case "podtemplates":
			data, err := clientset.CoreV1().PodTemplates("").List(metav1.ListOptions{})
			if !ErrorK8sClient(c, err) {
				c.JSON(http.StatusOK, data)
			}
		case "nodes":
			data, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
			if !ErrorK8sClient(c, err) {
				c.JSON(http.StatusOK, data)
			}
		case "configmaps":
			data, err := clientset.CoreV1().ConfigMaps("").List(metav1.ListOptions{})
			if !ErrorK8sClient(c, err) {
				c.JSON(http.StatusOK, data)
			}
		case "namespaces":
			data, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
			if !ErrorK8sClient(c, err) {
				c.JSON(http.StatusOK, data)
			}
		case "events":
			data, err := clientset.CoreV1().Events("").List(metav1.ListOptions{})
			if !ErrorK8sClient(c, err) {
				c.JSON(http.StatusOK, data)
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid request",
				"statue":  false,
			})
		}
	}
}

// Auth godoc
func Auth(c *gin.Context) {
	cluster := strings.TrimSpace(c.Param("cluster"))

	switch requestMethod := c.Request.Method; requestMethod {
	case "POST":
		authConfig := &AuthConfig{
			ServerName:         cluster,
			Server:             strings.TrimSpace(c.PostForm("server")),
			Token:              strings.TrimSpace(c.PostForm("token")),
			ServerCADataString: strings.TrimSpace(c.PostForm("serverCADataString")),
		}

		serverVersion, authError := postAuthConfig(authConfig)
		if !Error(c, authError) {
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("%v is updated successfully!", authConfig.ServerName),
				"statue":  true,
				"data":    serverVersion,
			})
		}

	case "GET":
		authConfig, authError := getAuthConfig(cluster)
		if !Error(c, authError) {
			c.JSON(http.StatusOK, authConfig)
		}

	case "DELETE":
		authError := deleteAuthConfig(cluster)
		if !Error(c, authError) {
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("%v is successfully removed!", cluster),
				"statue":  true,
			})
		}
	default:
	}
}
