package mp1

import "fmt"

type HiddenBlockResponse int

const (
	HiddenBlockAppears HiddenBlockResponse = iota
	HiddenBlockNotThere
)

func (h HiddenBlockResponse) String() string {
	switch h {
	case HiddenBlockAppears:
		return "A hidden block appears"
	case HiddenBlockNotThere:
		return "There is no hidden block"
	}
	return ""
}

//HiddenBlockEvent holds the implementation for hidden blocks.
type HiddenBlockEvent struct {
	Player int
}

func (h HiddenBlockEvent) Question(g *Game) string {
	return fmt.Sprintf("Did %s land on a hidden event block?",
		g.Players[h.Player].Char)
}

func (h HiddenBlockEvent) Type() EventType {
	return ENUM_EVT_TYPE
}

func (h HiddenBlockEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (h HiddenBlockEvent) Responses() []Response {
	return []Response{HiddenBlockAppears, HiddenBlockNotThere}
}

//Handle sets the hidden block action to be taken depending on r. If r is
//true, then the next event will be a EventDiceBlock. If r is false, then
//then the player will land on the space they're currently on.
func (h HiddenBlockEvent) Handle(r Response, g *Game) {
	isHiddenBlock := r.(HiddenBlockResponse)
	if isHiddenBlock == HiddenBlockAppears {
		g.NextEvent = EventDiceBlock{h.Player}
	} else {
		g.ActivateSpace(h.Player)
	}
}
