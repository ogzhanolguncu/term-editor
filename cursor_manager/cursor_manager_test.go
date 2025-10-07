package cursormanager

import (
	"testing"

	textbuffer "github.com/ogzhanolguncu/go_editor/text_buffer"
	"github.com/stretchr/testify/require"
)

func TestNewCursorManager(t *testing.T) {
	tb, err := textbuffer.NewTextBuffer(100)
	require.NoError(t, err)

	cm := NewCursorManager(tb)
	require.NotNil(t, cm)
	require.NotNil(t, cm.cursor)
	require.Equal(t, tb, cm.buffer)
	require.Equal(t, 0, cm.GetPosition())
}

func TestCursorGetPosition(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	cm := NewCursorManager(tb)

	require.Equal(t, 0, cm.GetPosition())

	// Manually set cursor position to test getter
	cm.cursor.position = 42
	require.Equal(t, 42, cm.GetPosition())

	cm.cursor.position = 100
	require.Equal(t, 100, cm.GetPosition())
}

func TestCursorSetPosition(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello World")
	cm := NewCursorManager(tb)

	// Valid positions
	err := cm.SetPosition(0)
	require.NoError(t, err)
	require.Equal(t, 0, cm.GetPosition())

	err = cm.SetPosition(5)
	require.NoError(t, err)
	require.Equal(t, 5, cm.GetPosition())

	err = cm.SetPosition(11) // At end of buffer
	require.NoError(t, err)
	require.Equal(t, 11, cm.GetPosition())

	// Invalid positions
	err = cm.SetPosition(-1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "position out of bounds")
	require.Equal(t, 11, cm.GetPosition()) // Should not change

	err = cm.SetPosition(12) // Beyond buffer
	require.Error(t, err)
	require.Contains(t, err.Error(), "position out of bounds")
	require.Equal(t, 11, cm.GetPosition()) // Should not change

	err = cm.SetPosition(100) // Way beyond
	require.Error(t, err)
	require.Contains(t, err.Error(), "position out of bounds")
	require.Equal(t, 11, cm.GetPosition()) // Should not change
}

func TestCursorSetPositionEmptyBuffer(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	cm := NewCursorManager(tb)

	// Empty buffer - only position 0 is valid
	err := cm.SetPosition(0)
	require.NoError(t, err)
	require.Equal(t, 0, cm.GetPosition())

	err = cm.SetPosition(1)
	require.Error(t, err)
	require.Contains(t, err.Error(), "position out of bounds")
}

func TestApplyTextChangeInsert(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello World")
	cm := NewCursorManager(tb)

	t.Run("insert before cursor", func(t *testing.T) {
		cm.SetPosition(6)                     // Position after "Hello "
		cm.ApplyTextChange(2, 3)              // Insert "XXX" at position 2
		require.Equal(t, 9, cm.GetPosition()) // Cursor shifted right by 3
	})

	t.Run("insert at cursor position", func(t *testing.T) {
		cm.SetPosition(5)
		cm.ApplyTextChange(5, 2)              // Insert "YY" at cursor position
		require.Equal(t, 7, cm.GetPosition()) // Cursor shifted right by 2
	})

	t.Run("insert after cursor", func(t *testing.T) {
		cm.SetPosition(3)
		cm.ApplyTextChange(8, 4)              // Insert "ZZZZ" after cursor
		require.Equal(t, 3, cm.GetPosition()) // Cursor unchanged
	})

	t.Run("insert at beginning", func(t *testing.T) {
		cm.SetPosition(5)
		cm.ApplyTextChange(0, 1)              // Insert "A" at beginning
		require.Equal(t, 6, cm.GetPosition()) // Cursor shifted right by 1
	})

	t.Run("cursor at position 0", func(t *testing.T) {
		cm.SetPosition(0)
		cm.ApplyTextChange(0, 2)              // Insert "BB" at position 0
		require.Equal(t, 2, cm.GetPosition()) // Cursor shifted right by 2
	})
}

