package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
)

func main() {
	file, err := os.Open("/Users/pathbox/busi.png")
	if err != nil {
		log.Fatal(err)
	}
	t := time.Now()
	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(1000, 0, img, resize.Lanczos3)

	out, err := os.Create("/Users/pathbox/test_bbcc.png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)
	x := time.Now().Sub(t)
	fmt.Println(x)
}
