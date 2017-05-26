package main

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"sync"

	"github.com/rafax/rplace/img"
)

var store = img.NewMemoryStore()
var imageBytes []byte
var lock = sync.RWMutex{}

func main() {
	go updateImage()
	http.HandleFunc("/write", write)
	http.HandleFunc("/img", getImage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func write(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	x, err := strconv.Atoi(r.Form.Get("x"))
	y, err := strconv.Atoi(r.Form.Get("y"))
	c, err := strconv.Atoi(r.Form.Get("c"))
	u := r.Form.Get("u")
	write, err := store.Set(img.Pixel{X: x, Y: y, C: img.Colors[c]}, u)
	if err != nil {
		if _, ok := err.(img.UserRateLimitedError); ok {
			w.WriteHeader(429)
			log.Println("Cooldown served for u")
		}
		w.WriteHeader(500)
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, write)
}

func getImage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	lock.RLock()
	io.Copy(w, bytes.NewBuffer(imageBytes))
	lock.RUnlock()
}

func updateImage() {
	for {
		select {
		case <-time.After(1 * time.Second):
			curr := store.GetImage()
			buf := &bytes.Buffer{}
			png.Encode(buf, curr)
			lock.Lock()
			imageBytes = buf.Bytes()
			lock.Unlock()
		}
	}
}
