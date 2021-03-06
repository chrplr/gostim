// author: Jacky Boen

package main

import (
	"fmt"
	"os"
	"github.com/veandco/go-sdl2/sdl"
)


func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var screen_width, screen_height int32
	var rect sdl.Rect
	var xpos int32
	var skip int32
	var timings []uint64

	window, err := sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		0, 0, sdl.WINDOW_FULLSCREEN_DESKTOP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED | sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	defer renderer.Destroy()

	screen_width, screen_height, _  = renderer.GetOutputSize()

	xpos = 0
	skip = 8
	
	running := true

	sdl.Delay(1000)

	
	for running {

		renderer.SetDrawColor(0, 0, 0, 255) // clear screen
		renderer.Clear()

		rect.X = xpos
		rect.Y = 0
		rect.W = skip
		rect.H = screen_height

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&rect)

		renderer.Present()
		// sdl.Delay(16)
		xpos = (xpos + skip) % screen_width
		timings = append(timings, sdl.GetPerformanceCounter())

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				running = false
			}
		}

	}

	h := float64(sdl.GetPerformanceFrequency())
	for i:=1; i < len(timings); i++ {
		fmt.Println(float64(timings [i]-timings[i - 1]) / h)
	}
	return 0
}

func main() {
	os.Exit(run())
}
