// Play a sound file
// Time-stamp: <2022-10-04 18:21:44 christophe@pallier.org>
package main

import (
	"os"
	"log"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/mix"
)



func OpenAudioSystem() (err error) {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		return err
	}

	if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 1, 2048); err != nil {
		return err
	}

	return nil
}

func CloseAudioSystem() {
	mix.CloseAudio()
	sdl.Quit()
}

func PlayWavFile(soundFileName string) (err error) {
	chunk, err := mix.LoadWAV(soundFileName)
	if err != nil {
		return err
	}
	defer chunk.Free()

	chunk.Play(1, 0)

	// Wait until the end of sound playing
	for mix.Playing(-1) == 1 {
		sdl.Delay(16)
	}

	return nil
}


func main() {
	if len(os.Args) < 2 {
		log.Println("expected a sound file name as argument")
		os.Exit(1)
	}

	soundFileName := os.Args[1]

	if err := OpenAudioSystem(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer CloseAudioSystem()
	
	if err := PlayWavFile(soundFileName); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
