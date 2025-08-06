package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

type Key int

const (
	// ASCII/Control codes with explicit values
	CtrlC     Key = 3
	Backspace Key = 127
	Enter     Key = 13
)

const (
	// Special keys using iota
	ArrowLeft Key = iota + 200
	ArrowRight
	ArrowUp
	ArrowDown
	Del
	Home
	End
	PageUp
	PageDown
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

	// Main loop
	for {
		key := readScreenInput()
		handleKey(key)
		screen.Refresh()
	}
}

// ##################### Key Handling #####################

func handleKey(key Key) {
	switch key {
	case CtrlC:
		screen.Clear()
		os.Exit(0)
	case Backspace:
		handleBackspace()
	case Enter, ArrowLeft, ArrowRight, ArrowUp, ArrowDown:
		handleCursorMove(key)
	default:
		handleCharInsert(key)
	}
}

func handleCursorMove(key Key) {
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
			if cX-1 < len(line) {
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
	case Enter:
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
}

func handleBackspace() {
	// Cursor is at beginning and if there is a line above
	if cX == 1 && cY > 1 {
		// Ensure buffer has enough lines
		for len(buffer.lines) <= cY-1 {
			buffer.lines = append(buffer.lines, "")
		}

		line := buffer.lines[cY-1]
		buffer.lines[cY-2] = fmt.Sprintf("%s%s", buffer.lines[cY-2], line)
		// Remove current line from the slice
		buffer.lines = append(buffer.lines[:cY-1], buffer.lines[cY:]...)
		// Move cursor to above line and end of it
		cY--
		cX = len(buffer.lines[cY-1]) + 1
	} else if cX > 1 {
		line := buffer.lines[cY-1]
		buffer.lines[cY-1] = fmt.Sprintf("%s%s", line[:cX-2], line[cX-1:])
		cX--
	} else {
		// Do nothing - can't backspace at start of file
		return
	}
}

func readScreenInput() Key {
	var buffer [1]byte
	if _, err := os.Stdin.Read(buffer[:]); err != nil {
		nuke(err)
	}

	if buffer[0] == '\n' || buffer[0] == '\r' {
		return Enter
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
						return Home
					case '3':
						return Del
					case '4':
						return End
					case '5':
						return PageUp
					case '6':
						return PageDown
					case '7':
						return Home
					case '8':
						return End
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
					return Home
				case 'F':
					return End
				}
			}
		case '0':
			switch seq[1] {
			case 'H':
				return Home
			case 'F':
				return End
			}
		}

		return '\x1b'
	}
	return Key(buffer[0])
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

func handleCharInsert(ch Key) {
	for len(buffer.lines) <= cY-1 {
		buffer.lines = append(buffer.lines, "")
	}

	line := buffer.lines[cY-1]
	// This is required for inline editing like editing an item from middle.
	line = line[:cX-1] + string(rune(ch)) + line[cX-1:]
	buffer.lines[cY-1] = line
	cX++
}

// ##################### SCREEN #####################

type Screen struct {
	buffer    *TextBuffer
	lastLines int
}

func NewScreen(textBuffer *TextBuffer) *Screen {
	screen := &Screen{
		buffer: textBuffer,
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

	// Calculate minimum lines to render: either up to cursor position or all buffer lines
	// This ensures cursor is always visible even if buffer is sparse
	maxLines := max(cY, len(s.buffer.lines))

	// Render all lines up to maxLines
	for i := range maxLines {
		line := ""
		// Safe access: use empty string if line doesn't exist in buffer
		if i < len(s.buffer.lines) {
			line = s.buffer.lines[i]
		}
		// Move to row i+1 col 1, clear line to end, then write content
		fmt.Fprintf(ab, "\x1b[%d;1H\x1b[K%s", i+1, line)
	}

	// Clear any leftover lines from previous render handles line deletion
	// If we previously rendered 5 lines but now only need 3, clear lines 4-5
	if s.lastLines > maxLines {
		for i := maxLines; i < s.lastLines; i++ {
			// Move to row i+1 col 1, clear entire line (remove old content)
			fmt.Fprintf(ab, "\x1b[%d;1H\x1b[K", i+1)
		}
	}

	// Move cursor to logical position
	fmt.Fprintf(ab, "\x1b[%d;%dH", cY, cX)

	// Remember how many lines we rendered for next refresh cycle
	s.lastLines = maxLines
	ab.WriteString("\x1b[?25h") // Show cursor
	ab.WriteTo(os.Stdout)
}
