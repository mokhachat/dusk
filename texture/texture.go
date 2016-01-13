package texture

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
    
    "os"
    
    _ "github.com/ftrvxmtrx/tga"
    "github.com/go-gl/gl/v3.3-core/gl"
)

func Create(rgba image.RGBA) (uint32, error) {

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
    gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1);
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR/*gl.LINEAR*/)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT/*CLAMP_TO_EDGE*/)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
    gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

func Load(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, err
	}
    defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	
    rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
    
    /*topng, _ := os.Create(file + ".png")
    defer topng.Close()

    png.Encode(topng, rgba)
    */
	return Create(*rgba)
}
