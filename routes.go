package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/webstrasuite/webstra-gateway/proxy"
)

type Router struct {
	router *echo.Echo
	server *http.Server
	proxy  proxy.Proxier
}

func NewRouter(port string, proxy proxy.Proxier) *Router {
	// Initialise router
	e := echo.New()

	// Logger config to skip logging of healthcheck calls
	loggerConfig := middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/health"
		},
	}

	// Register custom logger and standard recovery middleware
	e.Use(middleware.Recover(), middleware.LoggerWithConfig(loggerConfig))

	srv := &http.Server{
		Addr:    port,
		Handler: e,
	}

	return &Router{
		server: srv,
		router: e,
		proxy:  proxy,
	}
}

func (r *Router) RegisterRoutes() {
	// Health check endpoint for k8s liveness/readiness
	r.router.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	// Route any other requests through the reverse proxy / gateway
	r.router.Any("/api/*path", r.proxy.Handle)
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
