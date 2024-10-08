package insignals

import (
	// "fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Insignals struct {
	Quit      bool
	KeysState map[int]bool
	Input     []int
}

func New() Insignals {
	return Insignals{
		Quit:      false,
		KeysState: make(map[int]bool),
		Input:     make([]int, 0),
	}
}

func (i *Insignals) Update() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			i.Quit = true
		case *sdl.KeyboardEvent:
			k := event.(*sdl.KeyboardEvent)
			if k.State == sdl.PRESSED && k.Repeat == 0 {
				i.KeysState[int(k.Keysym.Scancode)] = true
				i.Input = append(i.Input, int(k.Keysym.Sym))
			}
			if k.State == sdl.RELEASED && k.Repeat == 0 {
				i.KeysState[int(k.Keysym.Scancode)] = false
			}
		}
	}
}

func (i *Insignals) Drain() []int {
	retval := i.Input
	i.Input = make([]int, 0)
	return retval

}
