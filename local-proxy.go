package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type LocalProxy struct{}

func NewLocalProxy() Proxier {
	return &LocalProxy{}
}

func (p *LocalProxy) ExtractService(path string) (string, error) {
	// The passed path should be /api/url/*
	split := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(split) <= 1 {
		return "", fmt.Errorf("failed to parse target service from path: %s", path)
	}

	if split[0] != "api" {
		return "", fmt.Errorf("failed to parse target service from path: %s", path)
	}

	url := split[1]
	if url == "" {
		return "", fmt.Errorf("failed to parse target  from path: %s", path)
	}

	return fmt.Sprintf(
		"http://%s/%s",
		url, strings.Join(split[2:], "/"),
	), nil
}

func (p *LocalProxy) Proxy(ctx *gin.Context) {
	proxy(p, ctx)
}
