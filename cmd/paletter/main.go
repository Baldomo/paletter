package main

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"runtime/pprof"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/Baldomo/paletter"
	"github.com/Baldomo/paletter/output"
)

var (
	htmlOut bool
	nColors int
	pngOut  bool
	outName string
)

func main() {
	flag.BoolVar(&htmlOut, "html", false, "Output an html page")
	flag.IntVar(&nColors, "colors", 7, "Number of colors to extract from the image")
	flag.BoolVar(&pngOut, "png", true, "Output a png image")
	flag.StringVar(&outName, "out", "", "Set output file name/path")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage: paletter [OPTIONS] <IMAGE>\n")
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	if flag.NArg() == 0 {
		fmt.Println("No arguments supplied!")
		flag.Usage()
		os.Exit(1)
	}

	if !htmlOut && !pngOut {
		fmt.Println("Either one of -html or -png must be set")
		os.Exit(1)
	}

	if nColors < 0 {
		fmt.Println("-colors must be a positive number")
		os.Exit(1)
	}

	pal, err := paletter.FromPath(flag.Arg(0), nColors)
	if err != nil {
		log.Fatal(err)
	}

	var out output.Output = output.NewPNG(flag.Arg(0), outName)
	if htmlOut {
		out = output.NewHTML(flag.Arg(0), outName)
	}

	err = pal.Generate(out)
	if err != nil {
		log.Fatal(err)
	}
}
