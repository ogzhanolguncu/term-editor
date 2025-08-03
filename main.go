package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

const (
	BACKSPACE = 127

	ArrowLeft = 1000 + iota
	ArrowRight
	ArrowUp
	ArrowDown
	DelKey
	HomeKey
	EndKey
	PageUp
	PageDown
	EnterKey
)

var (
	cX, cY int = 1, 1
	buffer *TextBuffer
	screen *Screen
)

func main() {
	// Switch to raw term mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	buffer = NewTextBuffer()
	screen = NewScreen(buffer)

	for {
		ch := readScreenInput()
		// TODO: CTRL-C to exit. Move to this to somewhere else later
		if ch == 3 {
			break
		}
		if ch == ArrowLeft ||
			ch == ArrowRight ||
			ch == ArrowUp ||
			ch == ArrowDown ||
			ch == EnterKey {
			moveCursor(ch)
		} else {
			insertChart(byte(ch))
		}
		screen.Refresh()

	}
}

func moveCursor(key int) {
	switch key {
	case ArrowLeft:
		if cX > 1 {
			cX--
		} else if cY > 1 {
			// Wrap to end of previous line
			cY--
			cX = len(buffer.lines[cY-1]) + 1
		}
	case ArrowRight:
		if cY-1 < len(buffer.lines) {
			line := buffer.lines[cY-1]
			if cX <= len(line) {
				cX++
			} else if cY < len(buffer.lines) {
				// Wrap to start of next line
				cY++
				cX = 1
			}
		}
	case ArrowUp:
		if cY > 1 {
			cY--
			line := buffer.lines[cY-1]
			// Clamp cursor to avoid going past line end
			if cX-1 > len(line) {
				cX = len(line) + 1
			}
		}
	case ArrowDown:
		if cY < len(buffer.lines) {
			cY++
			if cY-1 < len(buffer.lines) {
				line := buffer.lines[cY-1]
				// Clamp cursor to avoid going past line end
				if cX-1 > len(line) {
					cX = len(line) + 1
				}
			}
		}
	case EnterKey:
		// Ensure buffer has current line
		for len(buffer.lines) <= cY-1 {
			buffer.lines = append(buffer.lines, "")
		}

		line := buffer.lines[cY-1]
		leftPart := line[:cX-1]  // Before cursor
		rightPart := line[cX-1:] // After cursor

		buffer.lines[cY-1] = leftPart

		// Insert new line with right part
		newLine := []string{rightPart}
		buffer.lines = append(buffer.lines[:cY], append(newLine, buffer.lines[cY:]...)...)

		cY++
		cX = 1
	}
	// Sync terminal cursor with logical position
	fmt.Printf("\x1b[%d;%dH", cY, cX)
}

func readScreenInput() int {
	var buffer [1]byte
	if _, err := os.Stdin.Read(buffer[:]); err != nil {
		nuke(err)
	}

	if buffer[0] == '\n' || buffer[0] == '\r' {
		return EnterKey
	}

	if buffer[0] == '\x1b' {
		var seq [2]byte
		if cc, err := os.Stdin.Read(seq[:]); cc != 2 || err != nil {
			return '\x1b'
		}

		switch seq[0] {

		case '[':
			if seq[1] >= '0' && seq[1] <= '9' {
				if cc, err := os.Stdin.Read(buffer[:]); cc != 1 || err != nil {
					return '\x1b'
				}
				if buffer[0] == '~' {
					switch seq[1] {
					case '1':
						return HomeKey
					case '3':
						return DelKey
					case '4':
						return EndKey
					case '5':
						return PageUp
					case '6':
						return PageDown
					case '7':
						return HomeKey
					case '8':
						return EndKey
					}
				}
			} else {
				switch seq[1] {
				case 'A':
					return ArrowUp
				case 'B':
					return ArrowDown
				case 'C':
					return ArrowRight
				case 'D':
					return ArrowLeft
				case 'H':
					return HomeKey
				case 'F':
					return EndKey
				}
			}
		case '0':
			switch seq[1] {
			case 'H':
				return HomeKey
			case 'F':
				return EndKey
			}
		}

		return '\x1b'
	}
	return int(buffer[0])
}

func nuke(err error) {
	screen.Clear()
	log.Fatal(err)
}

// ##################### TEXT BUFFER #####################

type TextBuffer struct {
	lines []string
}

func NewTextBuffer() *TextBuffer {
	return &TextBuffer{
		lines: make([]string, 0),
	}
}

func insertChart(ch byte) {
	for len(buffer.lines) <= cY-1 {
		buffer.lines = append(buffer.lines, "")
	}

	line := buffer.lines[cY-1]
	// This is required for inline editing like editing an item from middle.
	line = line[:cX-1] + string(ch) + line[cX-1:]
	buffer.lines[cY-1] = line
	cX++
}

// ##################### SCREEN #####################

type Screen struct {
	buffer                   *TextBuffer
	dirty                    bool
	dirtyLines               map[int]bool
	lastCursorX, lastCursorY int
}

func NewScreen(textBuffer *TextBuffer) *Screen {
	screen := &Screen{
		buffer:      textBuffer,
		lastCursorX: cX,
		lastCursorY: cY,
	}
	screen.Clear()
	return screen
}

func (s *Screen) Clear() {
	fmt.Print("\x1b[2J") // Clear entire screen
	fmt.Print("\x1b[H")  // Move cursor to top-left (1,1)
}

func (s *Screen) Refresh() {
	ab := bytes.NewBufferString("\x1b[?25l") // Hide cursor
	ab.WriteString("\x1b[H")                 // Move cursor to top-left (1,1)

	// Draw all lines
	for i, line := range s.buffer.lines {
		fmt.Fprintf(ab, "\x1b[%d;1H\x1b[K%s", i+1, line) // \x1b[K clears to end of line
	}

	// Position cursor and show it
	fmt.Fprintf(ab, "\x1b[%d;%dH", cY, cX)
	ab.WriteString("\x1b[?25h") // Show cursor

	ab.WriteTo(os.Stdout)
}
