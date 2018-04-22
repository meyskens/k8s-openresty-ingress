package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/meyskens/k8s-openresty-ingress/controller/configgenerate"
	"github.com/meyskens/k8s-openresty-ingress/controller/connector"
)

type retryableFunc func(*connector.Client) error

func main() {
	log.Println("Starting OpenResty Ingress Controller...")

	client, err := connector.NewClient()
	if err != nil {
		panic(err)
	}
	ingress, err := client.GetIngresses()
	if err != nil {
		panic(err)
	}
	services, err := client.GetServiceMap()
	if err != nil {
		panic(err)
	}

	conf := configgenerate.GenerateConfigFileValuesFromIngresses(ingress, services)
	configgenerate.WriteFilesFromTemplate(conf, getTemplatePath(), getIngressPath())

	log.Println("Starting NGINX")
	startNginx()

	go runReloadOnChange(client)
	watchChanges(client)
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

func runAndRetry(fn retryableFunc, client *connector.Client) {
	for {
		err := fn(client)
		if err == nil {
			break
		}
		log.Println("Needs to retry because of", err)
		time.Sleep(time.Second) // sleep before retry
	}
}

func reload(client *connector.Client) error {
	ingress, err := client.GetIngresses()
	if err != nil {
		return err
	}
	services, err := client.GetServiceMap()
	if err != nil {
		return err
	}

	conf := configgenerate.GenerateConfigFileValuesFromIngresses(ingress, services)
	err = configgenerate.WriteFilesFromTemplate(conf, getTemplatePath(), getIngressPath())
	if err != nil {
		return err
	}

	nginx := exec.Command("nginx", "-s", "reload")
	nginx.Stderr = os.Stderr
	nginx.Stdout = os.Stdout
	nginx.Run()

	return nil
}
