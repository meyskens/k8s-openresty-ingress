package main

import (
	"log"
	"os"
	"os/exec"

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
}

func startNginx() {
	nginx := exec.Command("nginx")
	nginx.Stderr = os.Stderr
	nginx.Stdout = os.Stdout
	nginx.Start()
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
