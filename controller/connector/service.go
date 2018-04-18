package connector

import (
	"fmt"
	"time"

	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	changeChan := make(chan bool)

	go func() {
		for {
			err := c.watcher(changeChan, c.clientset.CoreV1().Services(""))
			fmt.Println(err)
			time.Sleep(300 * time.Millisecond) // backoff
		}
	}()

	return changeChan, nil
}
