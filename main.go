package main

import (
	"flag"
	svg "github.com/ajstarks/svgo"
	"log"
	"os"
)

var file = flag.String("file", "out.svg", "output file")

func main() {
	flag.Parse()

	f, err := os.Create(*file)
	if err != nil {
		log.Fatal(err)
	}

	width := 48
	height := 48
	canvas := svg.New(f)

	canvas.Startunit(width, height, "in")
	canvas.Circle(width/2, height/2, 100)
	canvas.Text(width/2, height/2, "Hello, SVG!", "text-anchor:middle;font-size:30px;fill:white")
	canvas.End()
}
