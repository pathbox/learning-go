package color_extractor

import (
	"image"
	"image/color"
	"math"
	"sort"
)

type bucket struct {
	Red   float64
	Green float64
	Blue  float64
	Count float64
}

type ByCount []bucket

func (c ByCount) Len() int           { return len(c) }
func (c ByCount) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCount) Less(i, j int) bool { return c[i].Count < c[j].Count }

type Config struct {
	DownSizeTo  float64
	SmallBucket float64
}

func ExtractColors(image image.Image) []color.Color {
	return ExtractColorsWithConfig(image, Config{
		DownSizeTo:  224.,
		SmallBucket: .01,
	})
}

func ExtractColorsWithConfig(image image.Image, config Config) []color.Color {
	width := image.Bounds().Max.X
	height := image.Bounds().Max.Y

	stepX := int(math.Max(float64(width)/config.DownSizeTo, 1))
	stepY := int(math.Max(float64(height)/config.DownSizeTo, 1))

	// load image's pixels into buckets
	var buckets [2][2][2]bucket
	totalCount := 0
	for x := 0; x < width; x += stepX {
		for y := 0; y < height; y += stepY {
			color := image.At(x, y)
			r, g, b, a := color.RGBA()
			r >>= 8
			g >>= 8
			b >>= 8
			i := r >> 7
			j := g >> 7
			k := b >> 7
			if a>>8 == 255 {
				buckets[i][j][k].Red += float64(r)
				buckets[i][j][k].Green += float64(g)
				buckets[i][j][k].Blue += float64(b)
				buckets[i][j][k].Count++
				totalCount++
			}
		}
	}

	// calculate bucket's averages
	var bucketsAverages []bucket
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				currentBucket := buckets[i][j][k]
				if currentBucket.Count > 0 {
					bucketsAverages = append(bucketsAverages, bucket{
						Count: currentBucket.Count,
						Red:   currentBucket.Red / currentBucket.Count,
						Green: currentBucket.Green / currentBucket.Count,
						Blue:  currentBucket.Blue / currentBucket.Count,
					})
				}
			}
		}
	}

	// sort buckets by bucket size
	sort.Sort(sort.Reverse(ByCount(bucketsAverages)))

	// export color.Color from bucket, ignore small buckets
	colors := []color.Color{}
	for _, avg := range bucketsAverages {
		if avg.Count/totalCount > config.SmallBucket {
			colors = append(colors, color.RGBA{
				R: uint8(avg.Red),
				G: uint8(avg.Green),
				B: uint8(avg.Blue),
				A: 255,
			})
		}
	}

	return colors
}
