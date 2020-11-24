package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	svg "github.com/ajstarks/svgo"
)

var (
	format  = flag.String("format", "nc", "output format (svg,nc)")
	zDepth  = flag.Float64("zdepth", 0.125, "material thickness")
	zTravel = flag.Float64("ztravel", 0.150, "safe travel height")
	bitSize = flag.Float64("bitsize", 0.125, "diameter of end mill")
	outFile = flag.String("out", "", "file output, empty for stdout")

	dpi       = 72
	oneEighth = int(math.Ceil(float64(dpi) / 8))
)

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
	bitSize float64
}

func (n nc) Circle(x, y, r int, s ...string) {
	fx := float64(x) / float64(dpi)
	fy := float64(y) / float64(dpi)
	fmt.Fprintf(n.f, "G0 X%.5f Y%.5f\n", fx, fy)
	// For the drill operation, plunge half the bit size at a time and then
	// resurface for material removal
	for d := 0-n.bitSize/2; d >= n.zDepth; d -= n.bitSize / 2 {
		fmt.Fprintf(n.f, "G1 Z%.5f F9.0\n", d)
		// If we're not at depth, we need to move up, otherwise this is a wasted
		// instruction since we will immediately move to travel height
		if d > n.zDepth {
			fmt.Fprintf(n.f, "G1 Z%.5f F9.0\n", 0.0)
		}
	}
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

	n.bitSize = 0.125

	fmt.Fprintf(n.f, "G20\n")
	fmt.Fprintf(n.f, "G90\n")
	fmt.Fprintf(n.f, "G1 Z%.5f F9.0\n", n.zTravel)
}

func mkNC(f *os.File, zDepth, zTravel, bitSize float64) *nc {
	return &nc{
		f:       f,
		zDepth:  zDepth,
		zTravel: zTravel,
		bitSize: bitSize,
	}
}

func main() {
	flag.Parse()

	width := oneEighth * 28 * 2
	height := oneEighth * 120

	out := os.Stdout
	if *outFile != "" {
		var err error
		out, err = os.Create(*outFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	var canvas cursor
	switch *format {
	case "svg":
		canvas = mkSVG(out)
	case "nc":
		canvas = mkNC(out, *zDepth, *zTravel, *bitSize)
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
