package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type ValuePoint struct {
	Lat, Lon float64
	Value    float64
}

type LatLonBox struct {
	MinLat, MinLon, MaxLat, MaxLon float64
}

func GetBoundingBox(points []ValuePoint, border float64) (box LatLonBox) {
	box.MinLat = math.Inf(+1)
	box.MinLon = math.Inf(+1)
	box.MaxLat = math.Inf(-1)
	box.MaxLon = math.Inf(-1)
	for p := range points {
		if points[p].Lat < box.MinLat {
			box.MinLat = points[p].Lat
		}
		if points[p].Lat > box.MaxLat {
			box.MaxLat = points[p].Lat
		}
		if points[p].Lon < box.MinLon {
			box.MinLon = points[p].Lon
		}
		if points[p].Lon > box.MaxLon {
			box.MaxLon = points[p].Lon
		}
	}
	lon_border := (box.MaxLon - box.MinLon) * border
	lat_border := (box.MaxLat - box.MinLat) * border
	box.MinLat -= lat_border
	box.MinLon -= lon_border
	box.MaxLat += lat_border
	box.MaxLon += lon_border
	return box
}

func getextremepoints(points []ValuePoint) (min, max float64) {
	minval := math.Inf(+1)
	maxval := math.Inf(-1)
	for i := range points {
		if minval > points[i].Value {
			minval = points[i].Value
		}
		if maxval < points[i].Value {
			maxval = points[i].Value
		}
	}
	return minval, maxval
}

func gaussian(height, dist, width float64) float64 {
	return height * math.Exp(-(dist*dist)/(2*width*width))
}

func gradient(val float64) color.RGBA {
	return color.RGBA{255-uint8(255*val),0,uint8(255*val),255}
}

func CreateHeatmap(points []ValuePoint, box LatLonBox, width, height int) *image.RGBA {
	minval, maxval := getextremepoints(points)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			posx := box.MinLat + float64(x)*(box.MaxLat - box.MinLat)/float64(width)
			posy := box.MinLon + float64(y)*(box.MaxLon - box.MinLon)/float64(height)
			val := float64(0.0)
			weight := float64(0.0)
			for i := range points {
				distsquared := (posx-points[i].Lat)*(posx-points[i].Lat) + (posy-points[i].Lon)*(posy-points[i].Lon)
				val += points[i].Value/distsquared
				weight += 1.0/distsquared
			}
			val = val/weight
			val = (val-minval)/(maxval-minval)
			img.SetRGBA(x, y, gradient(val))		}
	}
	return img
}

func main() {
	width := 2590
	height := 1541
	points := []ValuePoint{
		{1.0, 3.0, 12.0},
		{2.0, 1.0, 8.3},
		{5.0, 2.0, 16.9},
		{4.0, 2.2, 29.3},
		{3.5, 1.8, 4.5},
	}
	box := GetBoundingBox(points, 0.2)
	img := CreateHeatmap(points, box, width, height)
	outfile, _ := os.Create("heatmap.png")
	_ = png.Encode(outfile, img)
	outfile.Close()
}
