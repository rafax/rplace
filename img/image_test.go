package img

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func Test_toPng(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 1000, 1000))
	for i := 0; i < 1000*1000; i++ {
		var c color.NRGBA
		switch i % 4 {
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
