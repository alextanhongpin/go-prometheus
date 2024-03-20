package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	workers := 1
	count := 1000

	ch := make(chan struct{}, workers)

	var wg sync.WaitGroup
	wg.Add(count)

	var success, failure int

	for i := 0; i < count; i++ {
		ch <- struct{}{}

		sleep := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(sleep)

		go func() {
			defer wg.Done()
			defer func() {
				<-ch
			}()

			resp, err := http.Get("http://localhost:8000/")
			if err != nil {
				failure++
				log.Println(err)
				return
			}
			defer resp.Body.Close()
			_, _ = ioutil.ReadAll(resp.Body)

			log.Println(i)
			success++
		}()
	}

	wg.Wait()

	fmt.Println("success:", success)
	fmt.Println("failure:", failure)
}
