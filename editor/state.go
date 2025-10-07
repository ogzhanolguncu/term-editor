package editor

import "strconv"

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
)

type VimState struct {
	mode         Mode
	commandCount string
}

func NewVimState() *VimState {
	return &VimState{
		mode:         ModeNormal,
		commandCount: "",
	}
}

func (v *VimState) HandleDigit(r rune) bool {
	// '1'-'9' always start or continue a count
	if r >= '1' && r <= '9' {
		v.commandCount += string(r)
		return true
	}
	// '0' only continues a count if one exists
	if r == '0' && v.commandCount != "" {
		v.commandCount += "0"
		return true
	}
	return false
}

func (v *VimState) GetCountAndClear() int {
	if v.commandCount == "" {
		return 1
	}
	n, _ := strconv.Atoi(v.commandCount)
	v.commandCount = ""
	return n
}

func (v *VimState) ClearCount() {
	v.commandCount = ""
}
