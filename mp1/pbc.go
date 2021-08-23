package mp1

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
func pbcVisitSeed(g *Game, player, moves int) int {
	g.AwardCoins(player, -10, false)
	data := g.Board.Data.(pbcBoardData)
	data.SeedCount++
	if data.SeedCount > 4 { //Reset state
		data.SeedCount = 1
		data.BowserSeedPlanted = false
	}
	g.Board.Data = data
	if data.BowserSeedPlanted { //Can't recieve bowser, goto star
		g.Players[player].CurrentSpace = ChainSpace{0, 0}
	} else {
		if data.SeedCount == 4 { //Can't recieve toad, goto bowser
			g.Players[player].CurrentSpace = ChainSpace{1, 0}
		} else { //Could be either or
			g.NextEvent = pbcSeedCheck{player, moves}
		}
	}
	return moves - 1
}

//pbcVisitPiranha occurs when a player lands on a happening space. If the
//space is occupied by another player and the current player has a star,
//then the current player will give their star to the occupying player.
//If the space is unoccupied, then the player can make a decision to
//occupy it if they have 30 coins.
func pbcVisitPiranha(piranha int) func(*Game, int) {
	return func(g *Game, player int) {
		data := g.Board.Data.(pbcBoardData)
		if data.PiranhaOccupied[piranha] {
			owner := data.PiranhaPlant[piranha]
			if owner != player && g.Players[player].Stars > 0 {
				g.Players[player].Stars--
				g.Players[owner].Stars++
			}
		} else if g.Players[player].Coins >= 30 {
			g.NextEvent = pbcPiranhaDecision{player, piranha}
		}
	}
}

//PBC holds the data for Peach's Birthday Cake.
var PBC = Board{
	Chains: &[]Chain{
		{ //Main Path
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Star},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Chance},
			{Type: Blue},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Red},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(0)},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(1)},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(2)},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(3)},
			{Type: Blue},
			{Type: Red},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(4)},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(5)},
			{Type: Blue},
			{Type: Blue},
			{Type: Start},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(6)},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(7)},
			{Type: Red},
			{Type: Blue},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(8)},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(9)},
			{Type: Blue},
			{Type: Blue},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(10)},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(11)},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(12)},
			{Type: Happening, StoppingEvent: pbcVisitPiranha(13)},
			{Type: Blue},
			{Type: Invisible, PassingEvent: pbcVisitSeed},
		},
		{ //Bowser Path
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
			{Type: Mushroom},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
			{Type: Chance},
			{Type: BogusItem},
			{Type: Red},
			{Type: Blue},
			{Type: Invisible, PassingEvent: pbcVisitSeed},
		},
	},
	Links:       nil,
	BowserCoins: 20,
	Data:        pbcBoardData{},
}
