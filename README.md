# Paletter
Paletter is a CLI app and library to extract a color palette from an image through clustering (k-Means)

Example output:

![Unsplash (Anders Jildén), Vernazza](images/vernazza-paletted.png)

### Installation
To install paletter just use:
```sh
$ go get -u -v github.com/Baldomo/paletter
```
To install the CLI app use:
```sh
$ go install -v github.com/Baldomo/paletter/cmd/paletter
```

### Usage (CLI)
```
Usage: paletter [OPTIONS] <IMAGE>
Flags:
  -colors int
        Number of colors to extract from the image (default 7)
  -html
        Output an html page
  -out string
        Set output file name/path
  -png
        Output a png image (default true)
```
The CLI app outputs `png` images of resolution dependent on the source image (see `genimage.go` for calculations), with the extracted palette ordered by lightness from left to right

---

### Notes
- Paletter uses the [CIE L\*a\*b\* color space](https://en.wikipedia.org/wiki/CIELAB_color_space) both for the simplicity of its representation and it being device-indipendent
- The k-Means implementation included in paletter is the naive one ([muesli/kmeans](https://github.com/muesli/kmeans/))
