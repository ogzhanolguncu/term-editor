package main

import (
	"log"

	"github.com/ogzhanolguncu/go_editor/editor"
	"github.com/ogzhanolguncu/go_editor/screen"
)

func main() {
	editor, err := editor.New()
	if err != nil {
		log.Fatalf("failed to initialize editor: %v", err)
	}
	terminal, err := screen.NewScreen(editor)
	if err != nil {
		log.Fatalf("failed to initialize editor: %v", err)
	}
	defer terminal.Close()
	terminal.Run()
}
