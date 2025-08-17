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
// [x] Implement Ctrl+S save functionality with prompting at the bottom
// [x] Add command line argument support for opening file
// [x] Add file loading (Ctrl+O) - load file content, clear buffer, reset cursor
// [x] Track modified state - bool flag, set on edits, clear on save/load
// [x] Implement Ctrl+N for new file
// [x] Add Ctrl+Q quit with save prompt
// CRITICAL REFACTORING (Before Scrolling):
// [x] Move cursor state into Screen struct - remove global cX,cY, add cursorX,cursorY fields, create MoveCursor()/GetCursor() methods
// [x] Fix cursor management - replace scattered cX,cY assignments with Screen methods, buffer ops shouldn't touch cursor
// [x] Add bottom status line - "filename [Modified] | Line X, Col Y | N lines", real-time updates, truncate long names
// MAJOR FEATURES:
// [ ] Vertical scrolling - add scrollY offset, render visible lines only, auto-scroll on cursor move, Page Up/Down support
// [ ] Search functionality (Ctrl+F) - real-time highlighting, F3/Shift+F3 navigation, integration with scrolling
// MINOR:
// [ ] Handle Home/End keys properly
// [ ] Clipboard operations (Ctrl+A/X/C/V)
// [ ] Basic syntax highlighting
// [ ] Undo/Redo (Ctrl+Z/Y)
// [ ] Line numbers, word wrap, terminal resize handling

const Version = "GICO 0.1"

type Key int

const (
	// ASCII/Control codes with explicit values
	CtrlC     Key = 3
	CtrlN     Key = 14
	CtrlO     Key = 15
	CtrlQ     Key = 17
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
	buffer.OpenFile()
	if buffer.fPath != "" {
		screen.MoveCursorToEnd()
	}
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
	case CtrlQ:
		handleQuitWithSave()
	case CtrlS:
		handleSave()
	case CtrlO:
		handleLoadFile()
	case CtrlN:
		handleNewFile()
	case CtrlC:
		handleQuit()
	case Backspace:
		handleBackspace()
	case Enter, ArrowLeft, ArrowRight, ArrowUp, ArrowDown:
		handleCursorMove(key)
	default:
		handleCharInsert(key)
	}
}

func handleNewFile() {
	if screen.buffer.modified {
		screen.SetPrompt("Save changes? (y/n/c): ")
		for {
			key := readScreenInput()
			switch key {
			case 'y', 'Y':
				handleSave()
				screen.Restart()
			case 'n', 'N':
				screen.Restart()
			case 'c', 'C', CtrlC:
				screen.SetPrompt("")
				return
			default:
				continue
			}
		}
	}
	screen.Restart()
}

func handleQuit() {
	screen.Clear()
	os.Exit(0)
}

func handleCursorMove(key Key) {
	switch key {
	case ArrowLeft:
		if screen.cX > 1 {
			screen.cX--
		} else if screen.cY > 1 {
			// Wrap to end of previous line
			screen.cY--
			screen.cX = len(buffer.lines[screen.cY-1]) + 1
		}
	case ArrowRight:
		if screen.cY-1 < len(buffer.lines) {
			line := buffer.lines[screen.cY-1]
			if screen.cX-1 < len(line) {
				screen.cX++
			} else if screen.cY < len(buffer.lines) {
				// Wrap to start of next line
				screen.cY++
				screen.cX = 1
			}
		}
	case ArrowUp:
		if screen.cY > 1 {
			screen.cY--
			line := buffer.lines[screen.cY-1]
			// Clamp cursor to avoid going past line end
			if screen.cX-1 > len(line) {
				screen.cX = len(line) + 1
			}
		}
	case ArrowDown:
		if screen.cY < len(buffer.lines) {
			screen.cY++
			if screen.cY-1 < len(buffer.lines) {
				line := buffer.lines[screen.cY-1]
				// Clamp cursor to avoid going past line end
				if screen.cX-1 > len(line) {
					screen.cX = len(line) + 1
				}
			}
		}
	case Enter:
		line := buffer.lines[screen.cY-1]
		leftPart := line[:screen.cX-1]  // Before cursor
		rightPart := line[screen.cX-1:] // After cursor

		buffer.lines[screen.cY-1] = leftPart

		// Insert new line with right part
		newLine := []string{rightPart}
		buffer.lines = append(buffer.lines[:screen.cY], append(newLine, buffer.lines[screen.cY:]...)...)

		screen.cY++
		screen.cX = 1
		buffer.modified = true
	}
}

