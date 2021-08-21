package mp1

type lerBoardData struct {
	BlueUp bool
}

func lerRBRFork(g *Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.ExtraEvent = lerRedFork{player, moves}
	} else {
		g.Players[player].CurrentSpace = ChainSpace{5, 0}
	}
	return moves - 1
}

func lerRBFork(g *Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.Players[player].CurrentSpace = ChainSpace{4, 4}
	} else {
		g.Players[player].CurrentSpace = ChainSpace{4, 0}
	}
	return moves - 1
}

func lerBRFork1(g *Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.Players[player].CurrentSpace = ChainSpace{9, 0}
	} else {
		g.Players[player].CurrentSpace = ChainSpace{6, 10}
	}
	return moves - 1
}

func lerBRFork2(g *Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.Players[player].CurrentSpace = ChainSpace{7, 0}
	} else {
		g.Players[player].CurrentSpace = ChainSpace{6, 0}
	}
	return moves - 1
}

func lerBRFork3(g *Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.Players[player].CurrentSpace = ChainSpace{0, 0}
	} else {
		g.Players[player].CurrentSpace = ChainSpace{10, 0}
	}
	return moves - 1
}

func lerSwapGates(g *Game, player int) {
	bd := g.Board.Data.(lerBoardData)
	bd.BlueUp = !bd.BlueUp
	g.Board.Data = bd
}

func lerGotoIsland(space int) func(*Game, int) {
	return func(g *Game, player int) {
		g.Players[player].CurrentSpace = ChainSpace{8, space}
	}
}

func lerVisitRobot(g *Game, player, moves int) int {
	if g.Players[player].Coins >= 20 {
		g.ExtraEvent = lerRobot{player, moves}
	}
	return moves
}

var LER = Board{
	Chains: &[]Chain{
		{ //Start to first fork
			{Type: Blue},
			{Type: Red},
			{Type: Boo},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Start},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Red},
		},
		{ //Straight to red/blue/red fork
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: lerRBRFork},
		},
		{ //Offshoot to robot
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Invisible, PassingEvent: lerVisitRobot},
			{Type: Blue},
			{Type: Blue},
		},
		{ //Left of red/blue/red fork to red/blue fork
			{Type: Blue},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: BogusItem},
			{Type: Red},
			{Type: Happening, StoppingEvent: lerSwapGates},
			{Type: Invisible, PassingEvent: lerRBFork},
		},
		{ //Red/blue fork blue path
			{Type: Blue},
			{Type: Star},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Blue},
		},
		{ //Ahead of red/blue/red fork to blue/red fork 1
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Invisible, PassingEvent: lerBRFork1},
		},
		{ //Left of blue/red fork 1 to blue/red fork 2
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Happening, StoppingEvent: lerGotoIsland(0)},
			{Type: Happening, StoppingEvent: lerGotoIsland(1)},
			{Type: Invisible, PassingEvent: lerBRFork2},
		},
		{ //Past red gate of blue/red fork 2
			{Type: Blue},
		},
		{ //Top left Island
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Star},
		},
		{ //Right of blue/red fork 1 to blue/red fork 3
			{Type: Blue},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: lerBRFork3},
		},
		{ //Blue exit of blue/red fork 3
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Star},
			{Type: Bowser},
			{Type: Happening, StoppingEvent: lerSwapGates},
			{Type: Blue},
		},
		{ //Right of red/blue/red fork
			{Type: Blue},
			{Type: Blue},
			{Type: Red},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: lerVisitRobot},
			{Type: Happening, StoppingEvent: lerSwapGates},
		},
	},
	Links: &map[int]*[]ChainSpace{
		0:  {{1, 0}, {2, 0}},
		2:  {{1, 3}},
		4:  {{0, 0}},
		7:  {{0, 0}},
		8:  {{0, 0}},
		10: {{0, 0}},
		11: {{9, 7}},
	},
	BowserCoins: 19,
	Data:        lerBoardData{},
}
