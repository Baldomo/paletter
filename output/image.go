package output

import (
	"image"
	"image/png"
	"math"
	"os"
	"path"
	"path/filepath"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/draw"
)

const (
	border            = 4
	paletteHeightPerc = 0.5
)

// Describes a simple PNG file output. See "images" (or the main README) for sample outpts
type PNG struct {
	imgPath     string
	outFileName string
}

var _ Output = &PNG{}

func NewPNG(imgPath string, outFilename string) *PNG {
	return &PNG{
		imgPath, outFilename,
	}
}

// Calculates the height of the palette using a percentage
// of the height of the original image
func calcPaletteHeight(srcHeight int) int {
	return int(paletteHeightPerc * float32(srcHeight))
}

// Calculates width of a palette color and total remainder (gap between palette end and original image end)
func calcPaletteWidth(srcWidth int, nColors int) (width int, totalRemainder int) {
	calc := float64(srcWidth-(border*(nColors-1))) / float64(nColors)
	_, remainder := math.Modf(calc)
	return int(math.Floor(calc)), int(math.Round(remainder * float64(nColors)))
}

// Calculates the Rectangle for the final image
func imgRect(srcWidth int, srcHeight int) image.Rectangle {
	width := (border * 2) + srcWidth
	height := (border * 3) + srcHeight + calcPaletteHeight(srcHeight)
	return image.Rect(0, 0, width, height)
}

// Calculates the layout for the palette in the final image
func colorRects(srcWidth int, srcHeight int, nColors int) []image.Rectangle {
	var ret []image.Rectangle

	// Calculate width and height
	width, remainder := calcPaletteWidth(srcWidth, nColors)
	height := calcPaletteHeight(srcHeight)

	// Vertical offset is 2 border widths plus the height of the original image
	yOffset := (border * 2) + srcHeight
	for i := 0; i < nColors; i++ {
		xOffset := i*(width+border) + border

		// Cheat by stretching the last color to be flush with the original image
		if i == nColors-1 {
			width += remainder
		}

		ret = append(ret, image.Rect(xOffset, yOffset, width+xOffset, height+yOffset))
	}

	return ret
}

// Serializes a list of colors (palette) into an image with a defined layout, containing
// the original image and the palette
func (p PNG) Write(colors []colorful.Color) error {
	var outFilePath string

	srcImg, err := OpenImage(p.imgPath)
	if err != nil {
		return err
	}

	newImgRect := imgRect(srcImg.Bounds().Dx(), srcImg.Bounds().Dy())
	newImg := image.NewRGBA(newImgRect)

	// Fill background with white
	draw.Src.Draw(newImg, newImgRect, image.White, image.Point{})
	// Draw the original image onto the new one
	draw.Src.Draw(
		newImg,
		srcImg.
			Bounds().
			Add(image.Point{border, border}),
		srcImg,
		image.Point{},
	)

	for i, r := range colorRects(srcImg.Bounds().Dx(), srcImg.Bounds().Dy(), len(colors)) {
		draw.Src.Draw(newImg, r, &image.Uniform{colors[i]}, image.Point{})
	}

	// Place outFile in the same directory as the provided image
	if p.outFileName == "" {
		outFilePath = path.Join(filepath.Dir(p.imgPath), "palette.png")
	} else {
		outFilePath, _ = filepath.Abs(p.outFileName)
	}
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Write final image as png file
	return png.Encode(outFile, newImg)
}
