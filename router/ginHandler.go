package router

import (
	"context"
	"fmt"
	core_config "github.com/sf7293/Drone-Navigation-System/config/core"
	"net/http"
	"os"
	"sync"

	"github.com/sf7293/Drone-Navigation-System/router/controller"
	"github.com/sf7293/Drone-Navigation-System/utils/logger"

	"github.com/gin-gonic/gin"
)

func newGinRouter() (r *gin.Engine) {
	DNSController := controller.NewDNSController()

	// init gin
	r = gin.New()
	r.Use(gin.Recovery())

	Env := os.Getenv("env")
	if len(Env) == 0 {
		Env = core_config.EnvDev
	}

	switch Env {
	case core_config.EnvDev:
		gin.SetMode(gin.DebugMode)
		r.Use(gin.Logger())
	case core_config.EnvTest:
		gin.SetMode(gin.TestMode)
		r.Use(gin.Logger())
	case core_config.EnvProd:
		gin.SetMode(gin.ReleaseMode)
	}

	// setup routes
	v1 := r.Group("/v1")
	{
		dns := v1.Group("/dns")
		{
			dns.POST("/location", DNSController.CalculateLocation)
		}
	}

	/*
		monit := r.Group("/monit")
		{
			monit.GET("/liveness", monitController.GetHealth)
		}
	*/

	return r
}

// Run runs HTTP server and listens on port. It will block the current goroutine.
// When exit channel is closed, HTTP server will try to gracefully shutdown.
func Run(exit <-chan struct{}, wg *sync.WaitGroup, port int) (err error) {
	wg.Add(1)
	defer wg.Done()

	r := newGinRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	// setup graceful shutdown
	go func() {
		<-exit
		logger.ZSLogger.Info("Shutting HTTP server down")

		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logger.ZSLogger.Errorf("HTTP server failed to shutdown gracefully: %v", err)
		}
	}()

	return srv.ListenAndServe()
}