func TestApplyTextChangeDelete(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello World")
	cm := NewCursorManager(tb)

	t.Run("delete before cursor", func(t *testing.T) {
		cm.SetPosition(8)                     // Position in "World"
		cm.ApplyTextChange(2, -3)             // Delete 3 chars starting at position 2
		require.Equal(t, 5, cm.GetPosition()) // Cursor shifted left by 3
	})

	t.Run("delete after cursor", func(t *testing.T) {
		cm.SetPosition(3)
		cm.ApplyTextChange(7, -2)             // Delete 2 chars after cursor
		require.Equal(t, 3, cm.GetPosition()) // Cursor unchanged
	})

	t.Run("delete cursor inside deleted range", func(t *testing.T) {
		cm.SetPosition(5)
		cm.ApplyTextChange(3, -4)             // Delete 4 chars starting at position 3 (covers cursor)
		require.Equal(t, 3, cm.GetPosition()) // Cursor clamped to deletion start
	})

	t.Run("cursor at deletion start", func(t *testing.T) {
		cm.SetPosition(4)
		cm.ApplyTextChange(4, -2)             // Delete 2 chars starting at cursor
		require.Equal(t, 4, cm.GetPosition()) // Cursor stays at deletion point
	})

	t.Run("delete at beginning", func(t *testing.T) {
		cm.SetPosition(5)
		cm.ApplyTextChange(0, -2)             // Delete 2 chars from beginning
		require.Equal(t, 3, cm.GetPosition()) // Cursor shifted left by 2
	})
}

func TestApplyTextChangeReplace(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello World")
	cm := NewCursorManager(tb)

	t.Run("replace same length", func(t *testing.T) {
		cm.SetPosition(8)
		cm.ApplyTextChange(6, 0)              // Replace with same length at position 6
		require.Equal(t, 8, cm.GetPosition()) // Cursor unchanged
	})

	t.Run("replace with longer text", func(t *testing.T) {
		cm.SetPosition(8)
		cm.ApplyTextChange(2, 3)               // Replace with 3 additional chars
		require.Equal(t, 11, cm.GetPosition()) // Cursor shifted right by 3
	})

	t.Run("replace with shorter text", func(t *testing.T) {
		cm.SetPosition(8)
		cm.ApplyTextChange(2, -2)             // Replace with 2 fewer chars
		require.Equal(t, 6, cm.GetPosition()) // Cursor shifted left by 2
	})
}

func TestApplyTextChangeEdgeCases(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	cm := NewCursorManager(tb)

	t.Run("empty buffer", func(t *testing.T) {
		cm.SetPosition(0)
		cm.ApplyTextChange(0, 5)              // Insert 5 chars in empty buffer
		require.Equal(t, 5, cm.GetPosition()) // Cursor shifted right
	})

	t.Run("zero delta", func(t *testing.T) {
		tb.InsertString(0, "test")
		cm.SetPosition(2)
		cm.ApplyTextChange(1, 0)              // No net change
		require.Equal(t, 2, cm.GetPosition()) // Cursor unchanged
	})

	t.Run("cursor at end of buffer", func(t *testing.T) {
		tb2, _ := textbuffer.NewTextBuffer(100)
		tb2.InsertString(0, "Hello")
		cm2 := NewCursorManager(tb2)
		cm2.SetPosition(5)                     // At end
		cm2.ApplyTextChange(2, 3)              // Insert before end
		require.Equal(t, 8, cm2.GetPosition()) // Cursor shifted
	})
}

func TestGetLineColumn(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "line1\nline2\n\nline4")
	cm := NewCursorManager(tb)

	// Position 0: start of "line1"
	cm.SetPosition(0)
	line, col := cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 0, col)

	// Position 3: middle of "line1"
	cm.SetPosition(3)
	line, col = cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 3, col)

	// Position 5: newline at end of "line1"
	cm.SetPosition(5)
	line, col = cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 5, col)

	// Position 6: start of "line2"
	cm.SetPosition(6)
	line, col = cm.GetLineColumn()
	require.Equal(t, 1, line)
	require.Equal(t, 0, col)

	// Position 9: middle of "line2"
	cm.SetPosition(9)
	line, col = cm.GetLineColumn()
	require.Equal(t, 1, line)
	require.Equal(t, 3, col)

	// Position 12: empty line
	cm.SetPosition(12)
	line, col = cm.GetLineColumn()
	require.Equal(t, 2, line)
	require.Equal(t, 0, col)

	// Position 13: start of "line4"
	cm.SetPosition(13)
	line, col = cm.GetLineColumn()
	require.Equal(t, 3, line)
	require.Equal(t, 0, col)

	// Position 18: end of buffer
	cm.SetPosition(18)
	line, col = cm.GetLineColumn()
	require.Equal(t, 3, line)
	require.Equal(t, 5, col)
}

func TestGetLineColumnEmptyBuffer(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	cm := NewCursorManager(tb)

	line, col := cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 0, col)
}

