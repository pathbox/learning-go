package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nfnt/resize"
)

func main() {
	// imagePath := "/Users/pathbox/bbcc.jpeg"
	imagePath := "/Users/pathbox/econtract_tmp_file/user_photo/200000437-IDPhotoBack.jpeg"
	stream, err := ioutil.ReadFile(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	r := bytes.NewReader(stream)

	// bufReader := bufio.NewReader(r)
	t := time.Now()
	img, _, err := image.Decode(r)
	if err != nil {
		fmt.Println("1")
		log.Fatal(err)
	}

	m := resize.Resize(800, 0, img, resize.Lanczos3)

	out, err := os.Create("./image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	jpeg.Encode(out, m, nil)
	x := time.Now().Sub(t)
	fmt.Println(x)
}
