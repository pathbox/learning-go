package main

import (
	"github.com/fogleman/gg"
)

func main() {
	dc := gg.NewContext(1000, 1000)
	dc.DrawCircle(500, 500, 300)
	dc.SetRGB(110, 110, 110)
	dc.Fill()
	dc.SavePNG("out.png")
}
