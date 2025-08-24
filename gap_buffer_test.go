package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountWords(t *testing.T) {
	gbuf := NewGapBuffer(30)

	require.True(t, gbuf.gapEnd == 30)
	require.True(t, gbuf.gapStart == 0)
	require.Equal(t, gbuf.Length(), 0)
	require.Equal(t, gbuf.ToString(), "")
}

func TestInsert(t *testing.T) {
	gbuf := NewGapBuffer(6)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')

	require.Equal(t, gbuf.ToString(), "Hel")
	require.Equal(t, gbuf.gapStart, 3)
	require.Equal(t, gbuf.gapEnd, 6)
}

func TestExpand(t *testing.T) {
	gbuf := NewGapBuffer(3)
	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')

	require.Equal(t, gbuf.ToString(), "Hell")
	require.Equal(t, gbuf.Length(), 4)
	require.Equal(t, gbuf.gapStart, 4)
	require.Equal(t, gbuf.gapEnd, 6)
	require.Equal(t, len(gbuf.buffer), 6)
}

func TestCursorMove(t *testing.T) {
	gbuf := NewGapBuffer(10)
	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	require.Equal(t, gbuf.gapStart, 5)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.MoveCursorTo(2)
	require.Equal(t, gbuf.ToString(), "Hello")
	require.Equal(t, gbuf.gapStart, 2)
	require.Equal(t, gbuf.gapEnd, 7)

	gbuf.Insert('X')
	require.Equal(t, gbuf.ToString(), "HeXllo")

	gbuf.MoveCursorTo(6)

	gbuf.Insert('T')
	require.Equal(t, gbuf.ToString(), "HeXlloT")
	require.Equal(t, gbuf.gapStart, 7)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.MoveCursorTo(10)
	gbuf.MoveCursorTo(8)

	gbuf.Insert('X')
	gbuf.Insert('Y')
	gbuf.Insert('Z')
	require.Equal(t, gbuf.gapStart, 10)
	require.Equal(t, gbuf.gapEnd, 10)
}
