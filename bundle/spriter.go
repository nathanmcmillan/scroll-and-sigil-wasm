package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type sprite struct {
	x     int
	y     int
	image *image.RGBA
}

func upscaleSprite(from, to string) {
	dir, err := ioutil.ReadDir(from)
	if err != nil {
		panic(err)
	}
	for _, info := range dir {
		name := info.Name()
		extension := filepath.Ext(name)
		if extension != ".png" {
			continue
		}
		fmt.Println("upscale", name)
		original := getPNG(filepath.Join(from, name))
		width := original.Rect.Size().X
		height := original.Rect.Size().Y

		data := make([]int, width*height)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				sample := original.RGBAAt(x, y)
				data[x+y*width] = (int(sample.A) << 24) + (int(sample.B) << 16) + (int(sample.G) << 8) + int(sample.R)
			}
		}

		const scale = 2
		newWidth := width * scale
		newHeight := height * scale
		upscale := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

		xbr := superXBR(data, width, height)
		for x := 0; x < newWidth; x++ {
			for y := 0; y < newHeight; y++ {
				sample := xbr[x+y*newWidth]
				r := uint8(sample & 0xff)
				g := uint8((sample >> 8) & 0xff)
				b := uint8((sample >> 16) & 0xff)
				a := uint8((sample >> 24) & 0xff)
				upscale.Set(x, y, color.RGBA{r, g, b, a})
			}
		}

		writePNG(filepath.Join(to, name), upscale)
	}
}

func spritePacker(spriteName, from, to string) string {
	dir, err := ioutil.ReadDir(from)
	if err != nil {
		panic(err)
	}

	const limit = 1024

	var data strings.Builder
	x := 0
	y := 0
	atlasWidth := 0
	atlasHeight := 0

	images := make([]*sprite, 0)

	for _, info := range dir {
		name := info.Name()
		extension := filepath.Ext(name)
		base := strings.TrimSuffix(name, extension)
		if extension != ".png" {
			continue
		}
		fmt.Println(name)
		rgba := getPNG(filepath.Join(from, info.Name()))
		width := rgba.Rect.Size().X
		height := rgba.Rect.Size().Y

		if x+width+1 > limit {
			atlasWidth = x - 1
			x = 0
			y = atlasHeight + 1
		}

		save := &sprite{x: x, y: y, image: rgba}
		images = append(images, save)

		data.WriteString(base)
		data.WriteString("[")
		data.WriteString(strconv.Itoa(x))
		data.WriteString(",")
		data.WriteString(strconv.Itoa(y))
		data.WriteString(",")
		data.WriteString(strconv.Itoa(width))
		data.WriteString(",")
		data.WriteString(strconv.Itoa(height))
		data.WriteString("],")

		x += width + 1

		if y+height > atlasHeight {
			atlasHeight = y + height
		}
	}

	if x-1 > atlasWidth {
		atlasWidth = x - 1
	}

	power := 1
	for power < atlasWidth {
		power *= 2
	}
	atlasWidth = power

	power = 1
	for power < atlasHeight {
		power *= 2
	}
	atlasHeight = power

	bundle := image.NewRGBA(image.Rect(0, 0, atlasWidth, atlasHeight))

	for i := 0; i < len(images); i++ {
		source := images[i]
		point := image.Point{source.x, source.y}
		bounds := image.Rectangle{point, point.Add(source.image.Bounds().Size())}
		draw.Draw(bundle, bounds, source.image, image.Point{0, 0}, draw.Src)
	}

	writePNG(filepath.Join(to, spriteName+".png"), bundle)

	return data.String()
}

func getPNG(path string) *image.RGBA {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png, err := png.Decode(file)
	if err != nil {
		panic(err)
	}
	rgba := image.NewRGBA(png.Bounds())
	draw.Draw(rgba, rgba.Bounds(), png, image.Point{0, 0}, draw.Src)
	return rgba
}

func writePNG(path string, image *image.RGBA) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = png.Encode(file, image)
	if err != nil {
		panic(err)
	}
}
