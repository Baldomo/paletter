package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/Baldomo/paletter"
)

var (
	htmlOut bool
	pngOut  bool
	outName string
)

func main() {
	flag.BoolVar(&htmlOut, "html", false, "Output an html page")
	flag.BoolVar(&pngOut, "png", true, "Output a png image")
	flag.StringVar(&outName, "out", "", "Set output file name/path")
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage: paletter [OPTIONS] <IMAGE>\n")
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
	}

	if flag.NArg() == 0 {
		fmt.Println("No arguments supplied!")
		flag.Usage()
		os.Exit(1)
	}

	if htmlOut && pngOut {
		fmt.Println("Either one of -html or -png must be set")
		os.Exit(1)
	}

	// Open image
	img, err := paletter.OpenImage(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	obs := paletter.ImageToObservation(img)
	cs, _ := paletter.CalculatePalette(obs, 7)
	colors := paletter.ColorsFromClusters(cs)

	if htmlOut {
		err := paletter.WriteHTML(flag.Arg(0), outName, colors)
		if err != nil {
			log.Fatal(err)
		}
	} else if pngOut {
		err := paletter.WriteImage(flag.Arg(0), outName, colors)
		if err != nil {
			log.Fatal(err)
		}
	}
}
