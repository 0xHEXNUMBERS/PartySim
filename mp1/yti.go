package mp1

type ytiBoardData struct {
	Thwomps      [2]int
	StarPosition bool
}

func (y ytiBoardData) Copy() ExtraBoardData {
	return y
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
				g.Board.Links[playerPos.Chain],
			}
			return g
		}
		return g
	}
}

func ytiPayThwomp(thwomp int) func(Game, int, int) Game {
	return func(g Game, player, moves int) Game {
		bd := g.Board.Data.(ytiBoardData)
		g.ExtraEvent = PayThwompEvent{
			PayRangeEvent{
				player,
				bd.Thwomps[thwomp],
				g.Players[player].Coins,
				moves,
			},
			thwomp,
			(*g.Board.Links[thwomp+2])[0],
		}
		return g
	}
}

func ytiSwapStarPosition(g Game) Game {
	bd := g.Board.Data.(ytiBoardData)
	bd.StarPosition = !bd.StarPosition
	g.Board.Data = bd
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
			{Type: Star},
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
			{Type: Star},
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
	Links: map[int]*[]ChainSpace{
		0: {{2, 0}},
		1: {{3, 0}},
		2: {{1, 6}}, //Thwomp payments only have 1 link
		3: {{0, 7}}, //Thwomp payments only have 1 link
	},
	Data: ytiBoardData{[2]int{1, 1}, true},
}
