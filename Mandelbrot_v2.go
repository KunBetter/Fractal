package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"math/cmplx"
	"net/http"
)

func handle(w http.ResponseWriter, r *http.Request) {
	draw(w)
}

func draw(w io.Writer) {
	const size = 1000
	rec := image.Rect(0, 0, size, size)
	img := image.NewRGBA(rec)

	for y := 0; y < size; y++ {
		yy := 4 * (float64(y)/size - 0.5) // [-2, 2]
		for x := 0; x < size; x++ {
			xx := 4 * (float64(x)/size - 0.5) // [-2, 2]
			c := complex(xx, yy)

			img.Set(x, y, mandelbrot(c))
		}
	}

	png.Encode(w, img)
}

// z := z^2 + c
// 特点，如果 c in M，则 |c| <= 2; 反过来不一定成立
// 如果  c in M，则 |z| <= 2. 这个特性可以用来发现 c 是否属于 M
func mandelbrot(c complex128) color.Color {
	var z complex128
	const iterator = 254

	// 如果迭代 200 次发现 z 还是小于 2，则认为 c 属于 M
	for i := uint8(0); i < iterator; i++ {
		if cmplx.Abs(z) > 2 {
			return getColor(i)
		}
		z = z*z + c
	}

	return color.Black
}

// 根据迭代次数计算一个合适的像素值
func getColor(n uint8) color.Color {
	// 这里乘以 15 是为了提高颜色的区分度，即对比度
	return color.Gray{n * 15}
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8080", nil)
}
