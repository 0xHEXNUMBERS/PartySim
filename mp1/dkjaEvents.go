package mp1

type dkjaWhompEvent struct {
	Player int
	Moves  int
	Whomp  int
}

func (d dkjaWhompEvent) Responses() []Response {
	return []Response{true, false}
}

func (d dkjaWhompEvent) ControllingPlayer() int {
	return d.Player
}

func (d dkjaWhompEvent) Handle(r Response, g *Game) {
	pay := r.(bool)
	data := g.Board.Data.(dkjaBoardData)
	if pay {
		g.AwardCoins(d.Player, -10, false)
	} else {
		data.WhompPos[d.Whomp] = !data.WhompPos[d.Whomp]
		g.Board.Data = data
	}
	pos := dkjaGetWhompDestination(g, d.Whomp)
	g.Players[d.Player].CurrentSpace = pos
	g.MovePlayer(d.Player, d.Moves-1)
}

type dkjaCoinBranchEvent struct {
	Player   int
	Moves    int
	Blockade int
}

func (d dkjaCoinBranchEvent) Responses() []Response {
	return []Response{true, false}
}

func (d dkjaCoinBranchEvent) ControllingPlayer() int {
	return d.Player
}

func (d dkjaCoinBranchEvent) Handle(r Response, g *Game) {
	accept := r.(bool)
	data := g.Board.Data.(dkjaBoardData)
	var pos ChainSpace
	if accept {
		pos = data.CoinAcceptDestination[d.Blockade]
	} else {
		pos = data.CoinRejectDestination[d.Blockade]
	}
	g.Players[d.Player].CurrentSpace = pos
	g.MovePlayer(d.Player, d.Moves-1)
}
