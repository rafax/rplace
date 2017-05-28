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
	var wg sync.WaitGroup
	total := 1000 * 1000
	workers := 100
	wg.Add(total)
	cnt := int64(total)
	for i := 0; i < workers; i++ {
		seed := i
		go func(worker int) {
			for req := 0; req < total/workers; req++ {
				ind := worker*total/workers + req
				c := &http.Client{}
				seed++
				for {
					r, err := c.Get(fmt.Sprintf("http://127.0.0.1:8080/write?x=%v&y=%v&c=%v&u=%v", ind%1000, ind/1000, rand.Int()%4, "user"+strconv.Itoa(rand.Int()%100)))
					if err == nil {
						io.Copy(ioutil.Discard, r.Body)
						r.Body.Close()
						break
					}
					c = &http.Client{}
					log.Printf("Error %v")
					time.Sleep(time.Second)
				}
				atomic.AddInt64(&cnt, -1)
				if cnt%1000 == 0 {
					fmt.Println(cnt)
				}
				wg.Done()
			}
		}(i)
	}
	wg.Wait()
}
