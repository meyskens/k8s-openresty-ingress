package main

import (
	"log"
	"sync"
	"time"

	"github.com/meyskens/k8s-openresty-ingress/controller/connector"
)

var changes bool
var changesMutex = sync.Mutex{}

func watchChanges(client *connector.Client) {
	ingressWatch, err := client.WatchIngressForChanges()
	if err != nil {
		panic(err)
	}
	servicesWatch, err := client.WatchServicesForChanges()
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-ingressWatch:
			log.Println("Ingress update: reloading config...")
			changesMutex.Lock()
			changes = true
			changesMutex.Unlock()
			break
		case <-servicesWatch:
			log.Println("Service update: reloading config...")
			changesMutex.Lock()
			changes = true
			changesMutex.Unlock()
			break
		}
	}
}

func runReloadOnChange(client *connector.Client) {
	for {
		time.Sleep(time.Second)
		changesMutex.Lock()
		if changes {
			runAndRetry(reload, client)
		}
		changesMutex.Unlock()
	}
}
