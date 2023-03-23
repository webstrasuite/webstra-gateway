package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	router  *gin.Engine
	server  *http.Server
	gateway *Gateway
}

func NewRouter(port, serviceNamespace string) *Router {
	// Initialise router
	router := gin.New()

	// Disable logging for health check endpoint
	logger := gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: []string{"/api/health"}})

	// Use the logger and recovery middleware
	router.Use(logger, gin.Recovery())

	srv := &http.Server{
		Addr:    port,
		Handler: router,
	}

	return &Router{
		server:  srv,
		router:  router,
		gateway: NewGateway(serviceNamespace),
	}
}

func (r *Router) RegisterRoutes() {
	// Health check endpoint for k8s liveness/readiness
	r.router.GET("/health", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
	})

	// Route any other requests through the reverse proxy / gateway
	r.router.Any("/api/*path", r.gateway.proxy)
}

func (r *Router) Start() {
	go func() {
		// service connections
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}

	log.Println("Server exiting")
}
