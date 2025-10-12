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
	screen, err := screen.NewScreen(editor)
	if err != nil {
		log.Fatalf("failed to initialize editor: %v", err)
	}
	defer screen.Close()
	screen.Run()
}
