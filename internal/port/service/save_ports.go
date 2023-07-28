package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/fir1/port/internal/port/model"
	"github.com/fir1/port/internal/port/repository"
)

func (s PortService) SavePortsFromFile(ctx context.Context, filePath string, file multipart.File) error {
	// Check if both filePath and file are provided
	if filePath != "" && file != nil {
		return errors.New("both filePath and file cannot be provided at the same time")
	}

	// Check if neither filePath nor file are provided
	if filePath == "" && file == nil {
		return errors.New("either filePath or file must be provided")
	}

	// Create a cancel context and obtain a cancel function
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error)
	var wg sync.WaitGroup

	jsonStream := NewJSONStream()

	// Use a worker pool to handle port processing goroutines
	numWorkers := runtime.NumCPU()
	workerPool := make(chan struct{}, numWorkers)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		for data := range jsonStream.Watch() {
			if data.Error != nil {
				// If there is an error in JSON stream, cancel the context to stop further processing
				cancel()
				errChan <- data.Error
				return
			}

			workerPool <- struct{}{} // Acquire a worker slot from the pool
			wg.Add(1)
			go func(ctx context.Context, wg *sync.WaitGroup, id string, p model.Port) {
				defer wg.Done()
				defer func() {
					<-workerPool
				}() // Release the worker slot when done processing

				// Check if the cancel signal has been received
				select {
				case <-ctx.Done():
					return
				default:
				}

				// Check if the port already exists in the DB
				_, err := s.repository.Get(ctx, id)
				switch {
				case err == nil:
					err = s.repository.Update(ctx, id, p)
					if err != nil {
						errChan <- fmt.Errorf("error updating port with ID %s: %w", id, err)
						return
					}
				case errors.As(err, &repository.ErrObjectNotFound{}):
					// If not found, create a new record in the DB
					err = s.repository.Create(ctx, id, p)
					if err != nil {
						errChan <- fmt.Errorf("error creating port with ID %s: %w", id, err)
						return
					}
				default:
					errChan <- fmt.Errorf("error getting port with ID %s: %w", id, err)
					return
				}
			}(ctx, wg, data.PortCode, data.Port)
		}
	}(&wg)

	var err error
	if filePath != "" {
		if !filepath.IsAbs(filePath) {
			filePath = fmt.Sprintf("%s/%s", s.config.DataDir, filePath)
		}

		file, err = os.Open(filePath)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
	}

	jsonStream.Start(file)

	// Start a goroutine to close the error channel after all goroutines have completed
	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		return err
	}
	return nil
}
