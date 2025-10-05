package screen

import (
	"github.com/gdamore/tcell/v2"
	"github.com/ogzhanolguncu/go_editor/editor"
)

type Palette struct {
	lineNumStyle          tcell.Style
	gutterStyle           tcell.Style
	currentLineStyle      tcell.Style
	normalModeStyle       tcell.Style
	insertModeStyle       tcell.Style
	statusBarMessageStyle tcell.Style
	normalTextStyle       tcell.Style
}

func NewPalette() *Palette {
	s := tcell.StyleDefault

	editorBg := tcell.NewRGBColor(16, 16, 16)
	currentLineBg := tcell.NewRGBColor(22, 22, 22)
	textColor := tcell.NewRGBColor(255, 255, 255)
	lineNumColor := tcell.NewRGBColor(80, 80, 80)

	warmOrange := tcell.NewRGBColor(255, 179, 102)
	mintGreen := tcell.NewRGBColor(153, 255, 228)
	darkMint := tcell.NewRGBColor(72, 134, 119)

	return &Palette{
		lineNumStyle:          s.Foreground(lineNumColor).Background(editorBg),
		gutterStyle:           s.Foreground(mintGreen).Background(editorBg).Dim(true),
		currentLineStyle:      s.Background(currentLineBg).Foreground(textColor),
		normalModeStyle:       s.Background(darkMint).Foreground(textColor),
		insertModeStyle:       s.Background(warmOrange).Foreground(editorBg),
		statusBarMessageStyle: s.Background(mintGreen).Foreground(editorBg),
		normalTextStyle:       s.Foreground(textColor).Background(editorBg),
	}
}

func (p *Palette) StyleForLineNum() tcell.Style {
	return p.lineNumStyle
}

func (p *Palette) StyleForGutter() tcell.Style {
	return p.gutterStyle
}

func (p *Palette) StyleForCurrentLine() tcell.Style {
	return p.currentLineStyle
}

func (p *Palette) StyleForStatusBar(mode editor.Mode) tcell.Style {
	if mode == editor.ModeInsert {
		return p.insertModeStyle
	}
	return p.normalModeStyle
}

func (p *Palette) StyleForStatusMessage() tcell.Style {
	return p.statusBarMessageStyle
}

func (p *Palette) StyleForNormalText() tcell.Style {
	return p.normalTextStyle
}
