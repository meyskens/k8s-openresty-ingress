package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/meyskens/k8s-openresty-ingress/controller/configgenerate"
	"github.com/meyskens/k8s-openresty-ingress/controller/connector"
)

func main() {
	log.Println("Starting OpenResty Ingress Controller...")

	client, _ := connector.NewClient()
	ingress, _ := client.GetIngresses()
	services, _ := client.GetServiceMap()

	conf := configgenerate.GenerateConfigFileValuesFromIngresses(ingress, services)
	configgenerate.WriteFilesFromTemplate(conf, getTemplatePath(), getIngressPath())

	log.Println("Starting NGINX")
	startNginx()
	ingressWatch, err := client.WatchIngressForChanges()
	fmt.Println(err)
	servicesWatch, err := client.WatchServicesForChanges()
	fmt.Println(err)
	for {
		select {
		case <-ingressWatch:
			fmt.Println("Ingress update: reloading config...")
			reload(client)
			break
		case <-servicesWatch:
			fmt.Println("Service update: reloading config...")
			reload(client)
			break
		}
	}
}

func startNginx() *os.Process {
	nginx := exec.Command("nginx", "-c", "/etc/nginx/nginx.conf")
	nginx.Stderr = os.Stderr
	nginx.Stdout = os.Stdout
	nginx.Start()

	for {
		_, err := os.OpenFile("/run/nginx.pid", 'r', 0755)
		if err == nil {
			break // nginx is running
		}
		time.Sleep(100 * time.Millisecond)
		log.Println("Waiting on nginx.pid")
	}
	return nginx.Process
}

func getTemplatePath() string {
	envPath := os.Getenv("OPENRESTY_TEMPLATEPATH")
	if envPath != "" {
		return envPath
	}
	return "../template/ingress.tpl" // Dev fallback
}

func getIngressPath() string {
	envPath := os.Getenv("OPENRESTY_INGRESSATH")
	if envPath != "" {
		return envPath
	}
	return "../debug-out" // Dev fallback
}

func reload(client *connector.Client) {
	ingress, _ := client.GetIngresses()
	services, _ := client.GetServiceMap()

	conf := configgenerate.GenerateConfigFileValuesFromIngresses(ingress, services)
	configgenerate.WriteFilesFromTemplate(conf, getTemplatePath(), getIngressPath())

	nginx := exec.Command("nginx", "-s", "reload")
	nginx.Stderr = os.Stderr
	nginx.Stdout = os.Stdout
	nginx.Run()
}
