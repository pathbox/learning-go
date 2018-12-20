/*
 * Insert an image to a PDF file.
 *
 * Adds image to a specific page of a PDF.  xPos and yPos define the upper left corner of the image location, and width
 * is the width of the image in PDF coordinates (height/width ratio is maintained).
 *
 * Example go run pdf_add_image_to_page.go /tmp/input.pdf 1 /tmp/image.jpg 0 0 100 /tmp/output.pdf
 * adds the image to the upper left corner of the page (0,0).  The width is 100 (typical page width 612 with defaults).
 *
 * Syntax: go run pdf_add_image_to_page.go input.pdf <page> image.jpg <xpos> <ypos> <width> output.pdf
 */

package main

import (
	"fmt"
	"os"

	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
	pdf "github.com/unidoc/unidoc/pdf/model"
	"github.com/unidoc/unidoc/pdf/model/optimize"
)

func main() {
	// Use debug logging.
	unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	inputPath := "/Users/pathbox/busi.pdf"
	pageNum := 9
	imagePath := "/Users/pathbox/busi.png"

	xPos := float64(100)
	yPos := float64(150)
	iwidth := float64(420)
	outputPath := "/Users/pathbox/test-image-done.pdf"

	fmt.Printf("xPos: %v, yPos: %v\n", xPos, yPos)

	err := addImageToPdf(inputPath, outputPath, imagePath, pageNum, xPos, yPos, iwidth)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

}

// Add image to a specific page of a PDF.  xPos and yPos define the upper left corner of the image location, and iwidth
// is the width of the image in PDF document dimensions (height/width ratio is maintained).
func addImageToPdf(inputPath string, outputPath string, imagePath string, pageNum int, xPos float64, yPos float64, iwidth float64) error {
	c := creator.New()
	optim := optimize.New(optimize.Options{
		CombineDuplicateDirectObjects:   true,
		CombineIdenticalIndirectObjects: true,
		ImageUpperPPI:                   100.0,
		UseObjectStreams:                true,
		ImageQuality:                    50,
		CombineDuplicateStreams:         true,
	})
	c.SetOptimizer(optim)
	// Prepare the image.
	img, err := c.NewImageFromFile(imagePath)
	if err != nil {
		return err
	}
	img.ScaleToWidth(iwidth)
	img.SetPos(xPos, yPos)

	// Read the input pdf file
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f) // pdf  reader
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	// Load the page
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1) // read every page
		if err != nil {
			return nil
		}

		//Add the page
		err = c.AddPage(page) // add this page
		if err != nil {
			return err
		}

		img1, err := c.NewImageFromFile(imagePath)
		if err != nil {
			return err
		}
		img1.ScaleToWidth(100)
		img1.SetPos(100, 500)
		// If the specified page, or -1, apply the image to the page
		if i+1 == pageNum || pageNum == -1 {
			_ = c.Draw(img) // draw image
			// c.Draw(img1)
		}
	}

	err = c.WriteToFile(outputPath)
	return err
}
