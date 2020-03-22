package paletter

import (
	"github.com/muesli/clusters"
)

// Represent color as lightness (L), green (-) to red (+) (A) and blue (-) to yellow (+) (B).
// Using floats instead of converting to other color space makes computing distance much faster,
// while keeping the computations accurate as L*a*b* values are represented linearly.
// ColorLab implements clusters.Observation
type ColorLab struct {
	L, A, B float64
}

// Returns coordinates as a Coordinates object (slice of float64's)
func (c ColorLab) Coordinates() clusters.Coordinates {
	return clusters.Coordinates{c.L, c.A, c.B}
}

// Computes squared distance between `c` and another color (as Coordinates)
func (c ColorLab) Distance(pos clusters.Coordinates) float64 {
	dx := c.L - pos[0]
	dy := c.A - pos[1]
	dz := c.B - pos[2]

	return dx*dx + dy*dy + dz*dz
}

// Palette implements sort.Interface
type Palette []ColorLab

// Returns the number of elements in Palette
func (p Palette) Len() int {
	return len(p)
}

func (p Palette) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Sorts colors by lightness, in descending order
func (p Palette) Less(i, j int) bool {
	return p[i].L > p[j].L
}
