# Terminal Text Editor in Go

A minimal, functional terminal text editor built from scratch in Go. Handles raw keyboard input, ANSI escape sequences, and real-time screen rendering without external dependencies.

## Architecture

```
┌─────────────┐    ┌──────────────┐    ┌────────────┐
│ Input Loop  │───▶│ TextBuffer   │───▶│ Screen     │
│ (raw mode)  │    │ ([]string)   │    │ (display)  │
└─────────────┘    └──────────────┘    └────────────┘
```

**TextBuffer**: Stores text as slice of strings (one per line)
**Screen**: Handles ANSI positioning and buffered output
**Input System**: Parses keyboard events and escape sequences

## Controls

- **Arrow Keys**: Navigate cursor
- **Enter**: New line
- **Printable characters**: Insert at cursor position
- **Ctrl+C**: Exit editor

## Next Steps

**Core Editing:**
- [ ] Backspace and Delete key support
- [ ] Tab character handling
- [ ] Cut/Copy/Paste operations

**UI Components:**
- [ ] Status bar with file info and cursor position
- [ ] Line numbers in left gutter
- [ ] Message bar for notifications and prompts
- [ ] Welcome screen for new files

**File Operations:**
- [ ] File loading and saving (Ctrl+O, Ctrl+S)
- [ ] New file creation
- [ ] Dirty buffer tracking with modified indicator

**Editor Features:**
- [ ] Scrolling for large files (viewport management)
- [ ] Search and replace (Ctrl+F, highlighting matches)
- [ ] Undo/Redo system
- [ ] Find and goto line (Ctrl+G)

**Performance & Polish:**
- [ ] Gap buffer implementation for efficient insertions
- [ ] Dirty line tracking to minimize redraws
- [ ] Syntax highlighting framework
- [ ] Configuration system (tabstop, colors, keybindings)
