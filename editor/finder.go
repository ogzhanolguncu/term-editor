package editor

// Fuzzy Finder
// Why third: Fast file switching, complements tree
// Option A - Use fzf binary:
//
//  Shell out to fzf with file list
//  Parse selected file, open in buffer
//  Bind to Ctrl-p
//
// Option B - Pure Go
//
//  Use github.com/sahilm/fuzzy for matching
//  Popup UI: list files, filter as you type
//  Arrow keys navigate, Enter opens
//  Escape cancels
//  Show 10-15 results max
//  Search .git ignored files or all files (configurable)
//
