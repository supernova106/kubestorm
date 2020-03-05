package routers

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	b64 "encoding/base64"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

// postClusterAuth godoc
func postClusterAuth(stormCluster *StormCluster) (*kubernetes.Clientset, *AuthFriendlyErr) {

	var (
		serverName         = stormCluster.ServerName
		server             = stormCluster.Server
		token              = stormCluster.Token
		serverCADataString = stormCluster.ServerCADataString
		clientset          = &kubernetes.Clientset{}
		authFriendlyErr    = &AuthFriendlyErr{}
	)

	kubeconfigPath := homeDir() + "/.kubestorm/" + hash(server)
	serverHashPath := homeDir() + "/.kubestorm/" + serverName

	if isFileExist(kubeconfigPath) && !isFileExist(serverHashPath) {
		authFriendlyErr.Message = fmt.Sprintf("%v already exists", server)
		authFriendlyErr.HTTPStatus = "409"
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
			catchError(err)

			// create the clientset
			clientset, err = kubernetes.NewForConfig(config)
			catchError(err)

			// test the clientset
			serverVersion, err := clientset.ServerVersion()
			if err != nil {
				authFriendlyErr.Message = fmt.Sprintf("Authentication Failed. %v", err)
				authFriendlyErr.HTTPStatus = "401"
			} else {
				clientcmd.WriteToFile(*kubeconfig, kubeconfigPath)
				writeToFile([]byte(hash(server)), serverHashPath)
				authFriendlyErr.Message = fmt.Sprintf("%v", serverVersion)
			}
		} else {
			authFriendlyErr.Message = "Invalid Request"
			authFriendlyErr.HTTPStatus = "400"
		}
	}

	return clientset, authFriendlyErr
}

// getClusterAuth
func getClusterAuth(serverName string) *StormCluster {
	serverHashPath := homeDir() + "/.kubestorm/" + serverName
	hashID, err := ioutil.ReadFile(serverHashPath)
	catchError(err)
	kubeconfigPath := homeDir() + "/.kubestorm/" + string(hashID)
	kubeconfig, err := clientcmd.LoadFromFile(kubeconfigPath)
	catchError(err)
	return &StormCluster{
		Server:             kubeconfig.Clusters[serverName].Server,
		ServerName:         serverName,
		Token:              kubeconfig.AuthInfos[serverName].Token,
		ServerCADataString: b64.StdEncoding.EncodeToString(kubeconfig.Clusters[serverName].CertificateAuthorityData),
	}
}

// deleteClusterAuth
func deleteClusterAuth(serverName string) *AuthFriendlyErr {
	var (
		authFriendlyErr = &AuthFriendlyErr{}
	)

	serverHashPath := homeDir() + "/.kubestorm/" + serverName
	if !isFileExist(serverHashPath) {
		authFriendlyErr.Message = fmt.Sprintf("%v does not exists!", serverName)
		authFriendlyErr.HTTPStatus = "404"
	} else {
		hashID, err := ioutil.ReadFile(serverHashPath)
		catchError(err)
		kubeconfigPath := homeDir() + "/.kubestorm/" + string(hashID)

		err = os.Remove(kubeconfigPath)
		err = os.Remove(serverHashPath)
		if err != nil {
			authFriendlyErr.Message = fmt.Sprintf("%v", err)
		} else {
			authFriendlyErr.Message = fmt.Sprintf("%v is successfully removed", serverName)
			authFriendlyErr.HTTPStatus = "200"
		}
	}

	return authFriendlyErr
}

// getClientSet godoc
func getClientSet(serverName string) (*kubernetes.Clientset, *AuthFriendlyErr) {
	var (
		authFriendlyErr = &AuthFriendlyErr{}
		clientset       = &kubernetes.Clientset{}
	)

	serverHashPath := homeDir() + "/.kubestorm/" + serverName
	if isFileExist(serverHashPath) {
		hashID, err := ioutil.ReadFile(serverHashPath)
		catchError(err)
		kubeconfigPath := homeDir() + "/.kubestorm/" + string(hashID)
		kubeconfig, err := clientcmd.LoadFromFile(kubeconfigPath)
		catchError(err)
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags(kubeconfig.Clusters[serverName].Server, kubeconfigPath)
		catchError(err)

		// create the clientset
		clientset, err = kubernetes.NewForConfig(config)
		catchError(err)
	} else {
		authFriendlyErr.Message = fmt.Sprintf("%v cluster config is not found", serverName)
		authFriendlyErr.HTTPStatus = "404"
	}

	return clientset, authFriendlyErr
}

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
