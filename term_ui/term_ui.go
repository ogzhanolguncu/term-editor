package termui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/ogzhanolguncu/go_editor/editor"
)

type TerminalUI struct {
	screen  tcell.Screen
	editor  *editor.Editor
	palette *Palette

	width        int
	height       int
	scrollOffset int // Top line currently visible
}

func NewTerminalUI(editor *editor.Editor) (*TerminalUI, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	width, height := screen.Size()

	return &TerminalUI{
		screen:       screen,
		editor:       editor,
		palette:      NewPalette(),
		width:        width,
		height:       height,
		scrollOffset: 0,
	}, nil
}

func (ui *TerminalUI) Run() {
	for {
		ui.render()

		ev := ui.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			if !ui.handleKey(ev) {
				return // Quit signal
			}

		case *tcell.EventResize:
			ui.screen.Sync()
		}
	}
}

func (ui *TerminalUI) Close() {
	ui.screen.Fini()
}

func (ui *TerminalUI) render() {
	ui.screen.Clear()
	ui.width, ui.height = ui.screen.Size()

	bgStyle := ui.palette.StyleForNormalText()
	for y := 0; y < ui.height; y++ {
		for x := 0; x < ui.width; x++ {
			ui.screen.SetContent(x, y, ' ', nil, bgStyle)
		}
	}

	// Content height has to account for status bar thats why we skip the last line
	contentHeight := ui.height - 1
	cursorLine, cursorCol := ui.editor.GetLineColumn()

	if cursorLine < ui.scrollOffset {
		ui.scrollOffset = cursorLine
	}
	if cursorLine >= ui.scrollOffset+contentHeight {
		ui.scrollOffset = cursorLine - contentHeight + 1
	}

	// Calculate gutter width based on max line number
	gutterWidth := len(fmt.Sprintf("%d", ui.scrollOffset+contentHeight)) + 1 // +1 for space after number

	lines := ui.editor.GetVisibleContent(ui.scrollOffset, contentHeight)

	ui.renderLines(lines, gutterWidth, cursorLine)
	ui.renderStatusBar()

	// Position cursor accounting for gutter
	screenRow := cursorLine - ui.scrollOffset
	screenCol := cursorCol + gutterWidth + 2 // Offset by gutter width + separator + space
	ui.screen.ShowCursor(screenCol, screenRow)
	ui.screen.Show()
}

func (ui *TerminalUI) renderLines(lines []string, gutterWidth int, cursorLine int) {
	for row, lineContent := range lines {
		lineNum := ui.scrollOffset + row + 1

		// Draw line number
		gutterText := fmt.Sprintf("%*d ", gutterWidth-1, lineNum)
		for col, ch := range gutterText {
			ui.screen.SetContent(col, row, ch, nil, ui.palette.StyleForLineNum())
		}

		// Draw vertical separator
		ui.screen.SetContent(gutterWidth, row, '│', nil, ui.palette.StyleForGutter())

		// Draw text content
		textStartCol := gutterWidth + 2
		style := ui.palette.StyleForNormalText()
		if ui.scrollOffset+row == cursorLine {
			style = ui.palette.StyleForCurrentLine()
		}

		ui.drawLine(textStartCol, row, lineContent, style)
	}
}

func (ui *TerminalUI) renderStatusBar() {
	statusLine := ui.editor.GetStatusLine()
	ui.drawLine(0, ui.height-1, statusLine, ui.palette.StyleForStatusBar())

	// Show message if any
	if msg := ui.editor.GetMessage(); msg != "" {
		ui.drawLine(0, ui.height-1, msg, ui.palette.StyleForStatusMessage())
	}
}

func (ui *TerminalUI) drawLine(x, y int, text string, style tcell.Style) {
	col := x
	for _, ch := range text {
		if col >= ui.width {
			break
		}
		ui.screen.SetContent(col, y, ch, nil, style)
		col++
	}

	// Fill rest of line with spaces (for background color)
	for col < ui.width {
		ui.screen.SetContent(col, y, ' ', nil, style)
		col++
	}
}

func (ui *TerminalUI) handleKey(ev *tcell.EventKey) bool {
	switch ev.Key() {
	case tcell.KeyCtrlQ:
		return false // Quit

	case tcell.KeyBackspace, tcell.KeyBackspace2:
		ui.editor.Backspace()

	case tcell.KeyDelete:
		ui.editor.Delete()

	case tcell.KeyEnter:
		ui.editor.InsertChar('\n')

	case tcell.KeyLeft:
		ui.editor.MoveLeft()

	case tcell.KeyRight:
		ui.editor.MoveRight()

	case tcell.KeyUp:
		ui.editor.MoveUp()

	case tcell.KeyDown:
		ui.editor.MoveDown()

	case tcell.KeyHome:
		ui.editor.MoveToLineStart()

	case tcell.KeyEnd:
		ui.editor.MoveToLineEnd()

	case tcell.KeyCtrlA:
		ui.editor.MoveToStart()

	case tcell.KeyCtrlE:
		ui.editor.MoveToEnd()

	case tcell.KeyRune:
		ui.editor.InsertChar(ev.Rune())
	}

	return true
}
