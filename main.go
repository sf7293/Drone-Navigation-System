package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sf7293/Drone-Navigation-System/config"
	router "github.com/sf7293/Drone-Navigation-System/router"
	"github.com/sf7293/Drone-Navigation-System/service/metrix"
	"github.com/sf7293/Drone-Navigation-System/utils/logger"
)

func main() {
	exitChan := createExitChannel(10 * time.Second)
	wg := &sync.WaitGroup{}

	config.Init()

	jeagerConfig := metrix.JaegerConfigsData{
		ServerName:                  config.App.ServerName,
		JaegerAgentAddress:          config.Tracer.JaegerAgentAddress,
		JaegerReporterFlushInterval: config.Tracer.JaegerReporterFlushInterval,
		Testing:                     !config.Tracer.IsActive,
	}
	metrix.Init(jeagerConfig)

	defer func() {
		// close tracer
		_ = metrix.TracerCloser.Close()
	}()

	err := router.Run(exitChan, wg, config.App.HTTPPort)
	if err != nil {
		logger.ZSLogger.Errorf("HTTP server closed with error: %v", err)
		_ = logger.ZSLogger.Sync()
		os.Exit(1)
	}
	wg.Wait()
}

// createExitChannel creates and returns a channel that will close when an os interrupt is received.
// This is useful for handling graceful shutdown of different parts of the service.
// If an interrupt is received and graceful shutdown times out, or after receiving a second interrupt,
// createExitChannel will forcefully shut the server down.
func createExitChannel(timeout time.Duration) (exit <-chan struct{}) {
	exitChan := make(chan struct{})

	// handle os signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		// wait for an os signal, then close the exit channel to notify the system of the shutdown request
		<-sigChan
		close(exitChan)

		select {
		case <-sigChan:
			// second interrupt is received. terminate
			logger.ZSLogger.Errorf("graceful shutdown failed: killed the server because of user request")
			_ = logger.ZSLogger.Sync()
			os.Exit(1)
		case <-time.After(timeout):
			// timed out. terminate
			logger.ZSLogger.Errorf("graceful shutdown failed: killed the server because of timeout")
			_ = logger.ZSLogger.Sync()
			os.Exit(1)
		}
	}()

	return exitChan
}
