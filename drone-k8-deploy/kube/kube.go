package kube

import (
	"encoding/base64"
	"fmt"

	//"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Client : default kubernetes client
var defaultClient *kubernetes.Clientset

func GetDefaultClient() *kubernetes.Clientset {
	if defaultClient == nil {
		logrus.Fatal("Default Client was not inited, please call kube.Init to init defalt client")
	}
	return defaultClient
}

func InitOutClusterClient(kubeconfigFilePath string) error {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigFilePath)
	if err != nil {
		panic(err.Error())
	}
	defaultClient, err = kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}
	return nil
}

func InitDefaultClient(serverFullAddr string, certOfNamespace string, tokenOfNamespace string) (err error) {
	cert, err := base64.StdEncoding.DecodeString(certOfNamespace)
	if err != nil {
		return fmt.Errorf("failed to decode cert string")
	}

	config := clientcmdapi.NewConfig()
	config.Clusters["drone"] = &clientcmdapi.Cluster{
		Server:                   serverFullAddr,
		CertificateAuthorityData: cert,
		//InsecureSkipTLSVerify:    true,
	}
	config.AuthInfos["drone"] = &clientcmdapi.AuthInfo{
		Token: tokenOfNamespace,
	}

	config.Contexts["drone"] = &clientcmdapi.Context{
		Cluster:  "drone",
		AuthInfo: "drone",
	}

	config.CurrentContext = "drone"
	clientBuilder := clientcmd.NewNonInteractiveClientConfig(*config, "drone", &clientcmd.ConfigOverrides{}, nil)

	restCfg, err := clientBuilder.ClientConfig()
	//fmt.Println(actualCfg)
	if err != nil {
		return err
	}
	defaultClient, err = kubernetes.NewForConfig(restCfg)

	if err != nil {
		return err
	}
	return nil
}
