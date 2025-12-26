package messages

type ModeChangedMessage struct {
	Mode int
}
func (m ModeChangedMessage) IsModeSearch() bool {
	return true
}
