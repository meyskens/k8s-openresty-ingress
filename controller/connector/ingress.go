package connector

import (
	"fmt"
	"time"

	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetIngresses gets all ingresses from all namespaces
func (c *Client) GetIngresses() ([]extensions_v1beta1.Ingress, error) {
	list, err := c.clientset.ExtensionsV1beta1().Ingresses("").List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list.Items, nil
}

// WatchIngressForChanges watches changes on ingrress inside Kubernetes
func (c *Client) WatchIngressForChanges() (chan bool, error) {
	changeChan := make(chan bool)

	go func() {
		for {
			err := c.watcher(changeChan, c.clientset.ExtensionsV1beta1().Ingresses(""))
			fmt.Println(err)
			time.Sleep(300 * time.Millisecond) // backoff
		}
	}()

	return changeChan, nil
}
