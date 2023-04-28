// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main is a simple demo to show how to use the hilbert library
// When ran, this demo will create the following images:
//
//	hilbert.png, hilbert_animation.gif, peano.png, and peano_animation.gif
//
// It is suggested you optimise/compress both images before uploading.
//
//	go run demo/demo.go
//	zopflipng -y logo.png images/logo.png
//	zopflipng -y hilbert.png images/hilbert.png
//	zopflipng -y peano.png images/peano.png
//	gifsicle -O -o images/hilbert_animation.gif hilbert_animation.gif
//	gifsicle -O -o images/peano_animation.gif peano_animation.gif
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"

	"math"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/google/hilbert"
	"github.com/google/hilbert/demo/lib"
)

// spaceFillingImage facilitates the drawing of a space filing curve.
type spaceFillingImage struct {
	Curve hilbert.SpaceFilling

	// Size of each square in pixels
	SquareWidth  float64
	SquareHeight float64

	DrawGrid   bool
	DrawText   bool    // Should text be drawn on the image
	TextMargin float64 // Margin around text in pixels

	BackgroundColor color.Color
	GridColor       color.Color
	TextColor       color.Color
	SnakeColor      color.Color

	GridWidth  float64
	SnakeWidth float64
}

// createSpaceFillingImage returns a new SpaceFillingImage ready for drawing.
// squareWidth and squareHeight are the dimensions of each individual square in the resulting image.
func createSpaceFillingImage(curve hilbert.SpaceFilling, squareWidth, squareHeight float64) *spaceFillingImage {
	return &spaceFillingImage{
		Curve: curve,

		SquareWidth:  squareWidth,
		SquareHeight: squareHeight,

		// All the default values

		DrawGrid:   true,
		DrawText:   true,
		TextMargin: 3.0,

		BackgroundColor: color.RGBA{0xee, 0xee, 0xff, 0xff},
		GridColor:       color.White,
		TextColor:       color.RGBA{0x33, 0x33, 0x33, 0xff},
		SnakeColor:      color.RGBA{0x33, 0x33, 0x33, 0xff},

		GridWidth:  1.0,
		SnakeWidth: 2.0,
	}
}

func (h *spaceFillingImage) toPixel(x, y int) (float64, float64) {
	return float64(x) * h.SquareWidth, float64(y) * h.SquareHeight
}

func (h *spaceFillingImage) drawGrid(gc *gg.Context, width, height int) {

	// Draw grid, vertical then horizontal lines
	for x := 0; x <= width; x++ {
		gc.MoveTo(h.toPixel(x, 0))
		gc.LineTo(h.toPixel(x, height))
	}

	for y := 0; y < height; y++ {
		gc.MoveTo(h.toPixel(0, y))
		gc.LineTo(h.toPixel(width, y))
	}

	gc.SetLineWidth(h.GridWidth)
	gc.SetColor(h.GridColor)
	gc.Stroke()
}

// Draw uses the parameters in the hilbertImage and returns a Image
func (h *spaceFillingImage) Draw() (*gg.Context, error) {

	width, height := h.Curve.GetDimensions()
	pwidth, pheight := h.toPixel(width, height)

	gc := gg.NewContext(int(pwidth), int(pheight))
	gc.SetColor(h.BackgroundColor)
	gc.Clear()

	if h.DrawGrid {
		h.drawGrid(gc, width, height)
	}

	for t := 0; t < width*height; t++ {

		// Map the 1D number into the 2D space
		x, y, err := h.Curve.Map(t)
		if err != nil {
			return nil, err
		}

		px, py := h.toPixel(x, y)

		// Draw the grid for t
		if h.DrawText {
			text := strconv.Itoa(t)

			gc.SetColor(h.TextColor)
			gc.DrawStringAnchored(text, px+h.TextMargin, py, 0, 1)
		}

		// Move the snake along
		centerX, centerY := px+h.SquareWidth/2, py+h.SquareHeight/2
		if t == 0 {
			gc.MoveTo(centerX, centerY)
		} else {
			gc.LineTo(centerX, centerY)
		}
	}

	// Draw the snake at the end, to form one continuous line.
	gc.SetColor(h.SnakeColor)
	gc.SetLineWidth(h.SnakeWidth)

	gc.SetLineCap(gg.LineCapSquare)
	gc.SetLineJoin(gg.LineJoinRound)

	gc.Stroke()

	return gc, nil
}

func (h *spaceFillingImage) DrawSwap() (*gg.Context, error) {

	width, height := h.Curve.GetDimensions()
	pwidth, pheight := h.toPixel(width, height)

	gc := gg.NewContext(int(pwidth), int(pheight))
	gc.SetColor(h.BackgroundColor)
	gc.Clear()

	if h.DrawGrid {
		h.drawGrid(gc, width, height)
	}

	mat := make([][]int, height)
	for i := range mat {
		mat[i] = make([]int, width)
	}
	for t := 0; t < width*height; t++ {

		// Map the 1D number into the 2D space
		x, y, err := h.Curve.Map(t)
		if err != nil {
			return nil, err
		}
		mat[x][y] = t

		px, py := h.toPixel(x, y)
		// Move the snake along
		centerX, centerY := px+h.SquareWidth/2, py+h.SquareHeight/2
		if t == 0 {
			gc.MoveTo(centerX, centerY)
		} else {
			gc.LineTo(centerX, centerY)
		}
	}

	// stir(mat, width)
	rt := kbucket(mat, mat[3][3], 3, 3, width)
	printrt(rt)
	// bf(mat, width)
	// fmt.Printf("stired matrix:\n")
	// for i := 0; i < width; i++ {
	// 	fmt.Printf("%d\n", mat[i])
	// }
	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			px, py := h.toPixel(x, y)
			t := mat[x][y]
			// Draw the grid for t
			if h.DrawText {
				text := strconv.Itoa(t)

				gc.SetColor(h.TextColor)
				gc.DrawStringAnchored(text, px+h.TextMargin, py, 0, 1)
			}
		}
	}

	// Draw the snake at the end, to form one continuous line.
	gc.SetColor(h.SnakeColor)
	gc.SetLineWidth(h.SnakeWidth)

	gc.SetLineCap(gg.LineCapSquare)
	gc.SetLineJoin(gg.LineJoinRound)

	gc.Stroke()

	return gc, nil
}

