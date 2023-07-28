package main

import (
	"fmt"

	"github.com/fir1/port/config"
	http_rest "github.com/fir1/port/http"
	port "github.com/fir1/port/internal/port"

	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/fx"
)

func main() {
	err := fx.New(
		fx.Options(
			config.FxProvide,
			port.FxProvide,
			http_rest.FxProvide,
		),
		fx.Invoke(run),
	).Err()
	if err != nil {
		log.Panic(err)
	}
}

func run(restServer *http_rest.Service) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(interrupt)

	var wg sync.WaitGroup

	wg.Add(1)

	stop := make(chan struct{})
	errChan := make(chan error)

	go func() {
		defer wg.Done()

		err := restServer.Run(stop)
		if err != nil {
			errChan <- err
		}
	}()

	// Wait signal or error from services
	select {
	case <-interrupt:
	case err := <-errChan:
		return fmt.Errorf("webhook rest api http is down (error: %w)", err)
	}

	stop <- struct{}{}
	wg.Wait()

	return nil
}
