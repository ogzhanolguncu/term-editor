package highlighter

// # Chroma Integration Checklist
//
// ## **Setup**
// - [ ] `go get github.com/alecthomas/chroma/v2`
// - [ ] Create `highlighter/` package
//
// ## **1. highlighter/highlighter.go**
// - [ ] `Token` struct: Type, Value, Start, End
// - [ ] `Highlighter` struct: lexer, style, filename
// - [ ] `New(filename)` - auto-detect language
// - [ ] `TokenizeLine(lineNum, content)` - returns []Token
// - [ ] Internal `tokenize(content)` - chroma tokenization
// - [ ] `SetTheme(name)`, `GetTheme()`, `ListThemes()`
//
// ## **2. editor/editor.go**
// - [ ] Add `highlighter` field to Editor
// - [ ] Initialize in `New()`
// - [ ] `GetHighlighter()`, `SetTheme()`, `GetTheme()`, `ListThemes()` passthroughs
// - [ ] Invalidate on edits: `InvalidateLine()` in InsertChar/Delete, `InvalidateAll()` on newlines
//
// ## **3. screen/palette.go**
// - [ ] Add syntax style fields: keyword, string, comment, number, function, type, operator, identifier
// - [ ] Define colors in `NewPalette()`
// - [ ] `StyleForToken(chroma.TokenType)` - map tokens to styles
//
// ## **4. screen/screen.go**
// - [ ] Update `renderLines()`: get highlighter, call `TokenizeLine()` per line
// - [ ] For each char, find token, apply style + preserve current line background
//
// ## **5. screen/input.go**
// - [ ] Add `Ctrl-T` to cycle themes (temp testing)
//
// ## **6. Testing**
// - [ ] Test .go, .py, .txt files
// - [ ] Test typing, newlines, deletes
// - [ ] Test theme switching
// - [ ] Test large files for lag
//
// ## **Cache (Only if slow)**
// - [ ] `TokenCache` struct: map[int][]Token, maxSize
// - [ ] `Get()`, `Set()`, `InvalidateLine()`, `InvalidateFrom()`, `InvalidateAll()`
// - [ ] Add to Highlighter, check cache in `TokenizeLine()`
//
