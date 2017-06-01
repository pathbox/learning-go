package main

import (
	"github.com/tidwall/pinhole"
)

func main() {
	p := pinhole.New()
	p.DrawCube(-0.3, -0.3, -0.3, 0.3, 0.3, 0.3)
	p.SavePNG("cube.png", 500, 500, nil)
}
