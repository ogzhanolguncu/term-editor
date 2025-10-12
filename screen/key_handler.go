package screen

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/ogzhanolguncu/go_editor/editor"
)

func (s *Screen) handleKey(ev *tcell.EventKey) bool {
	mode := s.editor.GetMode()

	switch mode {
	case editor.ModeNormal:
		return s.handleNormal(ev)
	case editor.ModeInsert:
		return s.handleInsert(ev)
	}

	return true
}

func (s *Screen) handleNormal(ev *tcell.EventKey) bool {
	e := s.editor

	switch ev.Key() {
	case tcell.KeyCtrlC:
		return false
	case tcell.KeyEsc:
		return true
	case tcell.KeyLeft:
		e.MoveLeft()
	case tcell.KeyRight:
		e.MoveRight()
	case tcell.KeyUp:
		e.MoveUp()
	case tcell.KeyDown:
		e.MoveDown()
	}

	r := ev.Rune()

	// This one handles the command count
	if e.HandleDigit(r) {
		return true
	}

	switch ev.Rune() {
	case 'i':
		e.SetMode(editor.ModeInsert)
	case 'a':
		e.MoveRight()
		e.SetMode(editor.ModeInsert)
	case 'A':
		e.MoveToLineEnd()
		e.SetMode(editor.ModeInsert)
	case 'I':
		e.MoveToLineStart()
		e.SetMode(editor.ModeInsert)
	case 'x':
		e.Delete()
	case '0':
		e.MoveToLineStart()
	case '$':
		e.MoveToLineEnd()
	case 'G':
		e.MoveToEnd()
	case 'h':
		e.MoveLeft()
	case 'j':
		e.MoveDown()
	case 'k':
		e.MoveUp()
	case 'l':
		e.MoveRight()
	case 'w':
		e.MoveToNextWord()
	case 'b':
		e.MoveToPrevWord()
	}
	return true
}

func (s *Screen) handleInsert(ev *tcell.EventKey) bool {
	e := s.editor
	switch ev.Key() {
	case tcell.KeyEsc, tcell.KeyCtrlC:
		e.SetMode(editor.ModeNormal)
		return true
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		e.Backspace()
	case tcell.KeyDelete:
		e.Delete()
	case tcell.KeyEnter:
		e.InsertChar('\n')
	case tcell.KeyLeft:
		e.MoveLeft()
	case tcell.KeyRight:
		e.MoveRight()
	case tcell.KeyUp:
		e.MoveUp()
	case tcell.KeyDown:
		e.MoveDown()
	case tcell.KeyHome:
		e.MoveToLineStart()
	case tcell.KeyEnd:
		e.MoveToLineEnd()
	case tcell.KeyTAB:
		e.InsertString(strings.Repeat(" ", tabSize))
	case tcell.KeyRune:
		e.InsertChar(ev.Rune())
	}

	return true
}
