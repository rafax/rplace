package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	total := 5 * 1000 * 1000
	workers := 100
	wg.Add(total)
	sem := make(chan int, workers)
	cnt := int64(total)
	c := &http.Client{Transport: &http.Transport{
		MaxIdleConnsPerHost: workers,
	}}
	log.Println("Allocating permutation")
	perm := rand.Perm(total)
	log.Println("Allocated permutation")
	for _, i := range perm {
		setPixel(i, sem, &cnt, &wg, c)
	}
	log.Println("Started all clients, waiting")
	wg.Wait()
	end := time.Now()
	log.Printf("%v writes in %v using %v workers, rate %v/s", total, end.Sub(start), workers, float64(total)/end.Sub(start).Seconds())
}

func setPixel(i int, sem chan int, cnt *int64, wg *sync.WaitGroup, c *http.Client) {
	defer wg.Done()
	sem <- 1
	go func() {
		r, err := c.Get(fmt.Sprintf("http://127.0.0.1:8080/write?x=%v&y=%v&c=%v&u=%v", i%1000, (i/1000)%1000, rand.Int()%4, "user"+strconv.Itoa(rand.Int()%100000)))
		<-sem
		if err != nil {
			log.Printf("Error %v", err)
			return
		}
		defer func() {
			io.Copy(ioutil.Discard, r.Body)
			r.Body.Close()
		}()
		atomic.AddInt64(cnt, -1)
		if *cnt%10000 == 0 {
			log.Println(*cnt)
		}
	}()
}
