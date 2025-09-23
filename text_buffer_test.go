package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBufferNewlineCount(t *testing.T) {
	tb, err := NewTextBuffer(30)
	require.NoError(t, err)

	tb.Insert(0, 'a')
	tb.Insert(1, '\n')
	tb.Insert(2, 'b')
	tb.Insert(3, '\n')
	tb.Insert(4, 'c')

	require.Equal(t, []int{0, 2, 4}, tb.lineStarts)

	tb.Insert(1, '\n')

	require.Equal(t, "a\n\nb\nc", tb.String())
	require.Equal(t, []int{0, 2, 3, 5}, tb.lineStarts)
}

func TestBufferLineToChar(t *testing.T) {
	tb, err := NewTextBuffer(100)
	require.NoError(t, err)

	// Build complex text with mixed line endings and edge cases
	// "line1\n\n\nshort\nvery long line with many characters\n\nlast"
	text := "line1\n\n\nshort\nvery long line with many characters\n\nlast"
	for i, ch := range text {
		tb.Insert(i, ch)
	}

	// Verify the structure we built
	expected := "line1\n\n\nshort\nvery long line with many characters\n\nlast"
	require.Equal(t, expected, tb.String())

	// Test all line starts - corrected positions
	require.Equal(t, 0, tb.LineToChar(0))  // "line1\n"
	require.Equal(t, 6, tb.LineToChar(1))  // first empty line
	require.Equal(t, 7, tb.LineToChar(2))  // second empty line
	require.Equal(t, 8, tb.LineToChar(3))  // "short\n"
	require.Equal(t, 14, tb.LineToChar(4)) // "very long line..."
	require.Equal(t, 50, tb.LineToChar(5)) // empty line after long line
	require.Equal(t, 51, tb.LineToChar(6)) // "last" (no trailing newline)

	// Test dynamic insertion affecting line boundaries
	tb.Insert(5, 'X')                     // Insert before first newline: "line1X\n..."
	require.Equal(t, 0, tb.LineToChar(0)) // still starts at 0
	require.Equal(t, 7, tb.LineToChar(1)) // now starts at 7 (was 6)

	// Insert newline in middle of existing line
	tb.Insert(3, '\n')                    // "lin\ne1X\n..." - splits line1
	require.Equal(t, 0, tb.LineToChar(0)) // "lin\n"
	require.Equal(t, 4, tb.LineToChar(1)) // "e1X\n"
	require.Equal(t, 8, tb.LineToChar(2)) // first originally empty line

	// Boundary tests with the modified buffer
	require.Equal(t, 0, tb.LineToChar(-50))  // way negative
	require.Equal(t, 53, tb.LineToChar(100)) // way beyond - should be last line start

	// Test buffer with only newlines
	nlTb, _ := NewTextBuffer(10)
	nlTb.Insert(0, '\n')
	nlTb.Insert(1, '\n')
	nlTb.Insert(2, '\n') // "\n\n\n"

	require.Equal(t, 0, nlTb.LineToChar(0))  // first empty line
	require.Equal(t, 1, nlTb.LineToChar(1))  // second empty line
	require.Equal(t, 2, nlTb.LineToChar(2))  // third empty line
	require.Equal(t, 3, nlTb.LineToChar(3))  // fourth empty line (after last \n)
	require.Equal(t, 3, nlTb.LineToChar(10)) // clamp to last line
}

func TestBufferCharToLine(t *testing.T) {
	tb, err := NewTextBuffer(100)
	require.NoError(t, err)

	// Build: "abc\nde\n\nfgh\nij"
	text := "abc\nde\n\nfgh\nij"
	for i, ch := range text {
		tb.Insert(i, ch)
	}

	// Test positions within each line
	require.Equal(t, 0, tb.CharToLine(0)) // 'a'
	require.Equal(t, 0, tb.CharToLine(1)) // 'b'
	require.Equal(t, 0, tb.CharToLine(2)) // 'c'
	require.Equal(t, 0, tb.CharToLine(3)) // '\n'

	require.Equal(t, 1, tb.CharToLine(4)) // 'd'
	require.Equal(t, 1, tb.CharToLine(5)) // 'e'
	require.Equal(t, 1, tb.CharToLine(6)) // '\n'

	require.Equal(t, 2, tb.CharToLine(7)) // '\n' (empty line)

	require.Equal(t, 3, tb.CharToLine(8))  // 'f'
	require.Equal(t, 3, tb.CharToLine(9))  // 'g'
	require.Equal(t, 3, tb.CharToLine(10)) // 'h'
	require.Equal(t, 3, tb.CharToLine(11)) // '\n'

	require.Equal(t, 4, tb.CharToLine(12)) // 'i'
	require.Equal(t, 4, tb.CharToLine(13)) // 'j'

	// // Boundary clamping
	require.Equal(t, 0, tb.CharToLine(-1))   // negative -> first line
	require.Equal(t, 0, tb.CharToLine(-100)) // way negative -> first line
	require.Equal(t, 4, tb.CharToLine(14))   // at end -> last line
	require.Equal(t, 4, tb.CharToLine(100))  // way beyond -> last line
	//
	// // Edge case: position exactly at buffer length
	require.Equal(t, 4, tb.CharToLine(tb.gBuf.Length())) // at exact end
	//
	// // Dynamic test - insert newline and verify line mapping changes
	tb.Insert(2, '\n')                    // "ab\nc\nde\n\nfgh\nij"
	require.Equal(t, 0, tb.CharToLine(0)) // 'a'
	require.Equal(t, 0, tb.CharToLine(1)) // 'b'
	require.Equal(t, 0, tb.CharToLine(2)) // '\n' (new)
	require.Equal(t, 1, tb.CharToLine(3)) // 'c' (now on line 1)
	require.Equal(t, 1, tb.CharToLine(4)) // '\n'
	require.Equal(t, 2, tb.CharToLine(5)) // 'd' (now on line 2)
	//
	// // Empty buffer
	emptyTb, _ := NewTextBuffer(10)
	require.Equal(t, 0, emptyTb.CharToLine(0))
	require.Equal(t, 0, emptyTb.CharToLine(-1))
	require.Equal(t, 0, emptyTb.CharToLine(1))
}

