package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	fontPath        = "Inconsolata.ttf"
	fontSize        = 40
	maxResponseTime = 2000
)

type Trial struct {
	item   string
	isWord bool
}

var list = []Trial{
	Trial{item: "bonjour", isWord: true},
	Trial{item: "galaxie", isWord: true},
	Trial{item: "crocodile", isWord: true},
	Trial{item: "postral", isWord: false},
	Trial{item: "rontole", isWord: true},
	Trial{item: "callatrie", isWord: true},
}

func run() (err error) {
	var window *sdl.Window
	var font *ttf.Font
	var surface *sdl.Surface
	var text *sdl.Surface

	if err = ttf.Init(); err != nil {
		return
	}
	defer ttf.Quit()

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return
	}
	defer sdl.Quit()

	// Create a window for us to draw the text on
	if window, err = sdl.CreateWindow("Drawing text", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN); err != nil {
		return
	}
	defer window.Destroy()

	if surface, err = window.GetSurface(); err != nil {
		return
	}

	// Load the font for our text
	if font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		return
	}
	defer font.Close()

	// randomize the order of stimuli
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})

	var running bool = true // will get false if Quit.event is received

	for _, trial := range list {
		// clear screen
		surface.FillRect(nil, 0x00000000)
		window.UpdateSurface()

		// prepare the stimulus
		if text, err = font.RenderUTF8Blended(trial.item, sdl.Color{R: 255, G: 255, B: 255, A: 255}); err != nil {
			return
		}
		defer text.Free()

		// Draw the text around the center of the window
		if err = text.Blit(nil, surface, &sdl.Rect{X: 400 - (text.W / 2), Y: 300 - (text.H / 2), W: 0, H: 0}); err != nil {
			return
		}

		time.Sleep(2 * time.Second)
		// Update the window surface with what we have drawn
		window.UpdateSurface()

		// measure response & reaction time

		// flush event queue
		sdl.FlushEvents(sdl.FIRSTEVENT, sdl.LASTEVENT)

		var rt uint32 = 0
		var rt2 time.Duration
		var response sdl.Keycode

		start := time.Now()
		deadline := start.Add(2 * time.Second)

		for {
			if time.Now().After(deadline) || !running || rt != 0 {
				break
			}
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.QuitEvent:
					running = false
				case *sdl.KeyboardEvent:
					if t.Type == sdl.KEYDOWN {
						rt = t.Timestamp
						rt2 = time.Since(start)
						response = sdl.GetKeyFromScancode(t.Keysym.Scancode)
					}
				}
			}
			time.Sleep(5 * time.Millisecond)

		}

		fmt.Println(trial.item, response, rt, rt2)

		if !running {
			break
		}

	}

	return
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
