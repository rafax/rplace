package img

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"
)

const cooldownPeriod = 5 * time.Minute

const Size = 1000

type Pixel struct {
	X, Y int
	C    color.NRGBA
}

type PixelWrite struct {
	Sequence  uint64
	UserID    string
	WrittenAt time.Time
	Pixel
}

type ImageStore interface {
	Set(p Pixel, userID string) PixelWrite
	Get(x, y int) PixelWrite
	GetImage() image.Image
}

type UserStore interface {
	NextWriteTime(userID string) time.Time
	LastWrite(userID string) PixelWrite
}

var (
	Black     = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	White     = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	LightGray = color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	DarkGray  = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
	Colors    = []color.NRGBA{Black, White, LightGray, DarkGray}
)

type MemoryStore struct {
	pixels     []PixelWrite
	userWrites map[string]int64 // contains index of write in pixels
	seq        uint64
	lock       sync.Mutex
	seqLock    sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	pixels := make([]PixelWrite, Size*Size)
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			pixels[i*Size+j] = PixelWrite{Pixel: Pixel{X: i, Y: j, C: White}}
		}
	}
	return &MemoryStore{pixels: pixels, userWrites: map[string]int64{}}
}

func (s *MemoryStore) Set(p Pixel, userID string) (*PixelWrite, error) {
	now := time.Now()
	if ind, ok := s.userWrites[userID]; ok {
		w := s.pixels[ind]
		cooldownEndsAt := w.WrittenAt.Add(cooldownPeriod)
		if time.Now().After(cooldownEndsAt) {
			return nil, fmt.Errorf("Cooldown period has not passed, wait %s", cooldownEndsAt.Sub(now))
		}
	}
	s.seqLock.Lock()
	s.seq++
	seq := s.seq
	s.seqLock.Unlock()
	pw := PixelWrite{Pixel: p, Sequence: seq, UserID: userID, WrittenAt: now}
	s.pixels[p.Y*Size+p.X] = pw
	return &pw, nil
}

func (s *MemoryStore) Get(x, y int) *PixelWrite {
	return &s.pixels[y*Size+x]
}
func (s *MemoryStore) GetImage() image.Image {
	// TODO: add a helper that takes an array of colors and return an image.
	img := image.NewNRGBA(image.Rect(0, Size, 0, Size))
	for i := 0; i < Size; i++ {
		for j := 0; j < Size; j++ {
			img.SetNRGBA(i, j, s.pixels[j*Size+i].Pixel.C)
		}
	}
	return img
}
