package rplace

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"

	"github.com/rafax/rplace/img"
)

const cooldownPeriod = 5 * time.Minute
const size = 1000

type Pixel struct {
	x, y int16
	c    color.NRGBA
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

type MemoryStore struct {
	pixels     []PixelWrite
	userWrites map[string]int64 // contains index of write in pixels
	seq        uint64
	lock       sync.Mutex
	seqLock    sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	pixels := make([]PixelWrite, size*size)
	for i := int16(0); i < size; i++ {
		for j := int16(0); j < size; j++ {
			pixels[i*size+j] = PixelWrite{Pixel: Pixel{x: i, y: j, c: img.White}}
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
	s.pixels[p.y*size+p.y] = pw
	return &pw, nil
}

func (s *MemoryStore) Get(x, y int) *PixelWrite {
	return &s.pixels[y*size+x]
}
func (s *MemoryStore) GetImage() image.Image {
	// TODO: add a helper that takes an array of colors and return an image.
	return nil
}

func toImage([]color.NRGBA) image.Image {
	return nil
}
