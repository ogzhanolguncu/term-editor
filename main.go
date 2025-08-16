package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/term"
)

// TODO:
// [x] Implement Ctrl+S save functionality with propmting at the bottom. Also required for later iterations
// [x] Add command line argument support for opening file
// [x] Add file loading (Ctrl+O) - load file content, clear buffer, reset cursor
// [x] Track modified state - bool flag, set on edits, clear on save/load
// [ ] Add basic status line - show "filename [modified] | Line X, Col Y" at bottom
// [ ] Implement Ctrl+N for new file
// [ ] Add Ctrl+Q quit with save prompt

const Version = "GICO 0.1"

type Key int

const (
	// ASCII/Control codes with explicit values
	CtrlC     Key = 3
	CtrlO     Key = 15
	CtrlS     Key = 19
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

	var fPath string
	if len(os.Args) > 1 {
		fPath = os.Args[1]
	}

	buffer = NewTextBuffer(strings.TrimSpace(fPath))
	screen = NewScreen(buffer)
	screen.Refresh()

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
	case CtrlS:
		handleSave()
	case CtrlO:
		handleLoadFile()
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
		line := buffer.lines[cY-1]
		leftPart := line[:cX-1]  // Before cursor
		rightPart := line[cX-1:] // After cursor

		buffer.lines[cY-1] = leftPart

		// Insert new line with right part
		newLine := []string{rightPart}
		buffer.lines = append(buffer.lines[:cY], append(newLine, buffer.lines[cY:]...)...)

		cY++
		cX = 1
		buffer.modified = true
	}
}

func handleBackspace() {
	// Cursor is at beginning and if there is a line above
	if cX == 1 && cY > 1 {
		line := buffer.lines[cY-1]
		buffer.lines[cY-2] = fmt.Sprintf("%s%s", buffer.lines[cY-2], line)
		// Remove current line from the slice
		buffer.lines = append(buffer.lines[:cY-1], buffer.lines[cY:]...)
		// Move cursor to above line and end of it
		cY--
		cX = len(buffer.lines[cY-1]) + 1
		buffer.modified = true
	} else if cX > 1 {
		line := buffer.lines[cY-1]
		buffer.lines[cY-1] = fmt.Sprintf("%s%s", line[:cX-2], line[cX-1:])
		cX--
		buffer.modified = true
	} else {
		// Do nothing - can't backspace at start of file
		return
	}
}

func handleSave() {
	fPath := buffer.fPath
	if buffer.fPath == "" {
		fPath = promptForFilename("Save as: ")
		// User cancelled
		if fPath == "" {
			return
		}
		buffer.fPath = fPath
	}

	if err := trySaveFile(fPath); err != nil {
		screen.SetPrompt(fmt.Sprintf("Error saving: %v", err))
		time.Sleep(1 * time.Second)
		screen.SetPrompt("")
		return
	}
	buffer.modified = false

	// Get file stats for display
	stat, err := os.Stat(fPath)
	if err != nil {
		screen.SetPrompt(fmt.Sprintf("Saved but unable to get file info: %v", err))
		return
	}

	displayPath := getDisplayPath(fPath)
	lineCount := len(buffer.lines)
	screen.SetPrompt(fmt.Sprintf("\"%s\" %dL, %dB", displayPath, lineCount, stat.Size()))
	time.Sleep(200 * time.Millisecond)
	screen.SetPrompt("")
	// When new file opens move cursor to beginning
	cX, cY = 1, 1
}

func handleLoadFile() {
	fPath := promptForFilename("Open file: ")
	// User cancelled
	if fPath == "" {
		return
	}
	buffer.fPath = fPath

	var lines []byte
	lines, err := tryToOpenFile(fPath)
	if err != nil {
		screen.SetPrompt(fmt.Sprintf("Error opening: %v", err))
		time.Sleep(1 * time.Second)
		screen.SetPrompt("")
		return
	}

	buffer.lines = strings.Split(string(lines), "\n")
	buffer.modified = false

	stat, err := os.Stat(fPath)
	if err != nil {
		screen.SetPrompt(fmt.Sprintf("Loaded but unable to get file info: %v", err))
		return
	}

	displayPath := getDisplayPath(fPath)
	lineCount := len(buffer.lines)
	screen.SetPrompt(fmt.Sprintf("\"%s\" %dL, %dB", displayPath, lineCount, stat.Size()))
	time.Sleep(200 * time.Millisecond)
	screen.SetPrompt("")
	screen.MoveCursorToStart()
}

func getDisplayPath(fPath string) string {
	// Convert to absolute path first
	absPath, err := filepath.Abs(fPath)
	if err != nil {
		return fPath // fallback to original
	}

	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return absPath // fallback to absolute path
	}

	// Replace home directory with ~
	if strings.HasPrefix(absPath, homeDir) {
		return "~" + absPath[len(homeDir):]
	}

	return absPath
}

func tryToOpenFile(filename string) ([]byte, error) {
	if filename == "" {
		return nil, nil
	}
	return os.ReadFile(filename)
}

func trySaveFile(filename string) error {
	if filename == "" {
		return nil
	}
	c := strings.Join(buffer.lines, "\n")
	return os.WriteFile(filename, []byte(c), 0644)
}

