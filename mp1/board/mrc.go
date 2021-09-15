package board

import "github.com/0xhexnumbers/partysim/mp1"

//mrcBoardData holds all of the board specific data related to MRC.
type mrcBoardData struct {
	IsBowser bool
}

//mrcSwapCastleDir swaps toad/bowser's castle direction.
func mrcSwapCastleDir(g *mp1.Game, player int) {
	bd := g.Board.Data.(mrcBoardData)
	bd.IsBowser = !bd.IsBowser
	g.Board.Data = bd
}

//mrcVisitCastle handles the event when a player visits the castle.
func mrcVisitCastle(g *mp1.Game, player int, moves int) int {
	bd := g.Board.Data.(mrcBoardData)
	if bd.IsBowser {
		g.AwardCoins(player, -40, false)
	} else {
		if g.Players[player].Coins >= 20 {
			g.Players[player].Stars++
			g.AwardCoins(player, -20, false)
		}
	}
	g.Players[player].CurrentSpace = mp1.NewChainSpace(0, 0)
	mrcSwapCastleDir(g, player)
	return moves
}

//MRC holds the data for Mario's Rainbow Castle.
var MRC = mp1.Board{
	Chains: &[]mp1.Chain{
		{ //Start to first fork
			{Type: mp1.Invisible}, //Temp space so players can walk on Start
			{Type: mp1.Start},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //First fork: left
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Chance},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Boo},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
		},
		{ //First fork: right to second fork
			{Type: mp1.Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //Second fork: right
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: mp1.Blue},
		},
		{ //Second fork: left to end
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: mp1.Blue},
			{Type: mp1.Chance},
			{Type: mp1.Bowser},
			{Type: mp1.Invisible, PassingEvent: mrcVisitCastle},
		},
	},
	Links: &map[int]*[]mp1.ChainSpace{
		0: {mp1.NewChainSpace(1, 0), mp1.NewChainSpace(2, 0)},
		1: {mp1.NewChainSpace(2, 2)},
		2: {mp1.NewChainSpace(3, 0), mp1.NewChainSpace(4, 0)},
		3: {mp1.NewChainSpace(4, 3)},
	},
	Data: mrcBoardData{},
}
