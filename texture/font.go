package texture

import (
	"image"
	"image/draw"
	"io/ioutil"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func Text(fontfile string, size int, text string) (uint32, error) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		return 0, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return 0, err
	}

	// Draw the background and the guidelines.
	fg, bg := image.White, image.Transparent
	imgW, imgH := size*len([]rune(text)), size

	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	// Draw the text.
	dpi := 72
	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    float64(size),
			DPI:     float64(dpi),
			Hinting: font.HintingNone,
		}),
	}
	y := int(math.Ceil(float64(size*dpi) / 72))
	d.Dot = fixed.Point26_6{
		X: (fixed.I(imgW) - d.MeasureString(text)) / 2,
		Y: fixed.I(y),
	}
	d.DrawString(text)

	// TODO: export UV
	return Create(*rgba)
}
