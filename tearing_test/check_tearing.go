// Assess FrameRate
// author: Christophe Pallier
// Licence: GPL-3

package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var timerResolution float64

func toMs(ticks uint64) float64 {
	return 1000.0 * float64(ticks) / timerResolution
}

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var screenWidth, screenHeight int32
	var rect, bar sdl.Rect
	var xpos int32
	var skip int32
	var origin, start, end, dt uint64
	var timings []float64
	var startTimes []uint64
	var running, visibleRect bool

	window, err := sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		0, 0, sdl.WINDOW_FULLSCREEN_DESKTOP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	screenWidth, screenHeight, _ = renderer.GetOutputSize()
	timerResolution = float64(sdl.GetPerformanceFrequency())

	// rect properties
	rect.X = screenWidth / 2
	rect.Y = 0
	rect.W = 400
	rect.H = 200

	// bar properties
	xpos = 0
	skip = 8

	running = true
	visibleRect = true

	sdl.Delay(1000)

	origin = sdl.GetPerformanceCounter()

	for running {

		start = sdl.GetPerformanceCounter() - origin
		startTimes = append(startTimes, start)

		renderer.SetDrawColor(0, 0, 0, 255) // clear screen
		renderer.Clear()

		if visibleRect {
			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.FillRect(&rect)
		} else {
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.FillRect(&rect)
		}
		visibleRect = !visibleRect

		bar.X = xpos
		bar.Y = 0
		bar.W = skip
		bar.H = screenHeight
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&bar)

		xpos = (xpos + skip) % screenWidth

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				running = false
			}
		}

		end = sdl.GetPerformanceCounter() - origin
		dt = end - start
		// sdl.Delay(0.016 - dt/timerResolution)
		timings = append(timings, toMs(dt))
		renderer.Present()

	}

	/* for i := 1; i < len(timings); i++ {
		fmt.Println(float64(startTimes[i])/float64(timerResolution),
			float64(startTimes[i]-startTimes[i-1])/float64(timerResolution),
			float64(timings[i])/float64(timerResolution))
	}
	*/

	fmt.Println(hist(timings, 10, 1.0/60.0))
	return 0
}

func hist(timings []float64, nbins int, binsize float64) []uint {
	hist := make([]uint, nbins)
	for _, v := range timings {
		for i := 0; i < nbins; i++ {
			if v < float64(i+1)*binsize {
				hist[i]++
				break
			}
		}
	}
	return hist
}

func main() {
	os.Exit(run())
}
