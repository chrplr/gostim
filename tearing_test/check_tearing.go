// Displays a moving vertical bar between two horizontal flashing bar
// Author: Christophe Pallier <christophe@pallier.org>
// Licence: GPL-3

package main

import (
	"flag"
	"fmt"
	"log"

	//	"math"
	"os"
	"runtime"
	"slices"

	"github.com/jaypipes/ghw"
	"github.com/shirou/gopsutil/cpu"
	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

var window *sdl.Window
var renderer *sdl.Renderer
var timerResolution float64
var screenWidth, screenHeight int32

func toMilliseconds(ticks uint64) float64 {
	return 1000.0 * float64(ticks) / timerResolution
}

type SysInfo struct {
	HostName string
	OS       string
	CPU      string
	GPU      []string
}

func GetSystemInfo() SysInfo {
	info := SysInfo{
		HostName: "",
		OS:       "",
		CPU:      "",
		GPU:      []string{},
	}

	var err error

	info.HostName, err = os.Hostname()
	if err != nil {
		log.Printf("Error getting hostname: %v", err)
	}

	info.OS = runtime.GOOS

	cpu, err := cpu.Info()
	if err != nil {
		log.Printf("Error getting CPU info: %v", err)
	}

	info.CPU = cpu[0].ModelName

	gpu, err := ghw.GPU()
	if err != nil {
		log.Printf("Error getting GPU info: %v", err)
	}

	info.GPU = make([]string, len(gpu.GraphicsCards))
	for _, card := range gpu.GraphicsCards {
		info.GPU = append(info.GPU, fmt.Sprintf(" %v\n", card))
	}

	return info
}

type SdlInfo struct {
	Platform       string
	RAM            int
	Version        sdl.Version
	Revision       string
	AudioDriver    string
	VideoDriver    string
	NDisplays      int
	Displays       string
	GLSwapInterval int
	VBlankEnvVar   string
	GLVBLank       string
}

func GetSDLInfo() SdlInfo {
	info := SdlInfo{
		Platform:       "",
		RAM:            0,
		Version:        sdl.Version{},
		Revision:       "",
		AudioDriver:    "",
		VideoDriver:    "",
		NDisplays:      0,
		Displays:       "",
		GLSwapInterval: 0,
		VBlankEnvVar:   "",
		GLVBLank:       "",
	}

	var err error

	info.Platform = sdl.GetPlatform()
	info.Revision = sdl.GetRevision()
	sdl.GetVersion(&info.Version)
	info.RAM = sdl.GetSystemRAM()

	info.AudioDriver = sdl.GetCurrentAudioDriver()

	if info.VideoDriver, err = sdl.GetCurrentVideoDriver(); err != nil {
		panic(err)
	}

	// EnvVarsOfInterest = []string{"NvOptimusEnablement", "__NV_PRIME_RENDER_OFFLOAD", "SDL_VIDEODRIVER", "SDL_AUDIODRIVER", "vblank_mode", "__GL_SYNC_TO_VBLANK"}

	/* XDG_SESSION_TYPE ; XDG_SESSION_DESKTOP; DESKTOP_SESSION ; DRI_PRIME=1 ; __NV_PRIME_RENDER_OFFLOAD=1 ; __GLX_VENDOR_LIBRARY_NAME=nvidia */
	/* NvOptimusEnablement: optimus prime on intel/nvidia hybrid systems */
	/* SDL_VIDEODRIVER, SDL_AUDIODRIVER (see https://wiki.libsdl.org/SDL2/FAQUsingSDL) */
	/* vblank_sync (See https://www.opengl.org/wiki/Swap_Interval)

	   0. Never synchronize with vertical refresh, ignore application's choice
	   1. Initial swap interval 0, obey application's choice
	   2. Initial swap interval 1, obey application's choice
	   3. Always synchronize with vertical refresh application chooses the minimum swap interval
	   4. Adaptative vsync

	*/

	info.VBlankEnvVar = os.Getenv("vblank_mode")

	/*  GL SYNC TO VBLANK  (https://download.nvidia.com/XFree86/Linux-x86_64/304.137/README/openglenvvariables.html)  NVidia specific?

	The __GL_SYNC_TO_VBLANK (boolean) environment variable can be used to control whether swaps are synchronized to a display device's vertical refresh:

	0. allows glXSwapBuffers to swap without waiting for vblank.
	1. forces glXSwapBuffers to synchronize with the vertical blanking period. This is the default behavior.
	*/

	info.GLVBLank = os.Getenv("__GL_SYNC_TO_VBLANK")

	info.NDisplays, _ = sdl.GetNumVideoDisplays()

	// func GetDisplayDPI(displayIndex int) (ddpi, hdpi, vdpi float32, err error)

	//DisplayModes := make([]int, info.NDisplays)
	for displayIndex := 0; displayIndex < info.NDisplays; displayIndex++ {
		name, _ := sdl.GetDisplayName(displayIndex)
		info.Displays += fmt.Sprintf("Display:%s\n", name)
		mode, _ := sdl.GetCurrentDisplayMode(displayIndex)
		info.Displays += fmt.Sprintf("   Resolution: %dx%d\n", mode.W, mode.H)
		info.Displays += fmt.Sprintf("   Refresh: %dHz\n", mode.RefreshRate)
		info.Displays += fmt.Sprintf("   Format: %d (see https://wiki.libsdl.org/SDL_PixelFormatEnum)\n", mode.Format)
	}

	info.GLSwapInterval, _ = sdl.GLGetSwapInterval()
	return info
}

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

// SineWave generates an array containg a sine wave; freq & samplingFreq in Hz, duration in ms
/* func SineWave(freq int, duration int, amplitude int, samplingFreq int) []int16 {
	length := (duration * samplingFreq) / 1000
	signal := make([]int16, length)
	amplitudef := float64(amplitude)
	for i := 0; i < length; i++ {
		signal[i] = int16(amplitudef * math.Sin(2.0*math.Pi*float64(i)/float64(freq)))
	}
	return signal
}
*/

func AVtest() (timings []float64) {
	testSound, err := mix.LoadWAV("tone440.wav")
	if err != nil {
		log.Fatalln(err)
	}
	defer testSound.Free()

	rectTop := sdl.Rect{X: 0, Y: 0, W: screenWidth, H: 200}
	rectBottom := sdl.Rect{X: 0, Y: screenHeight - 200, W: screenWidth, H: 200}

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
		visibleRect = loop%60 < 30
		if visibleRect {
			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.FillRect(&rectTop)
			renderer.FillRect(&rectBottom)
		} else {
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.FillRect(&rectTop)
			renderer.FillRect(&rectBottom)
		}

		renderer.Present()
		delta := toMilliseconds(sdl.GetPerformanceCounter() - start)
		timings = append(timings, delta)

		if loop%30 == 0 {
			testSound.Play(-1, 0)
		}

		// Process key press
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				running = false
			}
		}
	}
	return timings[1:]
}

