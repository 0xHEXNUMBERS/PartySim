package mp1

type PayThwompEvent struct {
	PayRangeEvent
	Thwomp int
	Link   ChainSpace
	Moves  int
}

func (p PayThwompEvent) Handle(r Response, g Game) Game {
	g = p.PayRangeEvent.Handle(r, g)
	cost := r.(int)
	bd := g.Board.Data.(ytiBoardData)
	bd.Thwomps[p.Thwomp] = min(50, cost+1)
	g.Board.Data = bd
	g.Players[p.Player].CurrentSpace = p.Link
	g = MovePlayer(g, p.Player, p.Moves-1)
	return g
}
