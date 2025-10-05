package main

import (
	"log"

	"github.com/ogzhanolguncu/go_editor/editor"
	termui "github.com/ogzhanolguncu/go_editor/term_ui"
)

func main() {
	editor, err := editor.New()
	if err != nil {
		log.Fatalf("failed to initialize editor: %v", err)
	}
	terminal, err := termui.NewTerminalUI(editor)
	if err != nil {
		log.Fatalf("failed to initialize editor: %v", err)
	}
	defer terminal.Close()
	terminal.Run()
}
