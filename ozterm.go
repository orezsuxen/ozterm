package main

// packages:
// main -> management of other modules
// insignals -> reads input, matches keyboard layout and sends command and text to term and gui
// gui -> display of a bunch of chars(GfxRunes); images ...
// term -> part, deciding what and where to render the text => panel management ... cursor storage, scrollback ...
//
//
//?! some sort of configuration management module
//? maybe put execution of programs in its own module ... eg. runner ?
//? maybe some audio handling for termbell...
//? session -> a group and screens

import (
	"fmt"
	"log"
	"time"

	"local/ozterm/dog"
	"local/ozterm/gui"
	"local/ozterm/insignals"
	"local/ozterm/runner"
	"local/ozterm/term"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	//LINE: setup SDL

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// keypress handling
	ins := insignals.New()

	//Window and gui
	gui, err := gui.New()
	if err != nil {
		panic(err)
	}
	defer gui.Close()

	//panel managment
	// panel := term.TestPanel()

	//TEST: Logging
	var dogger dog.Dog
	dog.Init(&dogger)

	//LINE: start prog
	ended := make(chan bool)
	outputChan := make(chan byte, 1024)
	errChan := make(chan byte, 1024)
	go runner.RunProg3(ended, outputChan, errChan, nil)

	var parsy term.ByteParser
	var parsyerr term.ByteParser

	outString := ""
	gotOut := false
	errString := ""
	gotErr := false

	contend := make([]string, 0, 100)
	newContend := false

	framecount := 0

	done := false
	//REM: Main Loop =========================================================
	log.Println(">>> main loop start")
	for !done {
		framecount += 1
		select {

		case <-ended:
			done = true
			continue

		case b := <-outputChan:
			// fmt.Print(">>> something on outchan:")
			// fmt.Print(b)
			// fmt.Fprint(os.Stdout, string(b))
			outRune := parsy.Parse(b)
			outString += string(outRune)
			gotOut = true
			continue

		case c := <-errChan:
			// fmt.Print(">>> something on errchan:")
			// fmt.Print(c)
			// fmt.Fprint(os.Stdout, string(c))
			errRune := parsyerr.Parse(c)
			errString += string(errRune)
			gotErr = true
			continue

		default:

			ins.Update()
			if ins.Quit || ins.KeysState[sdl.SCANCODE_ESCAPE] {
				done = true
				continue
			}

			gui.Clear()

			if gotOut {
				contend = append(contend, outString)
				gotOut = false
				newContend = true
			}
			if gotErr {
				contend = append(contend, errString)
				gotErr = false
				newContend = true
			}

			if len(contend) > 0 && newContend {
				fmt.Print("\n >>> ##########################################################################\n")
				for _, s := range contend {
					// fmt.Print("\n", n, " >>> ==========================================================\n")
					fmt.Print(s)
				}
				newContend = false
			}

			gui.DisplayTextAtPos(fmt.Sprint(framecount), 0, 0)

			if len(contend) > 1 {
				gui.DisplayTextLinesAtPos(contend[len(contend)-1], 0, 50, 0)
			}

			gui.Present()

			time.Sleep(time.Millisecond * 100)
		}
	} // main loop end
	log.Print(">>> main loop end")
}
