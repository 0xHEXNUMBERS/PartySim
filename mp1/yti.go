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
		g.Board.Chains[0][23].Type = Star
		g.Board.Chains[1][18].Type = BlackStar
	} else {
		g.Board.Chains[0][23].Type = BlackStar
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
			{Type: Minigame},
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
			{Type: Minigame},
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
			{Type: Minigame},
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
			{Type: Minigame},
			{Type: Boo},
			{Type: Blue},
			{Type: Blue},
			{Type: Red},
			{Type: Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: Blue},
			{Type: Invisible, PassingEvent: ytiCheckThwomp(1)},
		},
	},
	Links: map[int][]ChainSpace{
		0: {{1, 6}},
		1: {{0, 7}},
	},
	Data: ytiBoardData{[2]int{1, 1}, true},
}
