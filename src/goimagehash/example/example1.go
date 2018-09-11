package main

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/corona10/goimagehash"
)

func main() {
	qFile, err := os.Open("pic1.jpeg") // according to your right path
	file1, err := os.Open("pic2.jpeg")
	file2, err := os.Open("pic3.jpg")
	if err != nil {
		fmt.Println("err1===")
		panic(err)
	}
	defer qFile.Close()
	defer file1.Close()
	defer file2.Close()

	imgQuery, err := jpeg.Decode(qFile)
	img1, err := jpeg.Decode(file1)
	img2, err := jpeg.Decode(file2)
	if err != nil {
		fmt.Println("err2===")
		panic(err)
	}
	queryHash, err := goimagehash.AverageHash(imgQuery)
	hash1, err := goimagehash.AverageHash(img1)
	hash2, err := goimagehash.AverageHash(img2)
	if err != nil {
		fmt.Println("err3===")
		panic(err)
	}
	distance1, err := queryHash.Distance(hash1)
	distance2, err := queryHash.Distance(hash2)
	if err != nil {
		fmt.Println("err4===")
		panic(err)
	}
	fmt.Printf("AverageHash Distance between images: %d %d\n", distance1, distance2)
	hash1, _ = goimagehash.DifferenceHash(img1)
	hash2, _ = goimagehash.DifferenceHash(img2)
	distance1, _ = queryHash.Distance(hash1)
	distance2, _ = queryHash.Distance(hash2)
	fmt.Printf("DifferenceHash Distance between images: %d %d\n", distance1, distance2)
	queryHash, _ = goimagehash.PerceptionHash(imgQuery)
	hash1, _ = goimagehash.PerceptionHash(img1)
	hash2, _ = goimagehash.PerceptionHash(img2)
	distance1, _ = queryHash.Distance(hash1)
	distance2, _ = queryHash.Distance(hash2)
	fmt.Printf("PerceptionHash Distance between images: %d %d\n", distance1, distance2)

	// distance is shorter, the result is more similar

}
