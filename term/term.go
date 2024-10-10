package term

//NOTE: terminology
// panel -> a text rectangle containing the output of a running program (hopefully shells)
// screen -> a view of one or more panels of a single group in splits
// group -> a collection of panels running programs
//
// infoline -> a status bar containing info of the current screen context(the Tab bar if you will)
// globalscreen -> the current screen and infoline (everything that needs to be renderd by gui)

//NOTE: com between term and gui in form of a extended 'Rune' containing color and stuff ?

import (
	"fmt"
	"unicode/utf8"
)

func TermTest() {
	fmt.Println("Message from term")
}

// Point
type xyPoint struct {
	X int
	Y int
}

// GfxRune ============================================================
type GfxRune struct {
	char   rune
	color  int // int for now
	style  int // int for now
	effect int // int for now
}

// GfxLine ============================================================
type GfxLine struct {
	Line []GfxRune
}

func (g GfxLine) String() string {
	var ret string
	for _, r := range g.Line {
		ret += string(r.char)
	}
	return ret

}

func NewGfxLine() GfxLine {
	retval := make([]GfxRune, 0)
	return GfxLine{
		Line: retval,
	}
}
func NewGfxLineFromString(in string) GfxLine {
	retval := make([]GfxRune, 0)
	for _, r := range in {
		retval = append(retval, GfxRune{
			char: r,
		})
	}
	return GfxLine{
		Line: retval,
	}
}

// Group ================================================================
type Group struct {
	id uint64

	Panels     []Panel
	focusPanel uint64
}

// Panel ================================================================
type Panel struct {
	id uint64

	// scrollback ! []GfxLine ... 10000 ?
	dimension xyPoint
	mode      int
	wraping   bool

	cursorPos            xyPoint
	cursorPosStoreageDEC xyPoint
	cursorPosStorrageSEC xyPoint
	//
	Lines []GfxLine
}

// resizing func that takes callback from gui

// Parsing =================================================================
// for now just print the event ...
// parser struct for buffering bytes ?
type ByteParser struct {
	buffer []byte
}

func (bp *ByteParser) Parse(inByte byte) (r rune, ok bool) {

	// first time setup
	if bp.buffer == nil {
		bp.buffer = make([]byte, 0, 4)
	}
	bp.buffer = append(bp.buffer, inByte)

	// try parsing
	if len(bp.buffer) == 0 {
		fmt.Println(">>> error byte buffer still 0!")
	}

	// r, runeLen := ProcessBytes(bp.buffer)
	start := utf8.RuneStart(bp.buffer[0])
	if !start {
		fmt.Println(">>> pars error not start")
		return 0, false
	}
	r, l := utf8.DecodeRune(bp.buffer)
	if r == utf8.RuneError {
		// fmt.Println(">>> pars error RuneError:", r, l)
		return 0, false
	}

	// remove consumed bytes
	bp.buffer = bp.buffer[l:]
	return r, true

}

// func ProcessBytes(inBytes []byte) (r rune, runeLen int) {
// 	start := utf8.RuneStart(inBytes[0])
// 	if !start {
// 		fmt.Println(">>> pars error not start")
// 		return 0, 1
// 	}
// 	r, l := utf8.DecodeRune(inBytes)
// 	if r == utf8.RuneError {
// 		// fmt.Println(">>> pars error RuneError:", r, l)
// 		return 0, 0
// 	}
//
// 	// fmt.Print(string(r))
// 	return r, l
//
// }

// Test stuff ==============================================================
func TestPanel() Panel {
	s := make([]GfxLine, 0)
	s = append(s, NewGfxLineFromString("this is a test line"))
	s = append(s, NewGfxLineFromString("here is another line"))
	return Panel{
		id:    333,
		Lines: s,
	}
}

func TestGroup() Group {
	var retval Group
	retval.id = 22
	retval.Panels = make([]Panel, 0)
	retval.Panels = append(retval.Panels, TestPanel())

	return retval
}

func (p *Panel) Update() {

}

func (p *Panel) SetContend(in string) {
	p.Lines = make([]GfxLine, 0)
	p.Lines = append(p.Lines, NewGfxLineFromString(in))
}
