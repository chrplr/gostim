// Displays a moving vertical bar between two horizontal flashing bar
// Author: Christophe Pallier <christophe@pallier.org>
// Licence: GPL-3

package main

import (
	"fmt"
	"slices"

	"github.com/veandco/go-sdl2/sdl"
)

var window *sdl.Window
var renderer *sdl.Renderer
var timerResolution float64
var screenWidth, screenHeight int32

func toMilliseconds(ticks uint64) float64 {
	return 1000.0 * float64(ticks) / timerResolution
}

func run() (timings []float64) {

	rectTop := sdl.Rect{X: 0, Y: 0, W: screenWidth, H: 200}
	rectBottom := sdl.Rect{X: 0, Y: screenHeight - 200, W: screenWidth, H: 200}

	bar := sdl.Rect{X: 0, Y: 250, W: 8, H: screenHeight - 500}
	var xpos int32 = 0
	var skip int32 = 8

	loop := 0
	running := true
	visibleRect := true

	for running {

		loop++

		start := sdl.GetPerformanceCounter()

		// clear screen
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// display horizontal rectangles
		visibleRect = loop%8 == 0
		if visibleRect {
			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.FillRect(&rectTop)
			renderer.FillRect(&rectBottom)
		} else {
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.FillRect(&rectTop)
			renderer.FillRect(&rectBottom)
		}

		// display vertical bar & move it
		bar.X = xpos
		bar.W = skip
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&bar)
		xpos = (xpos + skip) % screenWidth

		// detect any key press
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				running = false
			}
		}

		renderer.Present()

		delta := toMilliseconds(sdl.GetPerformanceCounter() - start)
		timings = append(timings, delta)
	}
	return timings
}

func printHistogram(timings []float64) {
	hist := make([]uint, 1+uint(slices.Max(timings)))
	for _, v := range timings {
		i := uint(v)
		hist[i]++
	}

	fmt.Println("time(ms) N")
	for index, num := range hist {
		fmt.Printf("%4d %5d\n", index, num)
	}
}

func main() {
	var err error

	//window, err := sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
	//	0, 0, sdl.WINDOW_FULLSCREEN_DESKTOP)
	if window, err = sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1920, 1080, sdl.WINDOW_FULLSCREEN); err != nil {
		panic(err)
	}
	defer window.Destroy()

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC); err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	sdl.Delay(3000)

	screenWidth, screenHeight, _ = renderer.GetOutputSize()
	timerResolution = float64(sdl.GetPerformanceFrequency())

	timings := run()
	printHistogram(timings)

}
