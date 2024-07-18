package photoLayouts

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

type Photo struct {
	File            string
	LayoutFile      string
	ContainerWidth  float64
	ContainerHeight float64
	PhotoWidth      float64
	PhotoHeight     float64
	Color           string
	Dpi             float64
}

func Layout(p *Photo) (string, error) {
	// convert millimeters (mm) to pixels (px)
	p.ContainerHeight = ConverMmToPx(p.Dpi, p.ContainerHeight)
	p.ContainerWidth = ConverMmToPx(p.Dpi, p.ContainerWidth)
	p.PhotoWidth = ConverMmToPx(p.Dpi, p.PhotoWidth)
	p.PhotoHeight = ConverMmToPx(p.Dpi, p.PhotoHeight)

	// If the photo size exceeds the container size, set the photo to the container size.
	if p.PhotoWidth >= p.ContainerWidth && p.PhotoHeight >= p.ContainerHeight {
		p.PhotoWidth, p.PhotoHeight = p.ContainerWidth, p.ContainerHeight
	}

	// Determine if the container needs to be rotated.
	isRotate := (p.ContainerWidth/p.PhotoHeight)*(p.ContainerHeight/p.PhotoWidth) > (p.ContainerWidth/p.PhotoWidth)*(p.ContainerHeight/p.PhotoHeight)
	// If rotation is needed, swap the container's width and height.
	if isRotate {
		p.ContainerHeight, p.ContainerWidth = p.ContainerWidth, p.ContainerHeight
	}
	// Calculate how many photos can fit in each direction of the container.
	floorX := math.Floor(p.ContainerWidth / p.PhotoWidth)
	floorY := math.Floor(p.ContainerHeight / p.PhotoHeight)

	// make container
	containerWidth := int(p.ContainerWidth)
	containerHeight := int(p.ContainerHeight)
	containerRectangle := image.Rect(0, 0, containerWidth, containerHeight)
	containerImg := image.NewNRGBA(containerRectangle)
	rgba, err := ParseHexColor(p.Color)
	if err != nil {
		return "", err
	}
	// Set the background color.
	for x := 0; x < containerWidth; x++ {
		for y := 0; y < containerHeight; y++ {
			containerImg.Set(x, y, rgba)
		}
	}

	file, err := os.Open(p.File)
	if err != nil {
		return "", err
	}
	defer file.Close()
	rawImg, err := jpeg.Decode(file)
	if err != nil {
		return "", err
	}

	// resize the photos
	rawImg = imaging.Resize(rawImg, int(p.PhotoWidth), int(p.PhotoHeight), imaging.Lanczos)

	// Calculate the spacing: (containerWidth−(imageWidth×imageCount))/(imageCount+1)(containerWidth−(imageWidth×imageCount))/(imageCount+1)
	gapX := int((p.ContainerWidth - (float64(rawImg.Bounds().Dx()) * floorX)) / (floorX + 1))
	gapY := int((p.ContainerHeight - (float64(rawImg.Bounds().Dy()) * floorY)) / (floorY + 1))

	// layout
	for i := 0; i < int(floorX); i++ {
		for j := 0; j < int(floorY); j++ {
			offsetX := i*(rawImg.Bounds().Dx()+gapX) + gapX
			offsetY := j*(rawImg.Bounds().Dy()+gapY) + gapY

			draw.Draw(containerImg, containerRectangle.Add(image.Pt(offsetX, offsetY)), rawImg, rawImg.Bounds().Min, draw.Src)
		}
	}

	// save to xx_layout.jpeg
	layoutFilePath := fmt.Sprintf("%s%s", strings.Split(p.File, ".")[0], "_layout.jpeg")
	saveFile, err := os.Create(layoutFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if err := jpeg.Encode(saveFile, containerImg, nil); err != nil {
		return "", err
	}

	return layoutFilePath, nil
}

const INCH float64 = 25.4

func ConverMmToPx(dpi, mm float64) (px float64) {
	return (dpi * mm) / INCH
}

var errInvalidFormat = errors.New("invalid format")

func ParseHexColor(hexColor string) (c color.RGBA, err error) {
	c.A = 0xff

	if hexColor[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(hexColor) {
	case 7:
		c.R = hexToByte(hexColor[1])<<4 + hexToByte(hexColor[2])
		c.G = hexToByte(hexColor[3])<<4 + hexToByte(hexColor[4])
		c.B = hexToByte(hexColor[5])<<4 + hexToByte(hexColor[6])
	case 4:
		c.R = hexToByte(hexColor[1]) * 17
		c.G = hexToByte(hexColor[2]) * 17
		c.B = hexToByte(hexColor[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}
