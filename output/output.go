package output

import "github.com/lucasb-eyer/go-colorful"

// A generic output image generator
type Output interface {
	// The Write method receives the colors found in an image and should perform
	// the necessary I/O such as generating a file and writing to disk
	Write(colors []colorful.Color) error
}