func run() (timings []float64) {

	soundOn := true
	// preloading sound
	testSound, err := mix.LoadWAV("tone440.wav")
	if err != nil {
		log.Println(err)
		soundOn = false
	}
	defer testSound.Free()

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
		bar.W = 8
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&bar)
		if (xpos+skip > 0) && (xpos+skip < screenWidth) {
			xpos = xpos + skip
		} else { // invert direction of movement & play tone
			skip = -skip
			if soundOn {
				if _, err := testSound.Play(-1, 0); err != nil {
					log.Println(err)
				}
			}
		}

		// Process key press
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					running = false
				case sdl.K_RIGHT:
					skip = skip + 8
					log.Printf("skip = %d W=%d\n", skip, bar.W)
				case sdl.K_LEFT:
					if Abs(skip) >= 16 {
						skip = skip - 8
						log.Printf("skip = %d W=%d\n", skip, bar.W)
					}
				}
			}
		}

		sdl.Delay(3)
		renderer.Present()

		delta := toMilliseconds(sdl.GetPerformanceCounter() - start)
		timings = append(timings, delta)
	}
	return timings[2:]
}

func run_fixedvelocity() (timings []float64) {

	soundOn := true
	// preloading sound
	testSound, err := mix.LoadWAV("tone440.wav")
	if err != nil {
		log.Println(err)
		soundOn = false
	}
	defer testSound.Free()

	rectTop := sdl.Rect{X: 0, Y: 0, W: screenWidth, H: 200}
	rectBottom := sdl.Rect{X: 0, Y: screenHeight - 200, W: screenWidth, H: 200}

	bar := sdl.Rect{X: 0, Y: 250, W: 8, H: screenHeight - 500}
	var xpos int32 = 0
	var skip int32 = 8

	loop := 0
	running := true
	visibleRect := true
	ticksLastFrame := sdl.GetPerformanceCounter()

	for running {

		loop++

		ticksNow := sdl.GetPerformanceCounter()
		timeStep := toMilliseconds(ticksNow - ticksLastFrame)
		ticksLastFrame = ticksNow

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
		bar.W = 8
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&bar)
		if (xpos+skip > 0) && (xpos+skip < screenWidth) {
			xpos = xpos + skip*int32(timeStep)
		} else { // invert direction of movement & play tone
			skip = -skip
			if soundOn {
				if _, err := testSound.Play(-1, 0); err != nil {
					log.Println(err)
				}
			}
		}

		// Process key press
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				switch t.Keysym.Sym {
				case sdl.K_ESCAPE:
					running = false
				case sdl.K_RIGHT:
					skip = skip + 8
					log.Printf("skip = %d W=%d\n", skip, bar.W)
				case sdl.K_LEFT:
					if Abs(skip) >= 16 {
						skip = skip - 8
						log.Printf("skip = %d W=%d\n", skip, bar.W)
					}
				}
			}
		}

		sdl.Delay(3)
		renderer.Present()

		delta := toMilliseconds(sdl.GetPerformanceCounter() - ticksLastFrame)
		timings = append(timings, delta)
	}
	return timings[2:]
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

