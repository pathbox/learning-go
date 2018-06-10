package main

import (
	"io/ioutil"
	"log"

	pigo "github.com/esimov/pigo/core"
)

func main() {
	cascadeFile, err := ioutil.ReadFile("/path/to/cascade/file")
	if err != nil {
		log.Fatalf("Error reading the cascade file: %v", err)
	}

	src, err := pigo.GetImage("/path/to/image")
	if err != nil {
		log.Fatalf("Cannot open the image file: %v", err)
	}

	sampleImg := pigo.RgbToGrayscale(src)

	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
	}
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y
	imgParams := pigo.ImageParams{sampleImg, rows, cols, cols}

	pigo := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := pigo.Unpack(cascadeFile)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := classifier.RunCascade(imgParams, cParams)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.2)
}
