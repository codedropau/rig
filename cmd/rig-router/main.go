package main

import (
	"github.com/alecthomas/kingpin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/codedropau/rig/internal/router"
)

var (
	cliAddr = kingpin.Flag("addr", "Address to respond to requests").Default(":80").String()
	cliRefresh = kingpin.Flag("refresh", "How often to refresh the list of Pods").Default("10s").Duration()

	cliAuthUsername = kingpin.Flag("auth-username", "Username to be used for basic authentication").Envar("RIG_ROUTER_AUTH_USERNAME").Required().String()
	cliAuthPassword = kingpin.Flag("auth-password", "Password to be used for basic authentication").Envar("RIG_ROUTER_AUTH_PASSWORD").Required().String()

	cliMaster    = kingpin.Flag("master", "URL of the Kubernetes master").String()
	cliKubecfg    = kingpin.Flag("kubecfg", "Path to a local Kuberneretes config file").Envar("KUBECONFIG").String()
)

func main() {
	kingpin.Parse()

	config, err := clientcmd.BuildConfigFromFlags(*cliMaster, *cliKubecfg)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	err = router.Run(*cliAddr, clientset, *cliRefresh, *cliAuthUsername, *cliAuthPassword)
	if err != nil {
		panic(err)
	}
}