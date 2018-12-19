package main

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
)

func main() {
	file, err := os.Open("/Users/pathbox/bbcc.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	t := time.Now()
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(600, 0, img, resize.Lanczos3)

	out, err := os.Create("./test_bbcc.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)
	x := time.Now().Sub(t)
	fmt.Println(x)
}