func TestGetLineColumnSingleLine(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello World")
	cm := NewCursorManager(tb)

	cm.SetPosition(0)
	line, col := cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 0, col)

	cm.SetPosition(6)
	line, col = cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 6, col)

	cm.SetPosition(11)
	line, col = cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 11, col)
}

func TestMoveLeft(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello\n")
	cm := NewCursorManager(tb)

	// Move to end first
	cm.SetPosition(6)

	// Move left 6 times should succeed
	require.True(t, cm.MoveLeft())
	require.True(t, cm.MoveLeft())
	require.True(t, cm.MoveLeft())
	require.True(t, cm.MoveLeft())
	require.True(t, cm.MoveLeft())
	require.True(t, cm.MoveLeft())

	// 7th move left should fail - already at position 0
	require.False(t, cm.MoveLeft())
	require.Equal(t, 0, cm.GetPosition())
}

func TestMoveUp(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello\nWorld\nTest")
	cm := NewCursorManager(tb)

	// Start at line 2, column 1 (position 13: "Hello\nWorld\nTe|st")
	cm.SetPosition(12)
	line, col := cm.GetLineColumn()
	require.Equal(t, 2, line)
	require.Equal(t, 0, col)

	require.True(t, cm.MoveUp())

	line, col = cm.GetLineColumn()
	require.Equal(t, 1, line)
	require.Equal(t, 0, col)

	require.True(t, cm.MoveUp())
	require.Equal(t, 0, cm.GetPosition())
	line, col = cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 0, col)
}

func TestMoveUpColumnClamping(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "abcd\nef\nghhr\nkeke")
	cm := NewCursorManager(tb)

	cm.SetPosition(16)
	line, col := cm.GetLineColumn()
	require.Equal(t, 3, line)
	require.Equal(t, 3, col)

	require.True(t, cm.MoveUp())
	require.Equal(t, 11, cm.GetPosition()) // Position 8 + 3
	line, col = cm.GetLineColumn()
	require.Equal(t, 2, line)
	require.Equal(t, 3, col)

	require.True(t, cm.MoveUp())
	require.Equal(t, 6, cm.GetPosition()) // Position 5 + 1
	line, col = cm.GetLineColumn()
	require.Equal(t, 1, line)
	require.Equal(t, 1, col)

	require.True(t, cm.MoveUp())
	require.Equal(t, 1, cm.GetPosition()) // Position 0 + 1
	line, col = cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 1, col)

	require.False(t, cm.MoveUp())
}

func TestMoveDown(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello\nWorld\nTest")
	cm := NewCursorManager(tb)

	// Start at line 0, column 0
	cm.SetPosition(0)
	line, col := cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 0, col)

	// Move down to line 1, column 0
	require.True(t, cm.MoveDown())
	require.Equal(t, 6, cm.GetPosition())
	line, col = cm.GetLineColumn()
	require.Equal(t, 1, line)
	require.Equal(t, 0, col)

	// Move down to line 2, column 0
	require.True(t, cm.MoveDown())
	require.Equal(t, 12, cm.GetPosition())
	line, col = cm.GetLineColumn()
	require.Equal(t, 2, line)
	require.Equal(t, 0, col)

	// Can't move down from last line
	require.False(t, cm.MoveDown())
	require.Equal(t, 12, cm.GetPosition())
}

func TestMoveDownColumnClamping(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "abcd\nef\nghhr\nkeke")
	cm := NewCursorManager(tb)

	// Start at line 0, column 3 (at 'd')
	cm.SetPosition(3)
	line, col := cm.GetLineColumn()
	require.Equal(t, 0, line)
	require.Equal(t, 3, col)

	// Move down to line 1 "ef" (only 2 chars) - column 3 doesn't exist, clamp to column 1
	require.True(t, cm.MoveDown())
	require.Equal(t, 6, cm.GetPosition()) // Position 5 + 1
	line, col = cm.GetLineColumn()
	require.Equal(t, 1, line)
	require.Equal(t, 1, col)

	// Move down to line 2 "ghhr" (4 chars) - column 1 exists
	require.True(t, cm.MoveDown())
	require.Equal(t, 9, cm.GetPosition()) // Position 8 + 1
	line, col = cm.GetLineColumn()
	require.Equal(t, 2, line)
	require.Equal(t, 1, col)

	// Move down to line 3 "keke" (4 chars) - column 1 exists
	require.True(t, cm.MoveDown())
	require.Equal(t, 14, cm.GetPosition()) // Position 13 + 1
	line, col = cm.GetLineColumn()
	require.Equal(t, 3, line)
	require.Equal(t, 1, col)

	// Can't move down from last line
	require.False(t, cm.MoveDown())
}

