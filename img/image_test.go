package img

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func Test_GetAfterSet(t *testing.T) {
	s := NewMemoryStore()
	p := Pixel{X: rand.Int() % 1000, Y: rand.Int() % 1000, C: LightGray}
	s.Set(p, strconv.Itoa(rand.Int()))
	written := s.Get(p.X, p.Y)
	if written.C != p.C {
		t.Errorf("Expected Get() to return what was Set() but %v != %v", written.C, p.C)
	}
	i := s.GetImage()
	if i.At(p.X, p.Y) != p.C {
		t.Errorf("Expected GetImage() to return what was Set() but %v != %v", written.C, p.C)
	}
}

func Test_toPng(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1000, 1000))
	for i := 0; i < 1000*1000; i++ {
		var c color.NRGBA
		switch rand.Int() % 4 {
		case 0:
			c = White
		case 1:
			c = LightGray
		case 2:
			c = DarkGray
		case 3:
			c = Black
		}
		img.Set(i%1000, i/1000, c)
	}

	f, err := os.Create("image.png")
	if err != nil {
		t.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

}

func Test_1M_parallel_writes(t *testing.T) {
	s := NewMemoryStore()
	uids := []string{}
	uidCount := 1000
	for i := 0; i < uidCount; i++ {
		uids = append(uids, strconv.Itoa(rand.Int()))
	}
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(10)
	for r := 0; r < 10; r++ {
		go func() {
			for i := 0; i < 1000*100; i++ {
				var c color.NRGBA
				switch rand.Int() % 4 {
				case 0:
					c = White
				case 1:
					c = LightGray
				case 2:
					c = DarkGray
				case 3:
					c = Black
				}
				s.Set(Pixel{X: rand.Int() % 1000, Y: rand.Int() % 1000, C: c}, uids[i%uidCount])
			}
			wg.Done()
		}()
	}
	wg.Wait()
	secs := time.Now().Sub(start).Seconds()
	t.Logf("1m writes in %v s", secs)
	if secs > 10.0 {
		t.Errorf("Expected to take less than 10s, took %v seconds", secs)
	}
}
