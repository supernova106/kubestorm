package routers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPods godoc
// @Summary Get Kubernetes Pods
// @Description get string by GET
// @Produce  json
// @Success 200 {object} Podlist
// @Router /pods/ [get]
func GetPods(c *gin.Context) {
	serverName := strings.TrimSpace(c.Param("name"))

	clientset, authFriendlyErr := getClientSet(serverName)
	if authFriendlyErr.HTTPStatus == "400" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    authFriendlyErr.Message,
			"httpStatus": authFriendlyErr.HTTPStatus,
		})
	} else if authFriendlyErr.HTTPStatus == "404" {
		c.JSON(http.StatusNotFound, gin.H{
			"message":    authFriendlyErr.Message,
			"httpStatus": authFriendlyErr.HTTPStatus,
		})
	} else {
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		catchError(err)
		c.JSON(http.StatusOK, pods)
	}
}

// GetNodes godoc
func GetNodes(c *gin.Context) {
	serverName := strings.TrimSpace(c.Param("name"))

	clientset, authFriendlyErr := getClientSet(serverName)
	if authFriendlyErr.HTTPStatus == "400" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    authFriendlyErr.Message,
			"httpStatus": authFriendlyErr.HTTPStatus,
		})
	} else if authFriendlyErr.HTTPStatus == "404" {
		c.JSON(http.StatusNotFound, gin.H{
			"message":    authFriendlyErr.Message,
			"httpStatus": authFriendlyErr.HTTPStatus,
		})
	} else {
		nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
		catchError(err)
		c.JSON(http.StatusOK, nodes)
	}
}

// GetNamespaces godoc
func GetNamespaces(c *gin.Context) {
	serverName := strings.TrimSpace(c.Param("name"))

	clientset, authFriendlyErr := getClientSet(serverName)
	if authFriendlyErr.HTTPStatus == "400" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    authFriendlyErr.Message,
			"httpStatus": authFriendlyErr.HTTPStatus,
		})
	} else if authFriendlyErr.HTTPStatus == "404" {
		c.JSON(http.StatusNotFound, gin.H{
			"message":    authFriendlyErr.Message,
			"httpStatus": authFriendlyErr.HTTPStatus,
		})
	} else {
		namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
		catchError(err)
		c.JSON(http.StatusOK, namespaces)
	}
}

// Auth godoc
func Auth(c *gin.Context) {
	serverName := strings.TrimSpace(c.Param("name"))

	switch requestMethod := c.Request.Method; requestMethod {
	case "POST":
		stormCluster := &StormCluster{
			ServerName:         serverName,
			Server:             strings.TrimSpace(c.PostForm("server")),
			Token:              strings.TrimSpace(c.PostForm("token")),
			ServerCADataString: strings.TrimSpace(c.PostForm("serverCADataString")),
		}

		_, authFriendlyErr := postClusterAuth(stormCluster)

		if authFriendlyErr.HTTPStatus == "400" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message":    authFriendlyErr.Message,
				"httpStatus": authFriendlyErr.HTTPStatus,
			})
		} else if authFriendlyErr.HTTPStatus == "401" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message":    authFriendlyErr.Message,
				"httpStatus": authFriendlyErr.HTTPStatus,
			})
		} else if authFriendlyErr.HTTPStatus == "409" {
			c.JSON(http.StatusConflict, gin.H{
				"message":    authFriendlyErr.Message,
				"httpStatus": authFriendlyErr.HTTPStatus,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("%v is updated successfully!", stormCluster.ServerName),
			})
		}
	case "GET":
		stormCluster := getClusterAuth(serverName)
		c.JSON(http.StatusOK, stormCluster)
	case "DELETE":
		authFriendlyErr := deleteClusterAuth(serverName)
		if authFriendlyErr.HTTPStatus == "404" {
			c.JSON(http.StatusNotFound, gin.H{
				"message":    authFriendlyErr.Message,
				"httpStatus": authFriendlyErr.HTTPStatus,
			})
		} else if authFriendlyErr.HTTPStatus == "200" {
			c.JSON(http.StatusOK, gin.H{
				"message":    authFriendlyErr.Message,
				"httpStatus": authFriendlyErr.HTTPStatus,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message":    "Something went wrong!",
				"httpStatus": http.StatusBadRequest,
			})
		}
	default:
	}
}
