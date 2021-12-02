package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"

	"github.com/fogleman/gg"
)

type edit_req struct {
	BgImgPath string
	FontPath  string
	FontSize  float64
	Text      string
}

type icondir struct {
	reserved  uint16
	imageType uint16
	numImages uint16
}

type icondirentry struct {
	imageWidth   uint8
	imageHeight  uint8
	numColors    uint8
	reserved     uint8
	colorPlanes  uint16
	bitsPerPixel uint16
	sizeInBytes  uint32
	offset       uint32
}

func txt_on_img(request edit_req) (image.Image, error) {
	bgImage, err := gg.LoadImage(request.BgImgPath)
	if err != nil {
		return nil, err
	}
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)

	if err := dc.LoadFontFace(request.FontPath, request.FontSize); err != nil {
		return nil, err
	}

	x := float64(imgWidth / 2)
	y := float64((imgHeight / 2))
	maxWidth := float64(imgWidth)
	dc.SetColor(color.White)
	dc.DrawStringWrapped(request.Text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignCenter)

	return dc.Image(), nil
}

func img_to_nrgba(im image.Image) *image.NRGBA {
	dst := image.NewNRGBA(im.Bounds())
	draw.Draw(dst, im.Bounds(), im, im.Bounds().Min, draw.Src)
	return dst
}

func newIcondir() icondir {
	var id icondir
	id.imageType = 1
	id.numImages = 1
	return id
}

func newIcondirentry() icondirentry {
	var ide icondirentry
	ide.colorPlanes = 1
	ide.bitsPerPixel = 32
	ide.offset = 22 //6 icondir + 16 icondirentry, next image will be this image size + 16 icondirentry, etc
	return ide
}

func convert_to_ico(w io.Writer, im image.Image) error {
	b := im.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, im, b.Min, draw.Src)

	id := newIcondir()
	ide := newIcondirentry()

	pngbb := new(bytes.Buffer)
	pngwriter := bufio.NewWriter(pngbb)
	png.Encode(pngwriter, m)
	pngwriter.Flush()
	ide.sizeInBytes = uint32(len(pngbb.Bytes()))

	bounds := m.Bounds()
	ide.imageWidth = uint8(bounds.Dx())
	ide.imageHeight = uint8(bounds.Dy())
	bb := new(bytes.Buffer)

	var e error
	binary.Write(bb, binary.LittleEndian, id)
	binary.Write(bb, binary.LittleEndian, ide)

	w.Write(bb.Bytes())
	w.Write(pngbb.Bytes())

	return e
}
