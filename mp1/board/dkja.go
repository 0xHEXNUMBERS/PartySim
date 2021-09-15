package board

import "github.com/0xhexnumbers/partysim/mp1"

//dkjaBoardData holds all of the board specific data related to DKJA.
type dkjaBoardData struct {
	WhompPos                 [3]bool
	WhompMainDestination     [3]mp1.ChainSpace
	WhompOffshootDestination [3]mp1.ChainSpace
	CoinAcceptDestination    [2]mp1.ChainSpace
	CoinRejectDestination    [2]mp1.ChainSpace
}

//dkjaGetWhompDestination gets the mp1.ChainSpace that the whomp is not
//currently blocking.
func dkjaGetWhompDestination(g *mp1.Game, whomp int) mp1.ChainSpace {
	data := g.Board.Data.(dkjaBoardData)
	var pos mp1.ChainSpace
	if data.WhompPos[whomp] {
		pos = data.WhompOffshootDestination[whomp]
	} else {
		pos = data.WhompMainDestination[whomp]
	}
	return pos
}

//dkjaCanPassWhomp checks to see if the player can pay the toll to pass
//the whomp. If so, the next event is set for the player to make that
//decision.
func dkjaCanPassWhomp(whomp int) func(*mp1.Game, int, int) int {
	return func(g *mp1.Game, player, moves int) int {
		if g.Players[player].Coins >= 10 {
			g.NextEvent = DKJAWhompEvent{
				player, moves, whomp,
			}
		} else {
			pos := dkjaGetWhompDestination(g, whomp)
			g.Players[player].CurrentSpace = pos
		}
		return moves - 1
	}
}

//dkjaCanPassCoinBlockade checks to see if the player has the amount of
//coins to pass the blockade. If so, the next event is set for the player
//to make that decision.
func dkjaCanPassCoinBlockade(blockade int) func(*mp1.Game, int, int) int {
	return func(g *mp1.Game, player, moves int) int {
		if g.Players[player].Coins >= 20 {
			g.NextEvent = DKJACoinBranchEvent{
				player, moves, blockade,
			}
		} else {
			data := g.Board.Data.(dkjaBoardData)
			g.Players[player].CurrentSpace = data.CoinRejectDestination[blockade]
		}
		return moves - 1
	}
}

//dkjaBoulder moves any players on the boulder's path to the end of the
//path.
func dkjaBoulder(g *mp1.Game, player int) {
	for i := 0; i < 4; i++ {
		pos := g.Players[i].CurrentSpace
		if pos.Chain == 7 || (pos.Chain == 5 && pos.Space != 0) {
			g.Players[i].CurrentSpace = mp1.NewChainSpace(0, 16)
		}
	}
}

//DKJA holds the data for Donkey Kong's Jungle Adventure.
var DKJA = mp1.Board{
	Chains: &[]mp1.Chain{
		{ //Last Offshoot to first thwomp fork
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Happening, StoppingEvent: dkjaBoulder},
			{Type: mp1.Happening, StoppingEvent: dkjaBoulder},
			{Type: mp1.Blue},
			{Type: mp1.Boo},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Start},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: dkjaCanPassWhomp(0)},
		},
		{ //First offshoot to coin blockade
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: dkjaCanPassCoinBlockade(0)},
		},
		{ //Through first coin blockade
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Red},
		},
		{ //Around first coin blockade
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: dkjaBoulder},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
		},
		{ //First main pathway
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: dkjaCanPassWhomp(1)},
		},
		{ //Second main pathway
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: dkjaBoulder},
			{Type: mp1.Happening, StoppingEvent: dkjaBoulder},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: dkjaBoulder},
			{Type: mp1.Boo},
			{Type: mp1.Invisible, PassingEvent: dkjaCanPassCoinBlockade(1)},
		},
		{ //Second offshoot pathway
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //Through second coin blockade
			{Type: mp1.MinigameSpace},
			{Type: mp1.Red},
			{Type: mp1.BogusItem},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: dkjaBoulder},
			{Type: mp1.Happening, StoppingEvent: dkjaBoulder},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //Around second coin blockade
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: dkjaCanPassWhomp(2)},
		},
		{ //Third Main Pathway
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Happening, StoppingEvent: dkjaBoulder},
			{Type: mp1.Blue},
		},
	},
	Links: &map[int]*[]mp1.ChainSpace{
		2: {mp1.NewChainSpace(0, 18)},
		3: {mp1.NewChainSpace(0, 16)},
		6: {mp1.NewChainSpace(5, 9)},
		7: {mp1.NewChainSpace(0, 15)},
		9: {mp1.NewChainSpace(0, 13)},
	},
	BowserCoins: 10,
	Data: dkjaBoardData{
		WhompPos: [3]bool{false, false, false},
		WhompMainDestination: [3]mp1.ChainSpace{
			mp1.NewChainSpace(4, 0),
			mp1.NewChainSpace(5, 0),
			mp1.NewChainSpace(9, 0),
		},
		WhompOffshootDestination: [3]mp1.ChainSpace{
			mp1.NewChainSpace(1, 0),
			mp1.NewChainSpace(6, 0),
			mp1.NewChainSpace(0, 0),
		},
		CoinAcceptDestination: [2]mp1.ChainSpace{
			mp1.NewChainSpace(2, 0),
			mp1.NewChainSpace(7, 0),
		},
		CoinRejectDestination: [2]mp1.ChainSpace{
			mp1.NewChainSpace(3, 0),
			mp1.NewChainSpace(8, 0),
		},
	},
}
