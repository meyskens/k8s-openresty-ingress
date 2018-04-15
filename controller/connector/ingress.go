package connector

import (
	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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
	chageChan := make(chan bool)

	w, err := c.clientset.ExtensionsV1beta1().Ingresses("").Watch(meta_v1.ListOptions{})
	if err != nil {
		return chageChan, err
	}

	go func() {
		for {
			event := <-w.ResultChan()
			if event.Type != watch.Error {
				chageChan <- true
			}
		}
	}()

	return chageChan, nil
}
