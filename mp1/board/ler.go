package board

import "github.com/0xhexnumbers/partysim/mp1"

//lerBoardData holds all of the board specific data related to LER.
type lerBoardData struct {
	BlueUp bool
}

//lerRBRFork handles the 3-way fork on the board.
func lerRBRFork(g *mp1.Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.NextEvent = LERRedFork{player, moves}
	} else {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(5, 0)
	}
	return moves - 1
}

//lerRBFork handles the Red/Blue fork on the board.
func lerRBFork(g *mp1.Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(4, 4)
	} else {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(4, 0)
	}
	return moves - 1
}

//lerBRFork1 handles the Blue/Red fork on the top-middle part of the board.
func lerBRFork1(g *mp1.Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(9, 0)
	} else {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(6, 10)
	}
	return moves - 1
}

//lerBRFork2 handles the Blue/Red fork on the top-left part of the board.
func lerBRFork2(g *mp1.Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(7, 0)
	} else {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(6, 0)
	}
	return moves - 1
}

//lerBRFork3 handles the Blue/Red fork on the top-right part of the board.
func lerBRFork3(g *mp1.Game, player, moves int) int {
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(0, 0)
	} else {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(10, 0)
	}
	return moves - 1
}

//lerSwapGates swaps which gates are up.
func lerSwapGates(g *mp1.Game, player int) {
	bd := g.Board.Data.(lerBoardData)
	bd.BlueUp = !bd.BlueUp
	g.Board.Data = bd
}

//lerGotoIsland sets the player's new position to the island at the
//top-left section of the board.
func lerGotoIsland(space int) func(*mp1.Game, int) {
	return func(g *mp1.Game, player int) {
		g.Players[player].CurrentSpace = mp1.NewChainSpace(8, space)
	}
}

//lerVisitRobot sets the next event to deciding to swap gates if the player
//has >= 20 coins.
func lerVisitRobot(g *mp1.Game, player, moves int) int {
	if g.Players[player].Coins >= 20 {
		g.NextEvent = LERRobot{mp1.Boolean{}, player, moves}
	}
	return moves
}

//LER holds the data for Luigi's Engine Room.
var LER = mp1.Board{
	Chains: &[]mp1.Chain{
		{ //Start to first fork
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Boo},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Start},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Red},
		},
		{ //Straight to red/blue/red fork
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: lerRBRFork},
		},
		{ //Offshoot to robot
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Invisible, PassingEvent: lerVisitRobot},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //Left of red/blue/red fork to red/blue fork
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.BogusItem},
			{Type: mp1.Red},
			{Type: mp1.Happening, StoppingEvent: lerSwapGates},
			{Type: mp1.Invisible, PassingEvent: lerRBFork},
		},
		{ //Red/blue fork blue path
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //Ahead of red/blue/red fork to blue/red fork 1
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Invisible, PassingEvent: lerBRFork1},
		},
		{ //Left of blue/red fork 1 to blue/red fork 2
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: lerGotoIsland(0)},
			{Type: mp1.Happening, StoppingEvent: lerGotoIsland(1)},
			{Type: mp1.Invisible, PassingEvent: lerBRFork2},
		},
		{ //Past red gate of blue/red fork 2
			{Type: mp1.Blue},
		},
		{ //Top left Island
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Star},
		},
		{ //Right of blue/red fork 1 to blue/red fork 3
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: lerBRFork3},
		},
		{ //Blue exit of blue/red fork 3
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Star},
			{Type: mp1.Bowser},
			{Type: mp1.Happening, StoppingEvent: lerSwapGates},
			{Type: mp1.Blue},
		},
		{ //Right of red/blue/red fork
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: lerVisitRobot},
			{Type: mp1.Happening, StoppingEvent: lerSwapGates},
		},
	},
	Links: &map[int]*[]mp1.ChainSpace{
		0:  {mp1.NewChainSpace(1, 0), mp1.NewChainSpace(2, 0)},
		2:  {mp1.NewChainSpace(1, 3)},
		4:  {mp1.NewChainSpace(0, 0)},
		7:  {mp1.NewChainSpace(0, 0)},
		8:  {mp1.NewChainSpace(0, 0)},
		10: {mp1.NewChainSpace(0, 0)},
		11: {mp1.NewChainSpace(9, 7)},
	},
	BowserCoins: 19,
	Data:        lerBoardData{},
}
