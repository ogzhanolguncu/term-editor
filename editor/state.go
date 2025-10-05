package editor

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
)

type VimState struct {
	mode Mode
}

func NewVimState() *VimState {
	return &VimState{
		mode: ModeNormal,
	}
}
