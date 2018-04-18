package connector

import (
	"fmt"
	"log"

	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// GetServiceMap gives all services in a map to look them up in (namespace)-(service) format
func (c *Client) GetServiceMap() (map[string]core_v1.Service, error) {
	servicesList, err := c.clientset.CoreV1().Services("").List(meta_v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	serviceMap := map[string]core_v1.Service{}
	for _, service := range servicesList.Items {
		serviceMap[fmt.Sprintf("%s-%s", service.GetObjectMeta().GetNamespace(), service.GetObjectMeta().GetName())] = service
	}

	return serviceMap, nil
}

// WatchServicesForChanges watches changes on services inside Kubernetes
func (c *Client) WatchServicesForChanges() (chan bool, error) {
	chageChan := make(chan bool)

	w, err := c.clientset.CoreV1().Services("").Watch(meta_v1.ListOptions{})
	if err != nil {
		return chageChan, err
	}

	go func() {
		for {
			event := <-w.ResultChan()
			if event.Type == watch.Error {
				log.Println(event.Object)
				break
			}
			chageChan <- true
		}
		w.Stop()
		// TO DO: restart watcher
	}()

	return chageChan, nil
}
