package main

import (
	"github.com/KunBetter/Fractal/core"
	"net/http"
	"log"
)

func main() {
	mandelbrot := &Fractal.Mandelbrot{}

	mux := http.NewServeMux()
	mux.Handle("/mandelbrot", http.HandlerFunc(mandelbrot.Render))
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}