func mainDrawOne(filename string, curve hilbert.SpaceFilling) error {
	log.Printf("Drawing one image %q", filename)

	img, err := createSpaceFillingImage(curve, 64, 64).DrawSwap()
	if err != nil {
		return err
	}
	return img.SavePNG(filename)
}

func mainDrawAnimation(filename string, newCurve func(n int) hilbert.SpaceFilling, min, max int) error {
	log.Printf("Drawing animation %q", filename)

	iterations := max - min
	imageWidth, imageHeight := 512.0, 512.0

	g := gif.GIF{
		Image:     make([]*image.Paletted, iterations),
		Delay:     make([]int, iterations),
		LoopCount: 0,
	}

	for i := 0; i < iterations; i++ {
		log.Printf("    Drawing frame %d", i)

		curve := newCurve(min + i)

		width, height := curve.GetDimensions()
		h := createSpaceFillingImage(curve, imageWidth/float64(width), imageHeight/float64(height))
		h.DrawText = false
		img, err := h.Draw()
		if err != nil {
			return err
		}

		g.Image[i] = lib.ConvertToPaletted(img.Image())
		g.Delay[i] = 200 // 200 x 100th of a second = 2 second
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return gif.EncodeAll(f, &g)
}

func mainDrawLogo(filename string, curve hilbert.SpaceFilling) error {
	const scale = 8

	log.Printf("Drawing logo %q", filename)

	h := createSpaceFillingImage(curve, math.Pow(2, scale), math.Pow(2, scale))
	h.DrawText = false
	h.DrawGrid = false
	h.SnakeWidth = math.Pow(2, scale-2)
	h.BackgroundColor = color.Transparent

	img, err := h.Draw()
	if err != nil {
		return err
	}
	return img.SavePNG(filename)
}

func main() {

	newHilbert := func(n int) hilbert.SpaceFilling {
		s, err := hilbert.NewHilbert(int(math.Pow(2, float64(n))))
		if err != nil {
			panic(fmt.Errorf("failed to create hilbert space: %s", err.Error()))
		}
		return s
	}

	if err := mainDrawOne("hilbert.png", newHilbert(4)); err != nil {
		log.Fatalf("Failed to draw image: %s", err.Error())
	}
}

func stir(arr [][]int, n int) {
	j := 0
	for i := n - 2; i >= n/2; i -= 2 {
		tmp := arr[i]
		arr[i] = arr[j]
		arr[j] = tmp
		j += 2
	}

	k := 0
	for i := n - 2; i >= n/2; i -= 2 {
		for j := 0; j < n; j++ {
			tmp := arr[j][i]
			arr[j][i] = arr[j][k]
			arr[j][k] = tmp
		}
		k += 2
	}

	for i := 0; i < n; i++ {
		fmt.Printf("%*d\n", 3, arr[i])
	}
}

type coordinates struct {
	x int
	y int
}

type routingrRecord struct {
	x    int
	y    int
	id   int
	dist int
}

func ManhattanDist(c1, c2 coordinates) int {
	return abs(c1.x-c2.x) + abs(c1.y-c2.y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func msb(n int) int {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n = ((n + 1) >> 1)
	return n
}

// func TrackNeighbors(mat [][]int) {

// }

// func BFS(mat [][]int, n, lo, up int, c coordinates) (int, int) {
// 	queue := make([]coordinates, n*n)
// 	queue = append(queue, 1)
// 	// Top (just get next element, don't remove it)
// 	x = queue[0]
// 	// Discard top element
// 	queue = queue[1:]

// }

func bf(mat [][]int, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			kbucket(mat, mat[i][j], i, j, n)
		}
	}
}

func kbucket(mat [][]int, n, x, y, width int) map[int][]routingrRecord {
	routingTable := make(map[int][]routingrRecord)
	for i := 0; i < width; i++ {
		for j := 0; j < width; j++ {
			dist := ManhattanDist(coordinates{x, y}, coordinates{i, j})
			id := mat[i][j]
			routingTable[msb(n^id)] = append(routingTable[msb(n^id)], routingrRecord{i, j, id, dist})
		}
	}
	return routingTable
}

func printrt(rt map[int][]routingrRecord) {
	for k, v := range rt {
		fmt.Printf("inside bucket %d:\n", k)
		closest := routingrRecord{0, 0, 0, int(^uint(0) >> 1)}
		for _, rec := range v {
			if rec.dist < closest.dist {
				closest = rec
			}
			fmt.Printf("x:%d, y:=%d, id:%d, dist:%d\n", rec.x, rec.y, rec.id, rec.dist)
		}
		fmt.Printf("x:%d, y:=%d, id:%d, dist:%d is most close to node\n", closest.x, closest.y, closest.id, closest.dist)
	}
}
