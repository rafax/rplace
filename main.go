package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/rafax/rplace/img"
)

var store = img.NewMemoryStore()

func main() {

	http.HandleFunc("/write", write)
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
	store.Set(img.Pixel{X: x, Y: y, C: img.Colors[c]}, "user"+strconv.Itoa(rand.Int()))
}
