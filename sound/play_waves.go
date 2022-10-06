// Read wav filenames from standard input and plays them
// Time-stamp: <2022-10-05 11:03:39 christophe@pallier.org>
package main

import (
	"io"
	"bufio"
	"os"
	"fmt"
	"log"
	"errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/mix"
)


var (
	sounds  = make(map[string](*mix.Chunk))
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



// reads a text file line by line and returns the lines in a array
func ReadLines(textReader io.Reader) ([]string, error) {
	var lines []string
	
        reader := bufio.NewReader(textReader)
        for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		
		if err != nil {
			return nil, err
		}

		lines = append(lines, string(line))
        }
	return lines, nil
}


func CheckFileExists(fname string) error {
	file, err := os.Open(fname)
	if errors.Is(err, os.ErrNotExist) {
		return err
	}
 	if err != nil {
 		return err
 	}
	file.Close()
	return nil
}

// load wav files into the global variable "sounds" 
func PreloadWavfiles(fnames []string) error {

 	for _, fname := range fnames {
		chunk, err := mix.LoadWAV(fname)
		if err != nil {
			return err
		}
		sounds[fname] = chunk
 	}

 	return nil
}


func PlaySound(name string) (error) {
	_, err := sounds[name].Play(1, 0)
	if err != nil {
		return err
	}

	// Wait until the end of sound playing
	for mix.Playing(-1) == 1 {
		sdl.Delay(16)
	}

	return nil
}


func main() {
	// setup audio 
	if err := OpenAudioSystem(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer CloseAudioSystem()

	// read wav files names from standard input
	fnames, err := ReadLines(os.Stdin)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// check that files exist
	for _, fname := range fnames {
		if err := CheckFileExists(fname); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}

	// play files sequentially
	// for i, fname := range fnames {
	// 	fmt.Println(i, fname)
	// 	if err := PlayWavFile(fname); err != nil {
	// 		log.Println(err)
	// 		os.Exit(1)
	// 	}
	// }

	PreloadWavfiles(fnames)
	
	for i, name := range fnames {
		fmt.Println(i, name)
		err := PlaySound(name)
		if err != nil {
			log.Println(err)
		}
	}
}