/* func isMouseMotionEvent(event sdl.Event, userdata interface{}) bool {
	return event.GetType() == sdl.MOUSEMOTION
}
*/

func main() {
	DesktopMode := flag.Bool("desktop", true, "true: FULLSCREEN_DESKTOP; false: FULLSCREEN")
	var err error

	flag.Parse()
	fmt.Printf("%+V\n", DesktopMode)

	fmt.Printf("System = %+v\n", GetSystemInfo())

	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalln(err)
	}
	defer sdl.Quit()

	if *DesktopMode {
		if window, err = sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 0, 0, sdl.WINDOW_FULLSCREEN_DESKTOP); err != nil {
			log.Fatalln(err)
		}
	} else {
		if window, err = sdl.CreateWindow("", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1920, 1080, sdl.WINDOW_FULLSCREEN); err != nil {
			log.Fatalln(err)
		}
	}
	defer window.Destroy()

	if renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC); err != nil {
		log.Fatalln(err)
	}
	defer renderer.Destroy()

	// Initialize mixer
	if err = mix.Init(mix.INIT_MP3); err != nil {
		log.Fatalln(err)
	}
	defer mix.Quit()

	// Open default playback device
	if err := mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT, mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE); err != nil {
		log.Fatalln(err)
	}
	defer mix.CloseAudio()

	sdl.DisableScreenSaver()
	sdl.ShowCursor(0)
	// sdl.SetEventFilter(isMouseMotionEvent, nil)

	fmt.Printf("SDL2 = %+v\n", GetSDLInfo())
	sdl.Delay(3000)

	screenWidth, screenHeight, _ = renderer.GetOutputSize()
	timerResolution = float64(sdl.GetPerformanceFrequency())

	timings := run()
	//timings := AVtest()
	printHistogram(timings)

}
