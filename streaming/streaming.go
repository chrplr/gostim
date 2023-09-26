// Displays png files from ./images folder
// Time-stamp: <2022-07-20 20:11:57 christophe@pallier.org>
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	screen_height = 768
	screen_width  = 1024
	picsWidth     = 768
	picsHeight    = 768
)

// returns the list of png files in folder (alphabetically sorted by name)
func getPics(folder, pattern string) []string {
	files, err := filepath.Glob(filepath.Join(folder, pattern))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	sort.Strings(files)
	return files
}

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var screen_width, screen_height int32
	var r sdl.Rect
	var pics []*sdl.Texture
	var timings []uint64

	// initialisation
	window, err := sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screen_width, screen_height, sdl.WINDOW_FULLSCREEN)
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

	screen_width, screen_height, _ = renderer.GetOutputSize()
	log.Println(screen_height, screen_width)

	// load the images
	for _, f := range getPics("images/", "*.png") {
		fmt.Println(f)
		tex, err := img.LoadTexture(renderer, f)
		if err != nil {
			fmt.Printf("Error (%s) loading %s\n", img.GetError(), f)
			return 2
		}
		// We should check the size of the image
		pics = append(pics, tex)
	}

	// main loop (display)
	running := true
	for running {

		renderer.SetDrawColor(0, 0, 0, 255) // clear screen
		renderer.Clear()

		r.W = picsWidth
		r.H = picsHeight
		r.X = (screen_width - r.W) / 2
		r.Y = (screen_height - r.H) / 2

		for _, tex := range pics {
			renderer.Copy(tex, nil, &r)
			// in case variable size images, we would need:
			// _, _, picWidth, picHeight, err := tex.Query()
			//
			renderer.Present()
			sdl.Delay(50)
			timings = append(timings, sdl.GetPerformanceCounter())
		}

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
	for i := 1; i < len(timings); i++ {
		fmt.Println(float64(timings[i]-timings[i-1]) / h)
	}
	return 0
}

func main() {
	os.Exit(run())
}
