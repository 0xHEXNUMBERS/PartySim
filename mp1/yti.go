package mp1

func ytiSwapStarPosition(g *Game) {
	s1, s2 := g.Board.Chains[1][19], g.Board.Chains[2][12]
	g.Board.Chains[1][19], g.Board.Chains[2][12] = s2, s1
}

func ytiCreateTwomp() func(game *Game, player, moves int) Event {
	coinsToPass := 1

	return func(game *Game, player, moves int) Event {
		playerPos := game.Players[player].CurrentSpace
		if game.Players[player].Coins >= coinsToPass {
			return BranchEvent{
				player,
				playerPos.Chain,
				moves,
			}
		}
		return nil
	}
}

var YTI = Board{
	Chains: []Chain{
		{ //Left island
			{Type: Blue}, //Branch #1 Dir A
			{Type: Happening, Event: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue}, //Branch #2 Dir B, links from before
			{Type: Blue},
			{Type: Minigame},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Chance},
			{Type: Red},
			{Type: Blue},
			{Type: Happening, Event: ytiSwapStarPosition},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Start},
			{Type: Happening, Event: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Minigame},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: ytiCreateTwomp()},
		},
		{ //Right island part 1
			{Type: Blue}, //Branch #2 Dir A
			{Type: Chance},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Happening, Event: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Blue}, //Branch #1 Dir B
			{Type: Red},
			{Type: Blue},
			{Type: Minigame},
			{Type: Blue},
			{Type: Happening, Event: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
			{Type: BlackStar},
			{Type: Happening, Event: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Minigame},
			{Type: Boo},
			{Type: Blue},
			{Type: Blue},
			{Type: Red},
			{Type: Happening, Event: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Invisible, PassingEvent: ytiCreateTwomp()},
		},
	},
	Links: map[int][]ChainSpace{
		0: {{1, 6}},
		1: {{0, 7}},
	},
}
