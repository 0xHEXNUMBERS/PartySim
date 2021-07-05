package mp1

type PayThwompEvent struct {
	PayRangeEvent
	Thwomp int
	Link   ChainSpace
}

func (p PayThwompEvent) Handle(r Response, g Game) Game {
	g = ResetGameExtras(g)
	g = p.PayRangeEvent.Handle(r, g)

	cost := r.(int)
	bd := g.Board.Data.(ytiBoardData)
	bd.Thwomps[p.Thwomp] = cost + 1
	g.Board.Data = bd
	g.Players[p.Player].CurrentSpace = p.Link
	g.ExtraMovement = Movement{p.Player, p.Moves - 1, false}
	return g
}
