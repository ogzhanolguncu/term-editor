package editor

// After Syntax HighlightingNext 4 Features (in order)1. Multiple Buffers + Buffer Switching
// Why first: Can't be a real editor without opening multiple files
//  Buffer struct: holds TextBuffer + metadata (filename, modified, cursor pos)
//  BufferManager in editor: map of buffers, active buffer ID
//  :e <file> - open file in new buffer
//  :bn, :bp - next/previous buffer
//  :bd - close buffer (warn if modified)
//  :ls or :buffers - list open buffers
//  Ctrl-6 - switch to last buffer
//  Status bar shows buffer name + modified flag
//  Handle closing last buffer
