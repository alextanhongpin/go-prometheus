package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	workers := 10
	var wg sync.WaitGroup
	wg.Add(2)

	releases := []string{"canary", "stable"}
	for _, release := range releases {
		{
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8000/", nil)
			if err != nil {
				panic(err)
			}
			req.Header.Set("x-release-header", release)

			go func() {
				defer wg.Done()

				concurrentGet(ctx, req, workers)
			}()
		}
	}

	wg.Wait()
	fmt.Println("terminated")
}

func concurrentGet(ctx context.Context, req *http.Request, workers int) {
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {

		go func() {
			defer wg.Done()

			get(ctx, req)
		}()
	}

	wg.Wait()
}

func get(ctx context.Context, req *http.Request) {
	n := rand.Intn(5) + 1
	for {
		select {
		case <-ctx.Done():
			return
		default:
			sleep := time.Duration(n*rand.Intn(100)) * time.Millisecond
			time.Sleep(sleep)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			_, _ = ioutil.ReadAll(resp.Body)
		}
	}
}
