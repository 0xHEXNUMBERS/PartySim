package mp1

type ytiThwompBranchEvent struct {
	Player int
	Moves  int
	Thwomp int
}

func (y ytiThwompBranchEvent) Responses() []Response {
	return []Response{true, false}
}

func (y ytiThwompBranchEvent) ControllingPlayer() int {
	return y.Player
}

func (y ytiThwompBranchEvent) Handle(r Response, g *Game) {
	pay := r.(bool)
	bd := g.Board.Data.(ytiBoardData)
	if pay {
		g.ExtraEvent = ytiPayThwompEvent{
			PayRangeEvent{
				y.Player,
				bd.Thwomps[y.Thwomp],
				min(50, g.Players[y.Player].Coins),
			},
			y.Moves,
			y.Thwomp,
		}
	} else {
		pos := bd.RejectThwompPos[y.Thwomp]
		g.Players[y.Player].CurrentSpace = pos
		g.MovePlayer(y.Player, y.Moves-1)
	}
}

type ytiPayThwompEvent struct {
	PayRangeEvent
	Moves  int
	Thwomp int
}

func (y ytiPayThwompEvent) Handle(r Response, g *Game) {
	y.PayRangeEvent.Handle(r, g)
	cost := r.(int)
	bd := g.Board.Data.(ytiBoardData)
	bd.Thwomps[y.Thwomp] = min(50, cost+1)
	pos := bd.AcceptThwompPos[y.Thwomp]
	g.Board.Data = bd
	g.Players[y.Player].CurrentSpace = pos
	g.MovePlayer(y.Player, y.Moves-1)
}
