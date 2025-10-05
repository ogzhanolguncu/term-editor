package editor

import (
	"fmt"

	cursor "github.com/ogzhanolguncu/go_editor/cursor_manager"
	textbuffer "github.com/ogzhanolguncu/go_editor/text_buffer"
)

type Editor struct {
	buffer   *textbuffer.TextBuffer // Text storage and line tacking
	cursor   *cursor.CursorManager  // Track cursor position
	filename string                 // Required for tracking file name when file is loaded from a file or saved to a file
	modified bool                   // Required for tracking file modified flag on status line
	message  string                 // Required for showing confirmation messages. e.g "Are you sure you want to save" etc...
}

func New() (*Editor, error) {
	buffer, err := textbuffer.NewTextBuffer(1024)
	if err != nil {
		return nil, fmt.Errorf("editor: failed to create text buffer: %w", err)
	}

	return &Editor{
		buffer:   buffer,
		cursor:   cursor.NewCursorManager(buffer),
		filename: "",
		modified: false,
		message:  "",
	}, nil
}

func (e *Editor) InsertChar(ch rune) {
	pos := e.cursor.GetPosition()
	e.buffer.Insert(pos, ch)
	e.cursor.ApplyTextChange(pos, +1)
	e.modified = true
}

func (e *Editor) InsertString(text string) {
	pos := e.cursor.GetPosition()
	e.buffer.InsertString(pos, text)
	e.cursor.ApplyTextChange(pos, len([]rune(text)))
	e.modified = true
}

func (e *Editor) Backspace() {
	pos := e.cursor.GetPosition()
	e.buffer.Delete(pos - 1)
	e.cursor.ApplyTextChange(pos-1, -1)
	e.modified = true
}

func (e *Editor) Delete() {
	pos := e.cursor.GetPosition()
	e.buffer.Delete(pos)
	e.cursor.ApplyTextChange(pos, -1)
	e.modified = true
}

// TODO: Will change this later
func (e *Editor) DeleteRange(start, end int) {
	e.buffer.DeleteRange(start, end)
	e.modified = true
}

func (e *Editor) MoveLeft() bool {
	return e.cursor.MoveLeft()
}

func (e *Editor) MoveRight() bool {
	return e.cursor.MoveRight()
}

func (e *Editor) MoveUp() bool {
	return e.cursor.MoveUp()
}

func (e *Editor) MoveDown() bool {
	return e.cursor.MoveDown()
}

func (e *Editor) MoveToLineStart() {
	e.cursor.MoveToLineStart()
}

func (e *Editor) MoveToLineEnd() {
	e.cursor.MoveToLineEnd()
}

func (e *Editor) MoveToStart() {
	e.cursor.MoveToStart()
}

func (e *Editor) MoveToEnd() {
	e.cursor.MoveToEnd()
}

func (e *Editor) GetCursorPosition() int {
	return e.cursor.GetPosition()
}

func (e *Editor) GetLineColumn() (int, int) {
	return e.cursor.GetLineColumn()
}

func (e *Editor) GetLine(lineNum int) string {
	return e.buffer.Line(lineNum)
}

func (e *Editor) GetLineCount() int {
	return e.buffer.LineCount()
}

func (e *Editor) GetLength() int {
	return e.buffer.Length()
}

func (e *Editor) GetContent() string {
	return e.buffer.String()
}

func (e *Editor) Find(needle string) []int {
	return e.buffer.Find(needle)
}

func (e *Editor) IsModified() bool {
	return e.modified
}

func (e *Editor) GetFilename() string {
	if e.filename == "" {
		return "[No Name]"
	}
	return e.filename
}

func (e *Editor) GetMessage() string {
	msg := e.message
	e.message = "" // Clear after reading
	return msg
}

func (e *Editor) GetVisibleContent(startLine, numLines int) []string {
	lines := make([]string, 0, numLines)
	totalLines := e.buffer.LineCount()

	for i := range numLines {
		lineNum := startLine + i
		if lineNum >= totalLines {
			break
		}
		lines = append(lines, e.buffer.Line(lineNum))
	}

	return lines
}

func (e *Editor) SetMessage(msg string) {
	e.message = msg
}

func (e *Editor) GetStatusLine() string {
	modFlag := ""
	if e.modified {
		modFlag = " [+]"
	}

	line, col := e.GetLineColumn()

	return fmt.Sprintf("%s%s | Line %d, Col %d | %d lines | %d chars",
		e.GetFilename(),
		modFlag,
		line+1, // Display as 1-indexed
		col+1,  // Display as 1-indexed
		e.GetLineCount(),
		e.GetLength())
}
