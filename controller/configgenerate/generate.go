package configgenerate

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"

	"log"

	core_v1 "k8s.io/api/core/v1"
	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"
)

// ConfigFileValues contains the values for one config file
type ConfigFileValues struct {
	Name      string
	Domain    string
	AllowHTTP bool
	Values    []ConfigValues
}

// ConfigValues contains the values for one path rule
type ConfigValues struct {
	Path string
	Host string
	Port int32
}

// GenerateConfigFileValuesFromIngresses gives back the ConfigFileValues for an ingress slice
func GenerateConfigFileValuesFromIngresses(ingresses []extensions_v1beta1.Ingress, serviceMap map[string]core_v1.Service) []ConfigFileValues {
	files := []ConfigFileValues{}
	for _, ingress := range ingresses {
		ingressName := fmt.Sprintf("%s-%s", ingress.GetObjectMeta().GetNamespace(), ingress.GetObjectMeta().GetName())
		allowHTTP := false

		allowHTTPValue, exists := ingress.Annotations["kubernetes.io/ingress.allow-http"]
		if exists {
			allowHTTP = allowHTTPValue == "true"
		}

		values := []ConfigValues{}
		for _, rule := range ingress.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				if path.Backend.ServicePort.Type == intstr.String {
					log.Println("String port values are not yet supported")
					continue
				}
				service, ok := serviceMap[fmt.Sprintf("%s-%s", ingress.GetObjectMeta().GetNamespace(), path.Backend.ServiceName)]
				if !ok {
					log.Printf("Service %s not found in namespace %s\n", path.Backend.ServiceName, ingress.GetObjectMeta().GetNamespace())
				}
				values = append(values, ConfigValues{
					Path: path.Path,
					Host: service.Spec.ClusterIP,
					Port: path.Backend.ServicePort.IntVal,
				})
			}

			files = append(files, ConfigFileValues{
				Name:      fmt.Sprintf("%s-%s", ingressName, rule.Host),
				Domain:    rule.Host,
				Values:    values,
				AllowHTTP: allowHTTP,
			})
		}
	}
	return files
}