func TestBufferInsertString(t *testing.T) {
	tb, err := NewTextBuffer(100)
	require.NoError(t, err)
	text := "abc\nde\n\nfgh\nij"
	for i, ch := range text {
		tb.Insert(i, ch)
	}
	// Verify initial state
	require.Equal(t, 14, tb.Length())
	require.Equal(t, 5, tb.LineCount())
	require.Equal(t, "abc\nde\n\nfgh\nij", tb.String())

	t.Run("insert at end (fast path)", func(t *testing.T) {
		tb2 := copyTextBuffer(tb)
		tb2.InsertString(14, "xyz\n123")
		require.Equal(t, 21, tb2.Length())
		require.Equal(t, 6, tb2.LineCount())
		require.Equal(t, "abc\nde\n\nfgh\nijxyz\n123", tb2.String())
	})

	t.Run("insert in middle (slow path)", func(t *testing.T) {
		tb2 := copyTextBuffer(tb)
		tb2.InsertString(7, "NEW\nLINE\n")
		require.Equal(t, 23, tb2.Length())
		require.Equal(t, 7, tb2.LineCount())
		require.Equal(t, "abc\nde\nNEW\nLINE\n\nfgh\nij", tb2.String())

		require.Equal(t, 4, tb2.LineToChar(1))
		require.Equal(t, 7, tb2.LineToChar(2))
		require.Equal(t, 11, tb2.LineToChar(3))
		require.Equal(t, 16, tb2.LineToChar(4))
		require.Equal(t, 17, tb2.LineToChar(5))
	})
}

func copyTextBuffer(original *TextBuffer) *TextBuffer {
	// Create new buffer with same capacity
	newTB, _ := NewTextBuffer(len(original.gBuf.buffer))

	// Copy the text content
	text := original.String()
	if len(text) > 0 {
		newTB.InsertString(0, text)
	}

	return newTB
}

func TestBufferGetLine(t *testing.T) {
	tb, err := NewTextBuffer(100)
	require.NoError(t, err)
	text := "abc\nde\n\nfgh\nij"
	for i, ch := range text {
		tb.Insert(i, ch)
	}
	require.Equal(t, 14, tb.Length())
	require.Equal(t, 5, tb.LineCount())
	require.Equal(t, "abc\nde\n\nfgh\nij", tb.String())

	require.Equal(t, "abc\n", tb.Line(0))
	require.Equal(t, "de\n", tb.Line(1))
	require.Equal(t, "\n", tb.Line(2))
	require.Equal(t, "fgh\n", tb.Line(3))
	require.Equal(t, "ij", tb.Line(4))

	require.Equal(t, "", tb.Line(-1))
	require.Equal(t, "", tb.Line(5))
	require.Equal(t, "", tb.Line(100))
}

func TestBufferLineLength(t *testing.T) {
	tb, err := NewTextBuffer(100)
	require.NoError(t, err)
	text := "abc\nde\n\nfgh\nij"
	for i, ch := range text {
		tb.Insert(i, ch)
	}
	require.Equal(t, 14, tb.Length())
	require.Equal(t, 5, tb.LineCount())
	require.Equal(t, "abc\nde\n\nfgh\nij", tb.String())

	require.Equal(t, 3, tb.LineLength(0))
	require.Equal(t, 2, tb.LineLength(1))
	require.Equal(t, 0, tb.LineLength(2))
	require.Equal(t, 3, tb.LineLength(3))
	require.Equal(t, 2, tb.LineLength(4))

	require.Equal(t, 0, tb.LineLength(-1))
	require.Equal(t, 0, tb.LineLength(5))
	require.Equal(t, 0, tb.LineLength(100))
}
