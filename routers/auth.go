package routers

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	b64 "encoding/base64"

	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
	clientcmd "k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// SetupKubeconfig godoc
// @Summary setup kubeconfig authentication
// @Description setup kubernetes client
func setupKubeconfigFromFile() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	catchError(err)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	catchError(err)

	return clientset
}

// postAuthConfig godoc
func postAuthConfig(authConfig *AuthConfig) (*version.Info, *AuthError) {

	var (
		serverName         = authConfig.ServerName
		server             = authConfig.Server
		token              = authConfig.Token
		serverCADataString = authConfig.ServerCADataString
		clientset          = &kubernetes.Clientset{}
		authError          = &AuthError{}
		serverVersion      = &version.Info{}
	)

	kubeconfigPath := homeDir() + "/.kubestorm/" + hash(server)
	serverHashPath := homeDir() + "/.kubestorm/" + serverName

	if isFileExist(kubeconfigPath) && !isFileExist(serverHashPath) {
		authError.Message = fmt.Sprintf("%v already exists", server)
		authError.Code = 409
		authError.Error = errors.New("")
	} else {
		if serverName != "" && server != "" && token != "" && serverCADataString != "" {
			serverCAData, _ := b64.StdEncoding.DecodeString(serverCADataString)
			kubeconfig := &clientcmdapi.Config{
				AuthInfos: map[string]*clientcmdapi.AuthInfo{
					serverName: {Token: token}},
				Clusters: map[string]*clientcmdapi.Cluster{
					serverName: {Server: server, CertificateAuthorityData: serverCAData, InsecureSkipTLSVerify: false}},
				Contexts: map[string]*clientcmdapi.Context{
					serverName: {AuthInfo: serverName, Cluster: serverName}},
				CurrentContext: serverName,
			}

			clientcmd.WriteToFile(*kubeconfig, "/tmp/"+serverName+"-kubeconfig")
			// use the current context in kubeconfig
			config, err := clientcmd.BuildConfigFromFlags(server, "/tmp/"+serverName+"-kubeconfig")
			if err != nil {
				authError.Message = "Not able to create clientcmdapi config"
				authError.Error = err
				authError.Code = 500
				return serverVersion, authError
			}

			// create the clientset
			clientset, err = kubernetes.NewForConfig(config)
			if err != nil {
				authError.Message = "Not able to create clientset"
				authError.Error = err
				authError.Code = 500
				return serverVersion, authError
			}

			// test the clientset
			serverVersion, err = clientset.ServerVersion()
			if err != nil {
				authError.Message = "Not able to retrieve server version. Something went wrong!"
				authError.Error = err
				authError.Code = 401
				return serverVersion, authError
			}

			clientcmd.WriteToFile(*kubeconfig, kubeconfigPath)
			if err := writeToFile([]byte(hash(server)), serverHashPath); err != nil {
				panic(err.Error())
			}

		} else {
			authError.Message = "Invalid Request"
			authError.Code = 400
			authError.Error = errors.New("Missing parameter")
		}
	}

	return serverVersion, authError
}

// getAuthConfig
func getAuthConfig(serverName string) (*AuthConfig, *AuthError) {
	var (
		authError  = &AuthError{}
		authConfig = &AuthConfig{}
	)

	hashID, err := getHashIDFromServerName(serverName)
	if err != nil {
		authError.Message = fmt.Sprintf("Not able to get config file for %v", serverName)
		authError.Code = 400
		authError.Error = err
		return authConfig, authError
	}

	kubeconfigPath := getKubeconfigPathFromHash(hashID)
	kubeconfig, err := clientcmd.LoadFromFile(kubeconfigPath)

	if err != nil {
		authError.Message = fmt.Sprintf("Not able to load kubeconfig for %v", serverName)
		authError.Code = 400
		authError.Error = err
		return authConfig, authError
	}

	authConfig.Server = kubeconfig.Clusters[serverName].Server
	authConfig.ServerName = serverName
	authConfig.Token = kubeconfig.AuthInfos[serverName].Token
	authConfig.ServerCADataString = b64.StdEncoding.EncodeToString(kubeconfig.Clusters[serverName].CertificateAuthorityData)

	return authConfig, authError
}

// deleteAuthConfig
func deleteAuthConfig(serverName string) *AuthError {
	var (
		authError = &AuthError{}
	)

	serverHashPath := homeDir() + "/.kubestorm/" + serverName
	if _, err := os.Stat(serverHashPath); err == nil {
		hashID, err := getHashIDFromServerName(serverName)
		if err != nil {
			authError.Message = fmt.Sprintf("Not able to read %v", serverHashPath)
			authError.Code = 400
			authError.Error = err
			return authError
		}
		kubeconfigPath := getKubeconfigPathFromHash(hashID)

		err = os.Remove(kubeconfigPath)
		if err != nil {
			authError.Message = fmt.Sprintf("Not able to remove %v", kubeconfigPath)
			authError.Code = 400
			authError.Error = err
			return authError
		}
		err = os.Remove(serverHashPath)
		if err != nil {
			authError.Message = fmt.Sprintf("Not able to remove %v", serverHashPath)
			authError.Code = 400
			authError.Error = err
			return authError
		}
	} else {
		authError.Message = fmt.Sprintf("%v does not exists!", serverHashPath)
		authError.Code = 404
		authError.Error = err
	}

	return authError
}

// getHashIDFromServerName godoc
func getHashIDFromServerName(serverName string) ([]byte, error) {
	return ioutil.ReadFile(homeDir() + "/.kubestorm/" + serverName)
}

// getKubeconfigPathFromHash godoc
func getKubeconfigPathFromHash(hashID []byte) string {
	return homeDir() + "/.kubestorm/" + string(hashID)
}

// getClientSet godoc
func getClientSet(serverName string) (*kubernetes.Clientset, *AuthError) {
	var (
		authError = &AuthError{}
		clientset = &kubernetes.Clientset{}
	)

	hashID, err := getHashIDFromServerName(serverName)
	if err != nil {
		authError.Message = fmt.Sprintf("Not able to get config for %v", serverName)
		authError.Code = 400
		authError.Error = err
		return clientset, authError
	}
	kubeconfigPath := getKubeconfigPathFromHash(hashID)
	kubeconfig, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		authError.Message = "Not able to load config from file"
		authError.Error = err
		authError.Code = 500
		return clientset, authError
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags(kubeconfig.Clusters[serverName].Server, kubeconfigPath)
	if err != nil {
		authError.Message = "Not able to create clientcmdapi config"
		authError.Error = err
		authError.Code = 500
		return clientset, authError
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		authError.Message = fmt.Sprintf("Not able to create clientset for %v", serverName)
		authError.Code = 400
		authError.Error = err
	}

	return clientset, authError
}

// generateKubeconfigAwsCli godoc
func generateKubeconfigAwsCli(clusterName string, roleArn string, awsRegion string) {
	//	Split the entire command up using ' ' as the delimeter
	parts := strings.Split("aws eks --region "+awsRegion+" update-kubeconfig --name "+strings.TrimSpace(clusterName)+" --role-arn "+strings.TrimSpace(roleArn), " ")
	bin := parts[0]
	args := parts[1:len(parts)]
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
