package discovery

import (
	"fmt"
	"log"
	"os"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

type ServiceDiscovery struct {
	client *consulapi.Client
}

func NewServiceDiscovery() (*ServiceDiscovery, error) {
	config := consulapi.DefaultConfig()
	config.Address = fmt.Sprintf("%s:8500", os.Getenv("CONSUL_HOST"))

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %v", err)
	}

	return &ServiceDiscovery{client: client}, nil
}

func (sd *ServiceDiscovery) RegisterService(serviceName string, servicePort int) error {
	registration := &consulapi.AgentServiceRegistration{
		ID:   fmt.Sprintf("%s-%d", serviceName, servicePort),
		Name: serviceName,
		Port: servicePort,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", os.Getenv("SERVICE_HOST"), servicePort),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	err := sd.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	log.Printf("Service %s registered successfully", serviceName)
	return nil
}

func (sd *ServiceDiscovery) DeregisterService(serviceID string) error {
	err := sd.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %v", err)
	}

	log.Printf("Service %s deregistered successfully", serviceID)
	return nil
}

func (sd *ServiceDiscovery) GetService(serviceName string) (*consulapi.AgentService, error) {
	services, err := sd.client.Agent().Services()
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	for _, service := range services {
		if service.Service == serviceName {
			return service, nil
		}
	}

	return nil, fmt.Errorf("service %s not found", serviceName)
}

func (sd *ServiceDiscovery) GetServiceURL(serviceName string) (string, error) {
	service, err := sd.GetService(serviceName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

func (sd *ServiceDiscovery) GetServicePort(serviceName string) (int, error) {
	service, err := sd.GetService(serviceName)
	if err != nil {
		return 0, err
	}

	return service.Port, nil
}

func GetServicePort() int {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		return 8080 // default port
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 8080
	}
	return port
}
