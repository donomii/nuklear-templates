package nktemplates

// From "github.com/cstegel/opengl-samples-golang/colors/gfx"

import (
	"errors"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io"

	"github.com/donomii/glim"

	//"log"
	"os"

	"github.com/go-gl/gl/v3.2-core/gl"
)

type Texture struct {
	Handle  uint32
	target  uint32 // same target as gl.BindTexture(<this param>, ...)
	texUnit uint32 // Texture unit that is currently bound to ex: gl.TEXTURE0
}

var errUnsupportedStride = errors.New("unsupported stride, only 32-bit colors supported")

var errTextureNotBound = errors.New("texture not bound")

func NewTextureFromFile(file string, wrapR, wrapS int32) (*Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	// Decode detexts the type of image as long as its image/<type> is imported
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	return NewGarbageTexture(img, wrapR, wrapS)
}

func NewTextureFromData(data []uint8, width, height int32) (*Texture, error) {
	imgin := glim.ImageToGFormat(int(width), int(height), data)
	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()
	go png.Encode(writer, imgin)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return NewGarbageTexture(img, width, height)
}

func NewGarbageTexture(img image.Image, wrapR, wrapS int32) (*Texture, error) {
	rgba := image.NewRGBA(img.Bounds())
	//log.Printf("Image size: %+v\n", img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 { // TODO-cs: why?
		return nil, errUnsupportedStride
	}

	var handle uint32
	gl.GenTextures(1, &handle)

	target := uint32(gl.TEXTURE_2D)
	internalFmt := int32(gl.SRGB_ALPHA)
	format := uint32(gl.RGBA)
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	pixType := uint32(gl.UNSIGNED_BYTE)
	dataPtr := gl.Ptr(rgba.Pix)

	texture := Texture{
		Handle: handle,
		target: target,
	}

	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	// set the texture wrapping/filtering options (applies to current bound texture obj)
	// TODO-cs
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_R, wrapR)
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_S, wrapS)
	gl.TexParameteri(texture.target, gl.TEXTURE_MIN_FILTER, gl.LINEAR) // minification filter
	gl.TexParameteri(texture.target, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // magnification filter

	gl.TexImage2D(target, 0, internalFmt, width, height, 0, format, pixType, dataPtr)

	gl.GenerateMipmap(texture.Handle)

	return &texture, nil
}

func RawTexture(data []byte, wrapR, wrapS int32, texture *Texture) (*Texture, error) {

	var handle uint32

	target := uint32(gl.TEXTURE_2D)
	internalFmt := int32(gl.RGBA)
	format := uint32(gl.RGBA)
	width := wrapR
	height := wrapS
	pixType := uint32(gl.UNSIGNED_BYTE)
	dataPtr := gl.Ptr(data)

	if texture == nil {
		gl.GenTextures(1, &handle)
		texture = &Texture{
			Handle: handle,
			target: target,
		}
	} else {
		handle = texture.Handle
	}

	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	// set the texture wrapping/filtering options (applies to current bound texture obj)
	// TODO-cs
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(texture.target, gl.TEXTURE_MIN_FILTER, gl.LINEAR) // minification filter
	gl.TexParameteri(texture.target, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // magnification filter

	gl.TexImage2D(target, 0, internalFmt, width, height, 0, format, pixType, dataPtr)
	dataPtr = nil

	gl.GenerateMipmap(texture.Handle)

	return texture, nil
}

func (tex *Texture) Bind(texUnit uint32) {
	gl.ActiveTexture(texUnit)
	gl.BindTexture(tex.target, tex.Handle)
	tex.texUnit = texUnit
}

func (tex *Texture) UnBind() {
	tex.texUnit = 0
	gl.BindTexture(tex.target, 0)
}

func (tex *Texture) SetUniform(uniformLoc int32) error {
	if tex.texUnit == 0 {
		return errTextureNotBound
	}
	gl.Uniform1i(uniformLoc, int32(tex.texUnit-gl.TEXTURE0))
	return nil
}

func loadImageFile(file string) (image.Image, error) {
	infile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Decode automatically figures out the type of immage in the file
	// as long as its image/<type> is imported
	img, _, err := image.Decode(infile)
	return img, err
}
