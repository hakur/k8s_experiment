package main

import (
	"flag"

	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	provisionerName   = "hakurei.cn/cephfs"
	provisionerIDAnn  = "cephFSProvisionerIdentity"
	cephShareAnn      = "cephShare"
	version           = "0.1.0"
	cephFSRootMntPath = "/mnt/cephFSRoot"
)

// kubeconfig kubernetes配置变量
var kubeconfig *string
var PVBaseDirName = "hakur_cephfs"

// Client 默认的kubernetes client
var Client *kubernetes.Clientset

func init() {
	kubeconfig = flag.String("kubeconfig", "./kubeconfig", "absolute path to the kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	Client, err = kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}
}

func main() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)

	cephFSProvisioner := NewProvisioner()
	pc := controller.NewProvisionController(
		Client,
		provisionerName,
		cephFSProvisioner,
		version,
		//controller.MetricsPort(int32(*metricsPort)),
	)

	pc.Run(wait.NeverStop)
}
