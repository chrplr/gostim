package main


import (
	"fmt"
	"os"
	"strings"
	"sort"
	"io/ioutil"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
)

// const (
// 	screenWidth = 640
// 	screenHeight = 480
// )

// returns the list of png files in folder (alphabetically sorted by name)
func getPics(folder string) []string {
	var pics []string
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".png") {
			pics = append(pics, f.Name())
		}
	}
	sort.Strings(pics)
	return pics
}




func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var screen_width, screen_height int32
	var r sdl.Rect
	var pics []*sdl.Texture
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

	for _, f := range getPics("42_texture_streaming") {
		tex, err := img.LoadTexture(renderer, "42_texture_streaming/" + f)
		if err != nil {
			fmt.Printf("Error loading %s\n", f)
			return 2
		}
		pics = append(pics, tex) 
	}
	
	
	running := true

	sdl.Delay(1000)

	
	for running {

		renderer.SetDrawColor(0, 0, 0, 255) // clear screen
		renderer.Clear()

		r.X = screen_width / 2
		r.Y = screen_height / 2 
		r.W = 64
		r.H = 205

		for _, tex := range pics {
			renderer.Copy(tex, nil, &r)
			renderer.Present()
			sdl.Delay(100)
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
	for i:=1; i < len(timings); i++ {
		fmt.Println(float64(timings [i]-timings[i - 1]) / h)
	}
	return 0
}

func main() {
	os.Exit(run())
}
