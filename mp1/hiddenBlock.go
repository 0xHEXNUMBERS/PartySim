package mp1

type HiddenBlockEvent struct {
	Player int
}

func (h HiddenBlockEvent) Responses() []Response {
	return []Response{true, false}
}

func (h HiddenBlockEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (h HiddenBlockEvent) Handle(r Response, g *Game) {
	isHiddenBlock := r.(bool)
	if isHiddenBlock {
		g.ExtraEvent = EventDiceBlock{h.Player}
	} else {
		g.ActivateSpace(h.Player)
	}
}
