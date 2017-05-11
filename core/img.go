package Fractal

import (
	"bytes"
	"encoding/base64"
	"flag"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"math"
)

var root = flag.String("root", ".", "file system path")

func Img() {
	http.HandleFunc("/blue/", blueHandler)
	http.HandleFunc("/red/", redHandler)
	http.HandleFunc("/fractal", fractalHandler)
	http.Handle("/", http.FileServer(http.Dir(*root)))
	log.Println("Listening on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func blueHandler(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	var img image.Image = m
	writeImage(w, &img)
}

func redHandler(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	blue := color.RGBA{255, 0, 0, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	var img image.Image = m
	writeImageWithTemplate(w, &img)
}

func repeat(z, c *Complex) color.RGBA {
	for k := 0; k < 256; k++ {
		v2 := z.real * z.real + z.imag * z.imag
		if v2 > 4 {
			f := int(math.Sqrt(v2)) - 2
			//TODO color
			return color.RGBA{uint8(k + f * 11 + 255), uint8(k * f * 5), uint8(k * f * 3 + 255), 255}
		} else {
			z = z.Multiply(z).Add(c)
		}
	}

	return color.RGBA{255, 255, 255, 255}
}

func fractalHandler(w http.ResponseWriter, r *http.Request) {
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

// writeImage encodes an image 'img' in jpeg format and writes it into ResponseWriter.
func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}