func promptForFilename(prompt string) string {
	screen.SetPrompt(prompt)
	var filename strings.Builder

	for {
		key := readScreenInput()
		switch key {
		case CtrlC:
			screen.Refresh()
			return ""
		case Enter:
			if filename.Len() == 0 {
				continue
			}
			screen.Refresh()
			return filename.String()
		case Backspace:
			if filename.Len() > 0 {
				str := filename.String()
				filename.Reset()
				// From beginning to last char
				filename.WriteString(str[:len(str)-1])
				screen.SetPrompt(prompt + filename.String())
			}
		default:
			if key >= 32 && key <= 126 {
				filename.WriteRune(rune(key))
				screen.SetPrompt(prompt + filename.String())
			}
		}
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
	lines    []string
	fPath    string
	modified bool
}

func NewTextBuffer(fPath string) *TextBuffer {
	tb := &TextBuffer{
		lines:    []string{""},
		fPath:    fPath,
		modified: false,
	}

	tb.openFile()
	return tb
}

func (tb *TextBuffer) openFile() {
	if tb.fPath == "" {
		return
	}

	content, err := os.ReadFile(tb.fPath)
	if err != nil {
		if os.IsNotExist(err) {
			cX, cY = 1, 1
			return
		}
		nuke(err)
	}
	tb.lines = strings.Split(strings.TrimSpace(string(content)), "\n")
	lastLine := len(tb.lines) - 1
	cX, cY = len(tb.lines[lastLine])+1, len(tb.lines)
}

func handleCharInsert(ch Key) {
	line := buffer.lines[cY-1]
	// This is required for inline editing like editing an item from middle.
	line = line[:cX-1] + string(rune(ch)) + line[cX-1:]
	buffer.lines[cY-1] = line
	buffer.modified = true
	cX++
}

// References for ANSI sequences
//
// \x1b[H          Move to top-left (1,1)
// \x1b[2J         Clear entire screen
// \x1b[K          Clear line from cursor right
// \x1b[2K         Clear entire line
// \x1b[%d;%dH     Move to row,col
// \x1b[%dA        Move up N lines
// \x1b[%dB        Move down N lines
// \x1b[?25l       Hide cursor
// \x1b[?25h       Show cursor
// \x1b[7m         Invert colors (highlight)
// \x1b[0m         Reset all formatting

// ##################### SCREEN #####################

type Screen struct {
	buffer    *TextBuffer
	lastLines int
	prompt    string
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

	// Get terminal width
	width := 80 // fallback
	if w, _, err := term.GetSize(int(os.Stdin.Fd())); err == nil {
		width = w
	}

	// Build 3-section header like UW PICO
	leftSection := fmt.Sprintf("  %s", Version)

	filename := "New Buffer"
	if s.buffer.fPath != "" {
		filename = getDisplayPath(s.buffer.fPath)
	}
	centerSection := fmt.Sprintf("File: %s", filename)

	rightSection := ""
	if s.buffer.modified {
		rightSection = fmt.Sprintf("%s  ", "MODIFIED")
	}

	// Calculate center position for middle section
	centerStart := (width - len(centerSection)) / 2

	// Build header with exact positioning
	header := make([]rune, width)
	for i := range header {
		header[i] = ' ' // Fill with spaces
	}

	// Place left section (starts at position 0)
	copy(header[0:], []rune(leftSection))

	// Place center section (centered)
	if centerStart >= 0 && centerStart+len(centerSection) <= width {
		copy(header[centerStart:], []rune(centerSection))
	}

	// Place right section (right-aligned)
	rightStart := width - len(rightSection)
	if rightStart >= 0 {
		copy(header[rightStart:], []rune(rightSection))
	}

	// Header with inverted colors covering full row
	fmt.Fprintf(ab, "\x1b[1;1H\x1b[7m%s\x1b[0m", string(header))

	// Render buffer content starting from line 2
	maxLines := max(cY, len(s.buffer.lines))
	for i := range maxLines {
		line := ""
		if i < len(s.buffer.lines) {
			line = s.buffer.lines[i]
		}

		// Current line highlighting + offset by 1 for header
		if i+1 == cY {
			// Pad line to full terminal width for complete background highlight
			paddedLine := fmt.Sprintf("%-*s", width, line)
			fmt.Fprintf(ab, "\x1b[%d;1H\x1b[K\x1b[48;5;236m%s\x1b[0m", i+2, paddedLine)
		} else {
			fmt.Fprintf(ab, "\x1b[%d;1H\x1b[K%s", i+2, line)
		}
	}

	// Clear leftover lines (offset for header)
	if s.lastLines > maxLines {
		for i := maxLines; i < s.lastLines; i++ {
			fmt.Fprintf(ab, "\x1b[%d;1H\x1b[K", i+2)
		}
	}

	// Cursor position (offset by 1 for header)
	fmt.Fprintf(ab, "\x1b[%d;%dH", cY+1, cX)
	s.lastLines = maxLines
	ab.WriteString("\x1b[?25h")
	ab.WriteTo(os.Stdout)
}

func (s *Screen) MoveCursorToStart() {
	cX, cY = 1, 1
}

func (s *Screen) MoveCursorToEnd() {
	lastLine := len(s.buffer.lines)
	lastCol := len(s.buffer.lines[len(s.buffer.lines)-1]) + 1
	cX, cY = lastCol, lastLine
}

// \x1b[%d;%dH     Move to row,col

func (s *Screen) SetPrompt(prompt string) {
	s.prompt = prompt
	fmt.Printf("\x1b[%d;1H\x1b[K%s\x1b[?25l", s.getTerminalHeight(), prompt)
}

// TODO: Later we can detect resize
func (s *Screen) getTerminalHeight() int {
	// Get terminal size - fallback to 24 if unable to determine
	if _, height, err := term.GetSize(int(os.Stdin.Fd())); err == nil {
		return height
	}
	return 24
}
