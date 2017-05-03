package img

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/rafax/rplace"
)

func Test_toPng(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1000, 1000))
	for i := 0; i < 1000*1000; i++ {
		var c color.NRGBA
		if i%2 == 0 {
			c = rplace.White
		} else {
			c = rplace.Black
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
