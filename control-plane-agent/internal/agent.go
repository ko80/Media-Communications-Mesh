package internal

import (
	"context"
	proxyapi "control-plane-agent/api/mp-api"
	"control-plane-agent/internal/registry"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	IsProd bool
}

func RunAgent() {
	// DEBUG
	cfg := Config{IsProd: false}

	// Setup logger
	if cfg.IsProd {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	} else {
		logrus.SetLevel(logrus.TraceLevel)
		logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	}

	// Setup logger for http standard output
	log.SetOutput(logrus.StandardLogger().Writer())

	logrus.Info("Mesh Control Plane Agent started")

	proxyAPI := proxyapi.NewAPI(proxyapi.Config{ListenPort: 8100})

	registry.MediaProxyRegistry.Init()
	registry.ConnRegistry.Init()
	registry.MultipointGroupRegistry.Init()
	registry.BridgeRegistry.Init()

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		proxyAPI.Run(ctx)
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := registry.MediaProxyRegistry.Run(ctx)
		if err != nil {
			logrus.Errorf("media proxy registry run err: %v", err)
		}
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := registry.ConnRegistry.Run(ctx)
		if err != nil {
			logrus.Errorf("conn registry run err: %v", err)
		}
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := registry.MultipointGroupRegistry.Run(ctx)
		if err != nil {
			logrus.Errorf("multipoint group registry run err: %v", err)
		}
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := registry.BridgeRegistry.Run(ctx)
		if err != nil {
			logrus.Errorf("bridge registry run err: %v", err)
		}
		cancel()
	}()

	// Start a goroutine for handling shutdown signals
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Setup notification of shutdown signals
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

		select {
		case <-signals:
			logrus.Info("Shutdown signal received")

		case <-ctx.Done():
		}

		cancel()
	}()

	wg.Wait()
	logrus.Info("Mesh Control Plane Agent exited gracefully")
}
