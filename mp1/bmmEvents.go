package mp1

type bmmBranchPay struct {
	Player     int
	Moves      int
	BowserPath ChainSpace
	StarPath   ChainSpace
}

func (b bmmBranchPay) Responses() []Response {
	return []Response{true, false}
}

func (b bmmBranchPay) ControllingPlayer() int {
	return b.Player
}

func (b bmmBranchPay) Handle(r Response, g *Game) {
	pay := r.(bool)
	if pay {
		g.AwardCoins(b.Player, -10, false)
		g.ExtraEvent = bmmBranchDecision{
			b.Player, b.Moves, b.BowserPath, b.StarPath,
		}
	} else {
		g.Players[b.Player].CurrentSpace = b.BowserPath
		g.MovePlayer(b.Player, b.Moves-1)
	}
}

type bmmBranchDecision struct {
	Player     int
	Moves      int
	BowserPath ChainSpace
	StarPath   ChainSpace
}

func (b bmmBranchDecision) Responses() []Response {
	return []Response{b.BowserPath, b.StarPath}
}

func (b bmmBranchDecision) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b bmmBranchDecision) Handle(r Response, g *Game) {
	dest := r.(ChainSpace)
	g.Players[b.Player].CurrentSpace = dest
	g.MovePlayer(b.Player, b.Moves-1)
}

type bmmBowserRoulette struct {
	Player int
	Moves  int
}

func (b bmmBowserRoulette) Responses() []Response {
	return []Response{true, false}
}

func (b bmmBowserRoulette) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b bmmBowserRoulette) Handle(r Response, g *Game) {
	starSteal := r.(bool)
	if starSteal {
		g.Players[b.Player].Stars--
	} else {
		g.AwardCoins(b.Player, -20, false)
	}
	g.MovePlayer(b.Player, b.Moves)
}
