package paletter

import (
	"image"
	"sort"

	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

// Controls delta threshold for k-means iterations
// (stops iterating when less than N percent points changed cluster assignment)
var DeltaThreshold = 0.05

// Extracts Observations from a given Image.
// Each pixel is converted to colorful.Color, then unpacked in its raw L*, a*
// and b* values as floats
func ImageToObservation(img image.Image) clusters.Observations {
	var ret clusters.Observations

	bounds := img.Bounds()
	// Iterate over every pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Extract color as L*a*b*
			colorfulColor, _ := colorful.MakeColor(img.At(x, y))
			l, a, b := colorfulColor.Lab()

			// Convert color to float values
			labColor := ColorLab{
				l, a, b,
			}

			ret = append(ret, labColor)
		}
	}

	return ret
}

// Extracts a given number of colors (as cluster centers) from Observations through naive k-Means clustering
func CalculatePalette(obs clusters.Observations, nColors int) (clusters.Clusters, error) {
	// Initalize a kmeans object with observations
	km, err := kmeans.NewWithOptions(DeltaThreshold, nil)
	if err != nil {
		return clusters.Clusters{}, err
	}

	// Run clustering
	c, err := km.Partition(obs, nColors)
	return c, err
}

// Converts Clusters centers to a color palette of a given number of colors.
// Each color gets 3 center values as floats (L*, a* and b* - the cluster space for an image
// is 3-dimensional)
func ColorsFromClusters(cs clusters.Clusters) []colorful.Color {
	var ret []colorful.Color
	var colors Palette

	for _, c := range cs {
		colors = append(colors, ColorLab{c.Center[0], c.Center[1], c.Center[2]})
	}

	// Sort colors by lightness
	sort.Sort(colors)

	for _, c := range colors {
		ret = append(ret, colorful.Lab(c.L, c.A, c.B).Clamped())
	}

	return ret
}
