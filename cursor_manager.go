package main

import (
	"fmt"
)

type Cursor struct {
	position int
}

type CursorManager struct {
	cursor *Cursor
	buffer *TextBuffer
}

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

func (cm *CursorManager) MoveToPosition(line, col int) error {
	if line < 0 || line >= cm.buffer.LineCount() {
		return fmt.Errorf("line out of bounds")
	}
	lineStart := cm.buffer.LineToChar(line)
	lineLength := cm.buffer.LineLength(line)
	if col < 0 || col > lineLength {
		return fmt.Errorf("column out of bounds")
	}
	cm.cursor.position = lineStart + col
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

func (cm *CursorManager) MoveDown() bool {
	line, col := cm.GetLineColumn()
	if line >= cm.buffer.LineCount()-1 {
		return false
	}

	targetLine := line + 1
	targetLineStart := cm.buffer.LineToChar(targetLine)
	targetLineLength := cm.buffer.LineLength(targetLine)

	targetCol := min(col, max(0, targetLineLength-1))

	cm.cursor.position = targetLineStart + targetCol
	return true
}

func (cm *CursorManager) MoveToLineStart() {
	line, _ := cm.GetLineColumn()
	lineStart := cm.buffer.LineToChar(line)
	cm.cursor.position = lineStart
}

func (cm *CursorManager) MoveToStart() {
	cm.cursor.position = 0
}

func (cm *CursorManager) MoveToLineEnd() {
	line, _ := cm.GetLineColumn()

	if line >= cm.buffer.LineCount()-1 {
		cm.cursor.position = cm.buffer.Length()
		return
	}

	nextLineStart := cm.buffer.LineToChar(line + 1)
	cm.cursor.position = nextLineStart - 1
}

func (cm *CursorManager) MoveToEnd() {
	cm.cursor.position = cm.buffer.Length()
}

func (cm *CursorManager) IsAtStart() bool {
	return cm.cursor.position == 0
}

func (cm *CursorManager) IsAtEnd() bool {
	return cm.cursor.position == cm.buffer.Length()
}
