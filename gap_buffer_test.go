package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountWords(t *testing.T) {
	gbuf, err := NewGapBuffer(30)
	require.NoError(t, err)

	require.True(t, gbuf.gapEnd == 30)
	require.True(t, gbuf.gapStart == 0)
	require.Equal(t, gbuf.Length(), 0)
	require.Equal(t, gbuf.String(), "")
}

func TestInsert(t *testing.T) {
	gbuf, err := NewGapBuffer(6)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')

	require.Equal(t, gbuf.String(), "Hel")
	require.Equal(t, gbuf.gapStart, 3)
	require.Equal(t, gbuf.gapEnd, 6)
}

func TestExpand(t *testing.T) {
	gbuf, err := NewGapBuffer(3)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')

	require.Equal(t, gbuf.String(), "Hell")
	require.Equal(t, gbuf.Length(), 4)
	require.Equal(t, gbuf.gapStart, 4)
	require.Equal(t, gbuf.gapEnd, 6)
	require.Equal(t, len(gbuf.buffer), 6)
}

func TestCursorMove(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	require.Equal(t, gbuf.gapStart, 5)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.MoveGapTo(2)
	require.Equal(t, gbuf.String(), "Hello")
	require.Equal(t, gbuf.gapStart, 2)
	require.Equal(t, gbuf.gapEnd, 7)

	gbuf.Insert('X')
	require.Equal(t, gbuf.String(), "HeXllo")

	gbuf.MoveGapTo(6)

	gbuf.Insert('T')
	require.Equal(t, gbuf.String(), "HeXlloT")
	require.Equal(t, gbuf.gapStart, 7)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.MoveGapTo(10)
	gbuf.MoveGapTo(8)

	gbuf.Insert('X')
	gbuf.Insert('Y')
	gbuf.Insert('Z')
	require.Equal(t, gbuf.gapStart, 10)
	require.Equal(t, gbuf.gapEnd, 10)
}

func TestBackspace(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	require.Equal(t, gbuf.gapStart, 5)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.Backspace()

	require.Equal(t, gbuf.gapStart, 4)
	require.Equal(t, gbuf.gapEnd, 10)

	require.Equal(t, gbuf.String(), "Hell")
}

func TestDelete(t *testing.T) {
	gbuf, err := NewGapBuffer(10)
	require.NoError(t, err)

	gbuf.Insert('H')
	gbuf.Insert('e')
	gbuf.Insert('l')
	gbuf.Insert('l')
	gbuf.Insert('o')

	require.Equal(t, gbuf.gapStart, 5)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.MoveGapTo(0)

	require.Equal(t, gbuf.gapStart, 0)
	require.Equal(t, gbuf.gapEnd, 5)

	gbuf.Delete()

	require.Equal(t, gbuf.String(), "ello")
	require.Equal(t, gbuf.gapStart, 0)
	require.Equal(t, gbuf.gapEnd, 6)

	gbuf.MoveGapTo(3)

	gbuf.Delete()

	require.Equal(t, gbuf.String(), "ell")
	require.Equal(t, gbuf.gapStart, 3)
	require.Equal(t, gbuf.gapEnd, 10)

	gbuf.Delete()
	require.Equal(t, gbuf.String(), "ell")
	require.Equal(t, gbuf.gapStart, 3)
	require.Equal(t, gbuf.gapEnd, 10)
}
