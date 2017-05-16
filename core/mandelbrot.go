package Fractal

import (
	"net/http"
	"image"
	"math"
	"image/color"
	"bytes"
	"image/jpeg"
	"encoding/base64"
	"log"
	"html/template"
	"time"
	"fmt"
)

type Mandelbrot struct {
}

func (mandelbrot *Mandelbrot) Render(w http.ResponseWriter, req *http.Request) {
	//纳秒
	startTime := time.Now().UnixNano()
	m := image.NewRGBA(image.Rect(0, 0, 640, 500))

	c := Complex{
		0.285,
		0.01,
	}
	/*
	 * -0.75,0
	 * 0.45, -0.1428
	 * 0.285, 0.01
	 * 0.285, 0
	 * -0.8, 0.156
	 * -0.835, -0.2321
	 * -0.70176, -0.3842
	 */

	for i := 0; i < 640; i++ {
		for j := 0; j < 500; j++ {
			z := Complex{
				float64(i - 320) / 200,
				float64(j - 250) / 200,
			}
			cr := repeat(&z, &c)
			m.Set(i, j, cr)
		}
	}

	var img image.Image = m
	writeImageWithTemplate(w, &img)
	//纳秒转化为毫秒
	costTime := (time.Now().UnixNano() - startTime) / 1e6
	fmt.Printf("costTime = %d millisecond.\n", costTime)
}

func repeat(z, c *Complex) color.RGBA {
	for k := 0; k < 256; k++ {
		v2 := z.real * z.real + z.imag * z.imag
		if v2 > 4 {
			var f float64 = float64(k - 25) / 256;
			return color.RGBA{uint8(255 * math.Pow(f, 0.3)), uint8(255 * math.Pow(f, 0.2)), uint8(255 * math.Pow(f, 0.5)), 255}
		} else {
			z = z.Multiply(z).Add(c)
		}
	}

	return color.RGBA{255, 255, 255, 255}
}

var ImageTemplate string = `<!DOCTYPE html>
			    <html lang="en">
			    <head>
			    </head>
			    <body>
			    <img src="data:image/jpg;base64,{{.Image}}">
			    </body>`

// Writeimagewithtemplate encodes an image 'img' in jpeg format and writes it into ResponseWriter using a template.
func writeImageWithTemplate(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Fatalln("unable to encode image.")
	}

	str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
		log.Println("unable to parse image template.")
	} else {
		data := map[string]interface{}{"Image": str}
		if err = tmpl.Execute(w, data); err != nil {
			log.Println("unable to execute template.")
		}
	}
}