package sandbox

func (s *Sandbox) Record(args ...string) {

}

type vhs struct{}

func (v vhs) Record(args ...string) {

}

func (v vhs) common() string {
	return `Set Theme "Snazzy"
Set FontSize 32
Set Width 1600
Set Height 400
Set Padding 0
Set CursorBlink false
Set TypingSpeed 0
`
}