func handleBackspace() {
	// Cursor is at beginning and if there is a line above
	if screen.cX == 1 && screen.cY > 1 {
		line := buffer.lines[screen.cY-1]
		buffer.lines[screen.cY-2] = fmt.Sprintf("%s%s", buffer.lines[screen.cY-2], line)
		// Remove current line from the slice
		buffer.lines = append(buffer.lines[:screen.cY-1], buffer.lines[screen.cY:]...)
		// Move cursor to above line and end of it
		screen.cY--
		screen.cX = len(buffer.lines[screen.cY-1]) + 1
		buffer.modified = true
	} else if screen.cX > 1 {
		line := buffer.lines[screen.cY-1]
		buffer.lines[screen.cY-1] = fmt.Sprintf("%s%s", line[:screen.cX-2], line[screen.cX-1:])
		screen.cX--
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
	screen.MoveCursorToStart()
}

func handleQuitWithSave() {
	screen.SetPrompt("Save changes? (y/n/c): ")

	for {
		key := readScreenInput()
		switch key {
		case 'y', 'Y':
			handleSave()
			handleQuit()
		case 'n', 'N':
			handleQuit()
		case 'c', 'C', CtrlC:
			screen.SetPrompt("")
			return
		default:
			continue
		}
	}
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

	return tb
}

func (tb *TextBuffer) OpenFile() {
	if tb.fPath == "" {
		return
	}

	content, err := os.ReadFile(tb.fPath)
	if err != nil {
		if os.IsNotExist(err) {
			screen.MoveCursorToStart()
			return
		}
		nuke(err)
	}
	tb.lines = strings.Split(strings.TrimSpace(string(content)), "\n")
}

func handleCharInsert(ch Key) {
	line := buffer.lines[screen.cY-1]
	// This is required for inline editing like editing an item from middle.
	line = line[:screen.cX-1] + string(rune(ch)) + line[screen.cX-1:]
	buffer.lines[screen.cY-1] = line
	buffer.modified = true
	screen.cX++
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
	cX, cY    int
	buffer    *TextBuffer
	lastLines int
	prompt    string
}

func NewScreen(textBuffer *TextBuffer) *Screen {
	screen := &Screen{
		buffer: textBuffer,
		cX:     1,
		cY:     1,
	}
	screen.Clear()
	return screen
}

func (s *Screen) Clear() {
	fmt.Print("\x1b[2J") // Clear entire screen
	fmt.Print("\x1b[H")  // Move cursor to top-left (1,1)
}

// TODO: keep it remove it idk
const (
	GruvboxBg0    = "235" // #282828 - dark background
	GruvboxBg1    = "237" // #3c3836 - lighter background
	GruvboxFg1    = "223" // #ebdbb2 - light foreground
	GruvboxOrange = "208" // #fe8019 - orange accent
	GruvboxYellow = "214" // #fabd2f - yellow
	GruvboxGreen  = "142" // #b8bb26 - green
	GruvboxBlue   = "109" // #83a598 - blue
)

func (s *Screen) renderStatusLine(ab *bytes.Buffer) {
	width, height := s.getTerminalSize()
	filename := "New Buffer"
	if s.buffer.fPath != "" {
		filename = getDisplayPath(s.buffer.fPath)
		maxFilenameWidth := width * 35 / 100 // 35% of terminal width

		if len(filename) > maxFilenameWidth {
			filename = "..." + filename[len(filename)-maxFilenameWidth+3:]
		}
	}

	modifiedFlag := ""
	if s.buffer.modified {
		modifiedFlag = " [Modified]"
	}

	time := time.Now().Format("15:04")
	leftStatus := fmt.Sprintf(
		"%s%s  |  Line %d, Col %d  |  %d Lines  |  (T) %s", filename,
		modifiedFlag,
		s.cY,
		s.cX,
		len(s.buffer.lines),
		time,
	)
	remainingWidth := width - len(leftStatus)

	// Render with right-aligned time
	fmt.Fprintf(ab, "\x1b[%d;1H\x1b[38;5;%s;48;5;%sm%s%*s\x1b[0m",
		height, GruvboxFg1, GruvboxBg1, leftStatus, remainingWidth, Version)
}

func (s *Screen) Refresh() {
	ab := bytes.NewBufferString("\x1b[?25l") // Hide cursor

	width, _ := s.getTerminalSize()

	// Render buffer content starting from line 1
	maxLines := max(s.cY, len(s.buffer.lines))
	for i := range maxLines {
		line := ""
		if i < len(s.buffer.lines) {
			line = s.buffer.lines[i]
		}

		// Current line highlighting
		if i+1 == s.cY {
			// Pad line to full terminal width for complete background highlight
			paddedLine := fmt.Sprintf("%-*s", width, line)
			fmt.Fprintf(ab, "\x1b[%d;1H\x1b[K\x1b[48;5;236m%s\x1b[0m", i+1, paddedLine)
		} else {
			fmt.Fprintf(ab, "\x1b[%d;1H\x1b[K%s", i+1, line)
		}
	}

	// Clear leftover lines (offset for header)
	if s.lastLines > maxLines {
		for i := maxLines; i < s.lastLines; i++ {
			fmt.Fprintf(ab, "\x1b[%d;1H\x1b[K", i+1)
		}
	}

	s.renderStatusLine(ab)

	// Cursor position (offset by 1 for header)
	fmt.Fprintf(ab, "\x1b[%d;%dH", s.cY, s.cX)
	s.lastLines = maxLines
	ab.WriteString("\x1b[?25h")
	ab.WriteTo(os.Stdout)
}

func (s *Screen) MoveCursorToStart() {
	s.cX, s.cY = 1, 1
}

func (s *Screen) MoveCursorToEnd() {
	lastLine := len(s.buffer.lines)
	lastCol := len(s.buffer.lines[len(s.buffer.lines)-1]) + 1
	s.cX, s.cY = lastCol, lastLine
}

func (s *Screen) ClearAndQuit() {
	screen.Clear()
	os.Exit(0)
}

func (s *Screen) Restart() {
	s.Clear()
	s.MoveCursorToStart()
	s.buffer.fPath = ""
	s.buffer.modified = false
	s.buffer.lines = []string{""}
}

func (s *Screen) SetPrompt(prompt string) {
	s.prompt = prompt
	_, height := s.getTerminalSize()
	fmt.Printf("\x1b[%d;1H\x1b[K%s\x1b[?25l", height, prompt)
}

func (s *Screen) getTerminalSize() (width, height int) {
	if width, height, err := term.GetSize(int(os.Stdin.Fd())); err == nil {
		return width, height
	}
	return 80, 24
}
