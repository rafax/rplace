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
)

func main() {
	var wg sync.WaitGroup
	total := 1000 * 1000
	workers := 100
	wg.Add(total)
	sem := make(chan int, workers)
	cnt := int64(total)
	log.Println("Allocating")
	for _, i := range rand.Perm(total) {
		go setPixel(i, sem, &cnt, &wg)
	}
	log.Println("Allocated, waiting")
	wg.Wait()
}

func setPixel(i int, sem chan int, cnt *int64, wg *sync.WaitGroup) {
	c := &http.Client{}
	defer wg.Done()
	sem <- 1
	r, err := c.Get(fmt.Sprintf("http://127.0.0.1:8080/write?x=%v&y=%v&c=%v&u=%v", i%1000, i/1000, rand.Int()%4, "user"+strconv.Itoa(rand.Int()%100000)))
	<-sem
	if err != nil {
		log.Printf("Error %v", err)
		return
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
	atomic.AddInt64(cnt, -1)
	if *cnt%10000 == 0 {
		log.Println(*cnt)
	}
}
