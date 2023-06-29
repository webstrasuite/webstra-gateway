package proxy

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/webstrasuite/webstra-gateway/pb"
)

type KubernetesProxy struct {
	serviceNamespace string
}

func NewKubernetes(serviceNamespace string) Proxier {
	return &KubernetesProxy{
		serviceNamespace: serviceNamespace,
	}
}

func (p *KubernetesProxy) Handler(authClient pb.AuthServiceClient) echo.HandlerFunc {
	return proxy(authClient, p)
}

func (p *KubernetesProxy) ExtractService(path string) (string, error) {
	// The passed path should be /api/{serviceName}/*
	split := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(split) <= 1 {
		return "", fmt.Errorf("failed to parse target service from path: %s", path)
	}

	if split[0] != "api" {
		return "", fmt.Errorf("failed to parse target service from path: %s", path)
	}

	serviceHost := fmt.Sprintf("svc-%s", split[1])
	if serviceHost == "" {
		return "", fmt.Errorf("failed to parse target  from path: %s", path)
	}

	// Return the interal k8s address for the found service
	return fmt.Sprintf(
		"http://%s.%s.svc.cluster.local/%s",
		serviceHost, p.serviceNamespace, split[2:]), nil
}