func TestMoveToLineStart(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello\nWorld\nTest")
	cm := NewCursorManager(tb)

	cm.SetPosition(3)
	cm.MoveToLineStart()

	// Test from end of first line (on newline)
	cm.SetPosition(5) // "Hello|\n"
	cm.MoveToLineStart()
	require.Equal(t, 0, cm.GetPosition())

	// Test from start of second line (already at start)
	cm.SetPosition(6) // "|World"
	cm.MoveToLineStart()

	// Test from middle of second line
	cm.SetPosition(9) // "Wor|ld"
	cm.MoveToLineStart()
	require.Equal(t, 6, cm.GetPosition())

	// Test from end of buffer
	cm.SetPosition(16) // "Test|"
	cm.MoveToLineStart()
	require.Equal(t, 12, cm.GetPosition())
}

func TestMoveToLineEnd(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello\nWorld\nTest")
	cm := NewCursorManager(tb)

	// Test from start of first line
	cm.SetPosition(0) // "|Hello"
	cm.MoveToLineEnd()
	require.Equal(t, 5, cm.GetPosition())

	// Test from middle of first line
	cm.SetPosition(3) // "Hel|lo"
	cm.MoveToLineEnd()
	require.Equal(t, 5, cm.GetPosition())

	// Test from end of first line (already at end, on newline)
	cm.SetPosition(5) // "Hello|\n"
	cm.MoveToLineEnd()
	require.Equal(t, 5, cm.GetPosition())

	// Test from start of second line
	cm.SetPosition(6) // "|World"
	cm.MoveToLineEnd()
	require.Equal(t, 11, cm.GetPosition())

	// Test from middle of second line
	cm.SetPosition(9) // "Wor|ld"
	cm.MoveToLineEnd()
	require.Equal(t, 11, cm.GetPosition())
	//
	// Test from end of buffer
	cm.SetPosition(12) // "|Test"
	cm.MoveToLineEnd()
	require.Equal(t, 16, cm.GetPosition())
	//
	// // Already at end of buffer
	cm.SetPosition(16)
	cm.MoveToLineEnd()
	require.Equal(t, 16, cm.GetPosition())
}

func TestIsAtStart(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello\nWorld")
	cm := NewCursorManager(tb)

	cm.SetPosition(0)
	require.True(t, cm.IsAtStart())

	cm.SetPosition(1)
	require.False(t, cm.IsAtStart())
}

func TestIsAtEnd(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello\nWorld")
	cm := NewCursorManager(tb)

	cm.SetPosition(0)
	require.False(t, cm.IsAtEnd())

	cm.SetPosition(11)
	require.True(t, cm.IsAtEnd())
}

func TestIsAtEndEmptyBuffer(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	cm := NewCursorManager(tb)

	require.True(t, cm.IsAtStart())
	require.True(t, cm.IsAtEnd())
}

func TestMoveToNextWord(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello World XXX YYY")
	cm := NewCursorManager(tb)

	cm.SetPosition(0)
	cm.MoveToNextWord()
	require.Equal(t, 6, cm.cursor.position)
	cm.MoveToNextWord()
	require.Equal(t, 12, cm.cursor.position)
	cm.MoveToNextWord()
	require.Equal(t, 16, cm.cursor.position)
	cm.MoveToNextWord()
	require.Equal(t, 19, cm.cursor.position)
}

func TestMoveToNextWordEdge(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello")
	cm := NewCursorManager(tb)

	cm.SetPosition(0)
	cm.MoveToNextWord()
	require.Equal(t, 5, cm.cursor.position)
}

func TestMoveToPrevWord(t *testing.T) {
	tb, _ := textbuffer.NewTextBuffer(100)
	tb.InsertString(0, "Hello World XXX YYY")
	cm := NewCursorManager(tb)

	cm.SetPosition(18)
	cm.MoveToPrevWord()
	require.Equal(t, 16, cm.cursor.position)
	cm.MoveToPrevWord()
	require.Equal(t, 12, cm.cursor.position)
	cm.MoveToPrevWord()
	require.Equal(t, 6, cm.cursor.position)
	cm.MoveToPrevWord()
	require.Equal(t, 0, cm.cursor.position)
}
