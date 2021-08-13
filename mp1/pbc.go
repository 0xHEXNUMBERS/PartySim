package mp1

type pbcBoardData struct {
	BowserSeedPlanted bool
	SeedCount         int
	PiranhaOccupied   [14]bool
	PiranhaPlant      [14]int
}

func pbcVisitSeed(g *Game, player, moves int) {
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
			g.ExtraEvent = pbcSeedCheck{player, moves}
		}
	}
}

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
			g.ExtraEvent = pbcPiranhaDecision{player, piranha}
		}
	}
}

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
