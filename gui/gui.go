package gui

//NOTE: terminology
// panel -> a text rectangle containing the output of a running program (hopefully shells)
// screen -> a view of one or more panels of a single group in splits
// group -> a collection of panels running programs
//
// infoline -> a status bar containing info of the current screen context(the Tab bar if you will)
// globalscreen -> the current screen and infoline (everything that needs to be renderd by gui)

import (
	"fmt"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func GuiTest() {
	fmt.Println("Message from gui")
}

// WIP Params
const charHeight = 15
const charWidth = 8

const windowHeight = 800
const windowWidth = 600

const fontPointSize = 12

// Screen ==================================================================
type Screen struct {
	id uint64
	// how to display a group
	// some floating ?

}

// Window ==================================================================
type Window struct {
	id       uint64
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Font     *ttf.Font
}

// name, dimension, maybe pos as params
func New() (term Window, err error) {
	var retval Window

	win, err := sdl.CreateWindow(
		"the best window",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowHeight, windowWidth,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return retval, err
	}

	rend, err := sdl.CreateRenderer(win, -1, 2)
	if err != nil {
		return retval, err
	}

	err = ttf.Init()
	if err != nil {
		return retval, err
	}

	fonty, err := ttf.OpenFont("./FiraCodeNerdFont-Regular.ttf", fontPointSize)
	if err != nil {
		return retval, err
	}

	retval.id = 1
	retval.Window = win
	retval.Renderer = rend
	retval.Font = fonty
	return retval, nil
}

func (tg *Window) Close() {

	tg.Font.Close()
	ttf.Quit()
	tg.Renderer.Destroy()
	tg.Window.Destroy()
}

func (tg *Window) Clear() {
	tg.Renderer.SetDrawColor(80, 80, 80, 255)
	tg.Renderer.Clear()
	tg.Renderer.SetDrawColor(255, 0, 255, 255)
}

func (tg *Window) Present() {
	tg.Renderer.Present()
}

func (tg *Window) DisplayText(text string) {

	var fg sdl.Color = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	rect2 := sdl.Rect{X: 0, Y: 0, W: int32(len(text) * charWidth), H: charHeight}

	msgsurface, err := tg.Font.RenderUTF8Blended(text, fg)
	if err != nil {
		panic(err)
	}

	msgtexture, err := tg.Renderer.CreateTextureFromSurface(msgsurface)
	if err != nil {
		panic(err)
	}

	msgsurface.Free()

	tg.Renderer.Copy(msgtexture, nil, &rect2)

	msgtexture.Destroy()
}

func (tg *Window) DisplayTextAtPos(text string, x int, y int) {

	var fg sdl.Color = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	rect2 := sdl.Rect{X: int32(x), Y: int32(y), W: int32(len(text) * charWidth), H: charHeight}

	msgsurface, err := tg.Font.RenderUTF8Blended(text, fg)
	if err != nil {
		panic(err)
	}
	defer msgsurface.Free()

	msgtexture, err := tg.Renderer.CreateTextureFromSurface(msgsurface)
	if err != nil {
		panic(err)
	}
	defer msgtexture.Destroy()

	tg.Renderer.Copy(msgtexture, nil, &rect2)
}

// display multiple lines at position
func (w *Window) DisplayTextLinesAtPos(text string, x int, y int, s int) { // s = spacing
	lines := strings.Split(text, "\n")

	var fg sdl.Color = sdl.Color{R: 255, G: 255, B: 255, A: 255}

	for i, l := range lines {
		rect2 := sdl.Rect{X: int32(x), Y: int32(y + (i*charHeight + s)), W: int32(len(l) * charWidth), H: charHeight}

		msgsurface, err := w.Font.RenderUTF8Blended(l, fg)
		if err != nil {
			panic(err)
		}
		defer msgsurface.Free()

		msgtexture, err := w.Renderer.CreateTextureFromSurface(msgsurface)
		if err != nil {
			panic(err)
		}
		defer msgtexture.Destroy()

		w.Renderer.Copy(msgtexture, nil, &rect2)

	}

}
