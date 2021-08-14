package mp1

type lerRedFork struct {
	Player int
	Moves  int
}

var lerRedForkDestinations = []Response{
	ChainSpace{3, 0},
	ChainSpace{11, 0},
}

func (l lerRedFork) Responses() []Response {
	return lerRedForkDestinations
}

func (l lerRedFork) ControllingPlayer() int {
	return l.Player
}

func (l lerRedFork) Handle(r Response, g *Game) {
	dest := r.(ChainSpace)
	g.Players[l.Player].CurrentSpace = dest
	g.MovePlayer(l.Player, l.Moves-1)
}

type lerRobot struct {
	Player int
	Moves  int
}

func (l lerRobot) Responses() []Response {
	return []Response{true, false}
}

func (l lerRobot) ControllingPlayer() int {
	return l.Player
}

func (l lerRobot) Handle(r Response, g *Game) {
	pay := r.(bool)
	if pay {
		g.AwardCoins(l.Player, -20, false)
		lerSwapGates(g, l.Player)
	}
	g.MovePlayer(l.Player, l.Moves)
}
