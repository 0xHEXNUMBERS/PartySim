package mp1

type pbcSeedCheck struct {
	Player int
	Moves  int
}

func (p pbcSeedCheck) Responses() []Response {
	return []Response{true, false}
}

func (p pbcSeedCheck) ControllingPlayer() int {
	return CPU_PLAYER
}

func (p pbcSeedCheck) Handle(r Response, g *Game) {
	bowser := r.(bool)
	if bowser {
		g.Players[p.Player].CurrentSpace = ChainSpace{1, 0}
		data := g.Board.Data.(pbcBoardData)
		data.BowserSeedPlanted = true
		g.Board.Data = data
	} else {
		g.Players[p.Player].CurrentSpace = ChainSpace{0, 0}
	}
	g.MovePlayer(p.Player, p.Moves-1)
}

type pbcPiranhaDecision struct {
	Player  int
	Piranha int
}

func (p pbcPiranhaDecision) Responses() []Response {
	return []Response{true, false}
}

func (p pbcPiranhaDecision) ControllingPlayer() int {
	return p.Player
}

func (p pbcPiranhaDecision) Handle(r Response, g *Game) {
	plantPiranha := r.(bool)
	data := g.Board.Data.(pbcBoardData)
	if plantPiranha {
		data.PiranhaPlant[p.Piranha] = p.Player
		data.PiranhaOccupied[p.Piranha] = true
		g.AwardCoins(p.Player, -30, false)
	}
	g.Board.Data = data
	g.EndCharacterTurn()
}
