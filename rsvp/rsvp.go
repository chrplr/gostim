// presentation of
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	onset    float64
	duration float64
	item     string
}

func newTrial(line string) Trial {
	var t Trial
	var err error

	tokens := strings.Split(line, "\t")
	//_, err := fmt.Sscanf(line, "%[^\t]\t%f\t%f", &t.item, &t.onset, &t.duration)

	t.item = tokens[0]

	t.onset, err = strconv.ParseFloat(tokens[1], 64)
	if err != nil {
		panic(err)
	}

	t.duration, err = strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		panic(err)
	}
	return t
}

func ReadListOfTrials(fname string) []Trial {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	// fileScanner.Scan() // skip the first line which contains the header (col names)

	list := make([]Trial, 0, 8000)

	for fileScanner.Scan() {
		list = append(list, newTrial(fileScanner.Text()))
	}

	return list
}

func WaitForKey() sdl.Keycode {
	sdl.FlushEvents(0x300, 0x400)
	var keyCode sdl.Keycode
	running := true
	for running {
		event := sdl.WaitEvent() // wait here until an event is in the event queue
		switch t := event.(type) {
		case sdl.QuitEvent:
			running = false
			keyCode = sdl.K_ESCAPE
		case sdl.KeyboardEvent:
			keyCode = t.Keysym.Sym
			running = false
		}
	}
	return keyCode
}

func WaitForMRISync() {
	running := true
	for running {
		running = WaitForKey() != sdl.K_t
	}
}

func run(tsvfile string) (err error) {
	var window *sdl.Window
	var font *ttf.Font
	var surface *sdl.Surface
	var text *sdl.Surface

	list := ReadListOfTrials(tsvfile)

	if err = ttf.Init(); err != nil {
		return
	}
	defer ttf.Quit()

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		return
	}
	defer sdl.Quit()

	// Create a window
	if window, err = sdl.CreateWindow("Drawing text", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN); err != nil {
		return
	}
	defer window.Destroy()

	if surface, err = window.GetSurface(); err != nil {
		return
	}

	// Load the font
	if font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		return
	}
	defer font.Close()

	var running bool = true // will get false if Quit.event is received

	start := time.Now()

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

		var delta time.Duration = time.Duration(trial.onset*1000.0*float64(time.Millisecond)) - (time.Since(start))
		fmt.Println(time.Since(start), delta, trial.onset*1000.0, time.Duration(trial.onset*1000.0*float64(time.Millisecond)))
		time.Sleep(delta)

		// Update the window surface with what we have drawn
		window.UpdateSurface()

		time.Sleep(time.Duration(trial.duration * 1000 * float64(time.Millisecond)))

		// Check if Quit button was pressed
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		if !running {
			break
		}

	}

	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./rsvp  wordlist.tsv\n\nwhere wolrdlist.tsv is a tab-separated file with three columns: Word, OnseTime, Duration")
		return
	}

	if err := run(os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
