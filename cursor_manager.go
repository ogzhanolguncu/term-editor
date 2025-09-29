// This handles cursor positioning.
package main

import (
	"fmt"
)

// CURSOR MOVEMENT OPERATIONS:
// [ ] - MoveDown() bool - move to same column on next line
// [ ] - MoveToLineStart() bool - move to beginning of current line
// [ ] - MoveToLineEnd() bool - move to end of current line

// TEXT-AWARE NAVIGATION:
// [ ] - MoveToLine(lineNum int) error - jump to start of specific line
// [ ] - IsAtStart() bool - check if cursor at position 0
// [ ] - IsAtEnd() bool - check if cursor at end of buffer

type Cursor struct {
	position int
}

type CursorManager struct {
	cursor *Cursor
	buffer *TextBuffer
}

// BASIC MOVEMENT:
// [ ] IsAtStart() bool
// [ ] IsAtEnd() bool

// LINE MOVEMENT:
// [ ] MoveDown() bool
// [ ] MoveToLineStart() bool
// [ ] MoveToLineEnd() bool
// [ ] MoveToLine(lineNum int) error

// PHASE 4 - INTERFACE COMPLIANCE (for future multi-cursor):
// [ ] GetCursors() []Cursor
// [ ] MoveCursor(index int, pos int) error

func NewCursorManager(buffer *TextBuffer) *CursorManager {
	return &CursorManager{
		cursor: &Cursor{position: 0},
		buffer: buffer,
	}
}

func (cm *CursorManager) GetPosition() int {
	return cm.cursor.position
}

func (cm *CursorManager) SetPosition(pos int) error {
	if pos < 0 || pos > cm.buffer.Length() {
		return fmt.Errorf("position out of bounds")
	}
	cm.cursor.position = pos
	return nil
}

func (cm *CursorManager) ApplyTextChange(changePos int, delta int) {
	// Text inserted
	if delta > 0 {
		// Applied change has to be smaller or right at the end of cursor
		if cm.cursor.position >= changePos {
			cm.cursor.position += delta
		}
		// Text deleted
	} else if delta < 0 {
		// e.g. If cursor is at pos=20 and user is trying to delete 35 to 38 that doesnn't change the cursor position so we don't really care
		if cm.cursor.position > changePos {
			// Cursor has to go back to changePos because it was within the deleted range
			if changePos+(-delta) >= cm.cursor.position {
				cm.cursor.position = changePos
			} else {
				// Cursor is after deleted range so we just change deduct delta from cursor
				cm.cursor.position -= (-delta)
			}
		}
	}
}

// GetLineColumn returns line, column
func (cm *CursorManager) GetLineColumn() (int, int) {
	line := cm.buffer.CharToLine(cm.cursor.position)
	lineStart := cm.buffer.LineToChar(line)
	// Column is the offset from the start of the line
	return line, cm.cursor.position - lineStart
}

func (cm *CursorManager) MoveRight() bool {
	pos := cm.cursor.position
	if pos >= cm.buffer.Length() {
		return false
	}
	cm.cursor.position++
	return true
}

func (cm *CursorManager) MoveLeft() bool {
	if cm.cursor.position == 0 {
		return false
	}
	cm.cursor.position--
	return true
}

func (cm *CursorManager) MoveUp() bool {
	line, col := cm.GetLineColumn()
	if line == 0 {
		return false
	}

	targetLine := line - 1
	targetLineStart := cm.buffer.LineToChar(targetLine)
	targetLineLength := cm.buffer.LineLength(targetLine)

	targetCol := min(col, max(0, targetLineLength-1))

	cm.cursor.position = targetLineStart + targetCol
	return true
}
