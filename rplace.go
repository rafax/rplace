package rplace

import (
	"fmt"
	"image"
	"image/color"
	"sync"
	"time"
)

var (
	Black = (color.NRGBA{R: 0, G: 0, B: 0, A: 255})
	White = (color.NRGBA{R: 255, G: 255, B: 255, A: 255})
)

var cooldownPeriod = 5 * time.Minute

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
	userWrites map[string]int64 // contains index in pixels
	seq        uint64
	lock       sync.Mutex
	seqLock    sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	// make 1000 a const
	pixels := make([]PixelWrite, 1000*1000)
	for i := int16(0); i < 1000; i++ {
		for j := int16(0); j < 1000; j++ {
			pixels[i*1000+j] = PixelWrite{Pixel: Pixel{x: i, y: j, c: White}}
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
	s.pixels[p.y*1000+p.y] = pw
	return &pw, nil
}

func (s *MemoryStore) Get(x, y int) *PixelWrite {
	return &s.pixels[y*1000+x]
}
func (s *MemoryStore) GetImage() image.Image {
	// TODO: add a helper that takes an array of colors and return an image.
	return nil
}
