package mp1

type ytiBoardData struct {
	Thwomps      [2]int
	StarPosition ChainSpace
}

func ytiCheckThwomp(thwomp int) func(Game, int, int) Game {
	return func(g Game, player, moves int) Game {
		bd := g.Board.Data.(ytiBoardData)
		playerPos := g.Players[player].CurrentSpace
		if g.Players[player].Coins >= bd.Thwomps[thwomp] {
			g.ExtraEvent = BranchEvent{
				player,
				playerPos.Chain,
				moves,
				(*g.Board.Links)[playerPos.Chain],
			}
			return g
		}
		return g
	}
}

func ytiPayThwomp(thwomp int) func(Game, int, int) Game {
	return func(g Game, player, moves int) Game {
		bd := g.Board.Data.(ytiBoardData)
		maxCoins := 50
		if maxCoins > g.Players[player].Coins {
			maxCoins = g.Players[player].Coins
		}
		g.ExtraEvent = PayThwompEvent{
			PayRangeEvent{
				player,
				bd.Thwomps[thwomp],
				maxCoins,
			},
			thwomp,
			(*(*g.Board.Links)[thwomp+2])[0],
			moves,
		}
		return g
	}
}

var ytiLeftIslandStar = ChainSpace{0, 19}
var ytiRightIslandStar = ChainSpace{1, 18}

func ytiSwapStarPosition(g Game) Game {
	bd := g.Board.Data.(ytiBoardData)
	if bd.StarPosition == ytiLeftIslandStar {
		bd.StarPosition = ytiRightIslandStar
	} else {
		bd.StarPosition = ytiLeftIslandStar
	}
	g.Board.Data = bd
	return g
}

func ytiGainStar(g Game, player, moves int) Game {
	bd := g.Board.Data.(ytiBoardData)
	if bd.StarPosition == g.Players[player].CurrentSpace {
		if g.Players[player].Coins >= 20 {
			g = AwardCoins(g, player, -20, false)
			g.Players[player].Stars++
			g = ytiSwapStarPosition(g)
		}
	} else { //Star at other island
		g = AwardCoins(g, player, -30, false)
	}
	return g
}

var YTI = Board{
	Chains: &[]Chain{
		{ //Left island
			{Type: Blue}, //Branch #1 Dir A
			{Type: Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue}, //Branch #2 Dir B, links from before
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Chance},
			{Type: Red},
			{Type: Blue},
			{Type: Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: Star, PassingEvent: ytiGainStar},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Start},
			{Type: Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: ytiCheckThwomp(0)},
		},
		{ //Right island part 1
			{Type: Blue}, //Branch #2 Dir A
			{Type: Chance},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Blue}, //Branch #1 Dir B
			{Type: Red},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
			{Type: Star, PassingEvent: ytiGainStar},
			{Type: Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Boo},
			{Type: Blue},
			{Type: Blue},
			{Type: Red},
			{Type: Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Invisible, PassingEvent: ytiCheckThwomp(1)},
		},
		{ //Temporary place for thwomp payment
			{Type: Invisible, PassingEvent: ytiPayThwomp(0)},
		},
		{ //Temporary place for thwomp payment
			{Type: Invisible, PassingEvent: ytiPayThwomp(1)},
		},
	},
	Links: &map[int]*[]ChainSpace{
		0: {{2, 0}},
		1: {{3, 0}},
		2: {{1, 6}}, //Thwomp payments only have 1 link
		3: {{0, 7}}, //Thwomp payments only have 1 link
	},
	Data: ytiBoardData{[2]int{1, 1}, ytiLeftIslandStar},
}
