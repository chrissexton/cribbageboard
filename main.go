package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	svg "github.com/ajstarks/svgo"
)

var format = flag.String("format", "nc", "output format (svg,nc)")
var dpi = 72

var oneEighth = int(math.Ceil(float64(dpi) / 8))

type cursor interface {
	Circle(x, y, r int, s ...string)
	End()
	Start(w, h int, ns ...string)
}

type nc struct {
	width  int
	height int

	f *os.File

	zTravel float64
	zDepth  float64
}

func (n nc) Circle(x, y, r int, s ...string) {
	fx := float64(x)/float64(dpi)
	fy := float64(y)/float64(dpi)
	fmt.Fprintf(n.f, "G0 X%.5f Y%.5f\n", fx, fy)
	fmt.Fprintf(n.f, "G1 Z%.5f F9.0\n", 0.0)
	fmt.Fprintf(n.f, "G1 Z%.5f F9.0\n", n.zDepth)
	fmt.Fprintf(n.f, "G1 Z%.5f F9.0\n", n.zTravel)
}

func (n nc) End() {
	fmt.Fprintf(n.f, "G20\n")
	fmt.Fprintf(n.f, "G90\n")
	fmt.Fprintf(n.f, "G1 Z%.5f F9.0\n", n.zTravel)
	fmt.Fprintf(n.f, "G0 X%.5f Y%.5f\n", 0.0, 0.0)
	fmt.Fprintf(n.f, "G4 P0.1\n")
}
func (n *nc) Start(w, h int, ns ...string) {
	n.width = w
	n.height = h

	n.zDepth = -0.125
	n.zTravel = 0.15

	fmt.Fprintf(n.f, "G20\n")
	fmt.Fprintf(n.f, "G90\n")
	fmt.Fprintf(n.f, "G1 Z%.5f F9.0\n", n.zTravel)
}

func mkNC(f *os.File) *nc {
	return &nc{f: f}
}

func main() {
	flag.Parse()

	width := oneEighth * 28 * 2
	height := oneEighth * 120

	var canvas cursor
	switch *format {
	case "svg":
		canvas = mkSVG(os.Stdout)
	case "nc":
		canvas = mkNC(os.Stdout)
	}

	canvas.Start(width, height)

	x := oneEighth * 5
	y := oneEighth

	groupSeparator := oneEighth * 3
	offset := oneEighth * 6

	row(canvas, x, y+offset)
	row(canvas, x+groupSeparator, y+offset)

	col(canvas, x+offset, y)
	col(canvas, x+offset, y+groupSeparator)

	bottomY := offset*2 + 101*oneEighth

	col(canvas, x+offset, y+bottomY)
	col(canvas, x+offset, y+groupSeparator+bottomY)

	x += oneEighth*11 + oneEighth*15*2

	row(canvas, x, y+offset)
	row(canvas, x+groupSeparator, y+offset)

	canvas.End()
}

func mkSVG(f *os.File) cursor {
	canvas := svg.New(f)

	return canvas
}

func col(canvas cursor, x, y int) {
	for i := 0; i < 3; i++ {
		clusterHorizontal(canvas, x, y)
		x += 2*oneEighth + 5*2*oneEighth
	}
}

func row(canvas cursor, x, y int) {
	for i := 0; i < 9; i++ {
		clusterVertical(canvas, x, y)
		y += 2*oneEighth + 5*2*oneEighth
	}
}

func clusterVertical(canvas cursor, x, y int) {
	for i := 0; i < 5; i++ {
		canvas.Circle(x, y, oneEighth/2)
		y = y + oneEighth*2
	}
}

func clusterHorizontal(canvas cursor, x, y int) {
	for i := 0; i < 5; i++ {
		canvas.Circle(x, y, oneEighth/2)
		x = x + oneEighth*2
	}
}
