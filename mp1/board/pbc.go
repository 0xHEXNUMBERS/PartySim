package board

import "github.com/0xhexnumbers/partysim/mp1"

//pbcBoardData holds all of the board specific data related to PBC.
type pbcBoardData struct {
	BowserSeedPlanted bool
	SeedCount         int
	PiranhaOccupied   [14]bool
	PiranhaPlant      [14]int
}

//pbcVisitSeed occurs when the player visit the seed game. If the seed the
//player will pick is known, then the decision is made for them. Otherwise,
//the next event is set to make that decision.
func pbcVisitSeed(g *mp1.Game, player, moves int) int {
	g.AwardCoins(player, -10, false)
	data := g.Board.Data.(pbcBoardData)
	data.SeedCount++
	if data.SeedCount > 4 { //Reset state
		data.SeedCount = 1
		data.BowserSeedPlanted = false
	}
	g.Board.Data = data
	if data.BowserSeedPlanted { //Can't receive bowser, goto star
		g.Players[player].CurrentSpace = mp1.NewChainSpace(0, 0)
	} else {
		if data.SeedCount == 4 { //Can't receive toad, goto bowser
			g.Players[player].CurrentSpace = mp1.NewChainSpace(1, 0)
		} else { //Could be either or
			g.NextEvent = PBCSeedCheck{mp1.Boolean{}, player, moves}
		}
	}
	return moves - 1
}

//pbcVisitPiranha occurs when a player lands on a happening space. If the
//space is occupied by another player and the current player has a star,
//then the current player will give their star to the occupying player.
//If the space is unoccupied, then the player can make a decision to
//occupy it if they have 30 coins.
func pbcVisitPiranha(piranha int) func(*mp1.Game, int) {
	return func(g *mp1.Game, player int) {
		data := g.Board.Data.(pbcBoardData)
		if data.PiranhaOccupied[piranha] {
			owner := data.PiranhaPlant[piranha]
			if owner != player && g.Players[player].Stars > 0 {
				g.Players[player].Stars--
				g.Players[owner].Stars++
			}
		} else if g.Players[player].Coins >= 30 {
			g.NextEvent = PBCPiranhaDecision{mp1.Boolean{}, player, piranha}
		}
	}
}

//PBC holds the data for Peach's Birthday Cake.
var PBC = mp1.Board{
	Chains: &[]mp1.Chain{
		{ //Main Path
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Star},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Chance},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(0)},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(1)},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(2)},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(3)},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(4)},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(5)},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Start},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(6)},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(7)},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(8)},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(9)},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(10)},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(11)},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(12)},
			{Type: mp1.Happening, StoppingEvent: pbcVisitPiranha(13)},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: pbcVisitSeed},
		},
		{ //Bowser Path
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Chance},
			{Type: mp1.BogusItem},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: pbcVisitSeed},
		},
	},
	Links:       nil,
	BowserCoins: 20,
	Data:        pbcBoardData{},
}
