package paletter

import (
	"html/template"
	"os"
	"path"
	"path/filepath"

	"github.com/lucasb-eyer/go-colorful"
)

type templateData struct {
	Title     string
	ImagePath string
	Colors    []string
}

func colorsToHex(colors []colorful.Color) []string {
	var ret []string
	for _, c := range colors {
		ret = append(ret, c.Hex())
	}

	return ret
}

// Creates an HTML document containing the original image and the color palette through
// a template with predefined layout
func WriteHTML(imgPath string, outFileName string, colors []colorful.Color) error {
	var outFilePath string

	templ, err := template.New("picture").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	// Get absolute path from relative path supplied as argument
	filePath, err := filepath.Abs(imgPath)
	if err != nil {
		return err
	}

	// Place outFile in the same directory as the provided image
	if outFileName == "" {
		outFilePath = path.Join(filepath.Dir(imgPath), "palette.png")
	} else {
		outFilePath, _ = filepath.Abs(outFileName)
	}
	outFile, err := os.Create(outFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = templ.Execute(outFile, templateData{
		filepath.Base(imgPath),
		filePath,
		colorsToHex(colors),
	})
	if err != nil {
		return err
	}

	return nil
}

const htmlTemplate = `<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>{{.Title}}</title>
        <body>
            <div style="display: flex;">
                <div style="display: flex; flex: 1;">
                    <img src="{{.ImagePath}}" width="100%">
                </div>
                <div style="display: flex; flex: 2;">
                    <table height="100%" width="100%">
                        <tr>
                            {{range .Colors}}
					        <td bgcolor="{{.}}"></td>
                            {{end}}
                        </tr>
                    </table>
                </div>
            </div>
        </body>
    </head>
</html>
`
