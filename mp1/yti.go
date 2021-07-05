package mp1

type ytiBoardData struct {
	Thwomps      [2]int
	StarPosition bool
}

func (y ytiBoardData) Copy() ExtraBoardData {
	return y
}

func (y *ytiBoardData) SwapStarPosition(g *Game) {
	y.StarPosition = !y.StarPosition

	if y.StarPosition {
		g.Board.Chains[0][19].Type = Star
		g.Board.Chains[1][18].Type = BlackStar
	} else {
		g.Board.Chains[0][19].Type = BlackStar
		g.Board.Chains[1][18].Type = Star
	}
}

func (y ytiBoardData) CanPassThwomp(game *Game,
	player, moves, thwomp int) Event {
	playerPos := game.Players[player].CurrentSpace
	if game.Players[player].Coins >= y.Thwomps[thwomp] {
		return BranchEvent{
			player,
			playerPos.Chain,
			moves,
			game.Board.Links[playerPos.Chain],
		}
	}
	return nil
}

func ytiCheckThwomp(thwomp int) func(*Game, int, int) Event {
	return func(g *Game, player, moves int) Event {
		bd := g.Board.Data.(ytiBoardData)
		return bd.CanPassThwomp(g, player, moves, thwomp)
	}
}

func ytiPayThwomp(thwomp int) func(*Game, int, int) Event {
	return func(g *Game, player, moves int) Event {
		bd := g.Board.Data.(ytiBoardData)
		return PayThwompEvent{
			PayRangeEvent{
				player,
				bd.Thwomps[thwomp],
				g.Players[player].Coins,
				moves,
			},
			thwomp,
			g.Board.Links[thwomp+2][0],
		}
	}
}

func ytiSwapStarPosition(g *Game) Event {
	bd := g.Board.Data.(ytiBoardData)
	bd.SwapStarPosition(g)
	return nil
}

var YTI = Board{
	Chains: []Chain{
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
			{Type: BlackStar},
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
	Links: map[int][]ChainSpace{
		0: {{2, 0}},
		1: {{3, 0}},
		2: {{1, 6}}, //Thwomp payments only have 1 link
		3: {{0, 7}}, //Thwomp payments only have 1 link
	},
	Data: ytiBoardData{[2]int{1, 1}, true},
}
