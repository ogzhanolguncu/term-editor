package main

import (
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

var cX, cY int = 1, 1

func main() {
	// Switch to raw term mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	clearScreen()
	for {
		key := readScreenInput()
		// TODO: CTRL-C to exit. Move to this to somewhere else later
		if key == 3 {
			break
		}
		if key == ArrowLeft ||
			key == ArrowRight ||
			key == ArrowUp ||
			key == ArrowDown ||
			key == EnterKey {
			moveCursor(key)
		} else {
			cX++
			fmt.Printf("%c", key)
		}
	}
}

func clearScreen() {
	fmt.Print("\x1b[2J") // Clear entire screen
	fmt.Print("\x1b[H")  // Move cursor to top-left (1,1)
}

func moveCursor(key int) {
	switch key {
	case ArrowLeft:
		if cX > 1 {
			cX--
		}
	case ArrowRight:
		// TODO: Later add bounds to check
		cX++
	case ArrowUp:
		if cY > 1 {
			cY--
		}
	case ArrowDown:
		// TODO: Later add bounds to check
		cY++
	case EnterKey:
		cY++
		cX = 1 // Go back to beginning of line
	}
	fmt.Printf("\x1b[%d;%dH", cY, cX)
}

func hideCursor() {
	fmt.Print("\x1b[?25l")
}

func showCursor() {
	fmt.Print("\x1b[?25h")
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
	clearScreen()
	log.Fatal(err)
}
