package connector

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	clientset *kubernetes.Clientset
}

// NewClient generates a client with the right configuration
func NewClient() (*Client, error) {
	clientset, err := getInClusterClientset()
	if err != nil {
		//return nil, err
		clientset, err = getLocalClientSet("lab") // TODO: change me
		if err != nil {
			return nil, err
		}
	}
	client := Client{
		clientset: clientset,
	}
	return &client, nil
}

func getInClusterClientset() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func getLocalClientSet(context string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		},
	).ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
