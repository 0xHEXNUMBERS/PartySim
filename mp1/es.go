package mp1

type esBoardData struct {
	StarTaken [7]bool
	Gate      int //0 if unknown
	Gate2or3  bool
}

var esEntrance1 = ChainSpace{0, 0}
var esStartingSpace = ChainSpace{0, 1}
var esEntrance4 = ChainSpace{7, 0}
var esEntrance5 = ChainSpace{8, 0}
var esEntrance6 = ChainSpace{9, 0}
var esEntrance7 = ChainSpace{10, 0}
var esEntrance8 = ChainSpace{11, 0}
var esEntrance9 = ChainSpace{12, 0}

func esVisitBowser(g *Game, player, moves int) int {
	if g.Players[player].Stars > 0 {
		g.Players[player].Stars--
	} else {
		g.AwardCoins(player, -20, false)
	}
	return moves
}

func esSendToStart(g *Game, player int) {
	for i := range g.Players {
		g.Players[i].CurrentSpace = esStartingSpace
	}
	//For some reason, happening also gives you 3/6 coins
	//because the game thinks you landed on a blue space as well.
	//Happening space, thankfully, still count towards happening
	//star.
	if g.LastFiveTurns() {
		g.AwardCoins(player, 6, false)
	} else {
		g.AwardCoins(player, 3, false)
	}
}

func esWarpC(g *Game, player, moves int) int {
	bd := g.Board.Data.(esBoardData)
	if bd.Gate2or3 {
		g.Players[player].CurrentSpace = esStartingSpace
	} else {
		switch bd.Gate {
		case 0:
			g.ExtraEvent = esWarpCDest{player, moves}
		case 1:
			g.Players[player].CurrentSpace = esEntrance7
		default:
			g.Players[player].CurrentSpace = esEntrance1
		}
	}
	return moves
}

func esWarpSpace(dest1, dest2, dest3 ChainSpace) func(*Game, int, int) int {
	return func(g *Game, player, moves int) int {
		bd := g.Board.Data.(esBoardData)
		switch bd.Gate {
		case 0:
			g.ExtraEvent = esWarpDest{
				player, moves, bd.Gate2or3,
				dest1,
				dest2,
				dest3,
			}
		case 1:
			g.Players[player].CurrentSpace = dest1
		case 2:
			g.Players[player].CurrentSpace = dest2
		case 3:
			g.Players[player].CurrentSpace = dest3
		}
		return moves
	}
}

func esBranchWithWarp(dest1, dest2, dest3 ChainSpace) func(*Game, int, int) int {
	return func(g *Game, player, moves int) int {
		g.ExtraEvent = esBranchEvent{
			player,
			moves,
			dest1,
			dest2,
			dest3,
		}
		return moves
	}
}

func esAllStarsCollected(e esBoardData) bool {
	for _, star := range e.StarTaken {
		if !star {
			return false
		}
	}
	return true
}

func esPassStarSpace(i int) func(*Game, int, int) int {
	return func(g *Game, player, moves int) int {
		bd := g.Board.Data.(esBoardData)
		if !bd.StarTaken[i] {
			if g.Players[player].Coins >= 20 {
				g.ExtraEvent = esVisitBabyBowser{
					player,
					moves,
					i,
				}
			}
			return moves
		}
		return moves - 1
	}
}

func esLandOnChanceTime(g *Game, player int) {
	g.Players[player].LastSpaceType = Chance
}

var ES = Board{
	Chains: &[]Chain{
		{ //0: Entrance 1
			{Type: Invisible}, //Warp 1 Entrance
			{Type: Blue},
			{Type: MinigameSpace},
		},
		{ //1: Entrance 1 Fork: Right Exit Through warp A to Entrance 2 Fork
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
		},
		{ //2: Entrance 1 Fork: Left Exit Through warp B to Entrance 3 Fork
			{Type: Blue},
			{Type: Red},
			{Type: Blue},
			{Type: Blue},
		},
		{ //3: Entrance 2 Fork: Left Exit to warp C
			{Type: Blue},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Bowser},
			{Type: Invisible, PassingEvent: esWarpC},
		},
		{ //4: Entrance 2 Fork: Right Exit to warp D
			{Type: Blue},
			{Type: Blue},
			{
				Type:          Invisible,
				PassingEvent:  esPassStarSpace(0),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: Blue},
			{Type: Happening, StoppingEvent: esSendToStart},
			{
				Type: Invisible,
				PassingEvent: esWarpSpace(
					esEntrance1,
					esEntrance7,
					esEntrance6,
				),
			},
		},
		{ //5: Entrance 3 Fork: Right exit to warp E
			{Type: Red},
			{
				Type:          Invisible,
				PassingEvent:  esPassStarSpace(1),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{
				Type: Invisible,
				PassingEvent: esWarpSpace(
					esEntrance6,
					esEntrance4,
					esEntrance1,
				),
			},
		},
		{ //6: Entrance 3 Fork: Left exit to warp F
			{Type: Blue},
			{Type: Happening, StoppingEvent: esSendToStart},
			{
				Type: Invisible,
				PassingEvent: esWarpSpace(
					esEntrance1,
					esEntrance6,
					esEntrance7,
				),
			},
		},
		{ //7: Entrance 4 to Warp H with Warp G branch
			{Type: Invisible}, //Movement space
			{Type: Blue},
			{
				Type:          Invisible,
				PassingEvent:  esPassStarSpace(2),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: Blue},
			{Type: Bowser},
			{ //Warp G
				Type: Invisible,
				PassingEvent: esBranchWithWarp(
					esEntrance9,
					esEntrance9,
					esEntrance8,
				),
			},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{
				Type:          Invisible,
				PassingEvent:  esPassStarSpace(3),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: Blue},
			{
				Type: Invisible,
				PassingEvent: esWarpSpace(
					esEntrance5,
					esEntrance8,
					esEntrance5,
				),
			},
		},
		{ //8: Entrance 5 to Warp I (always goes to start
			{Type: Invisible}, //Movement
			{Type: Red},
			{Type: Happening, StoppingEvent: esSendToStart},
			{
				Type:          Invisible,
				PassingEvent:  esPassStarSpace(4),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: Blue},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
		},
		{ //9: Entrance 6 to warp K with Warp J Branch
			{Type: Invisible}, //Tmp space for warp
			{Type: Blue},
			{
				Type:          Invisible,
				PassingEvent:  esPassStarSpace(5),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: Blue},
			{Type: Blue},
			{Type: Red},
			{Type: Blue},
			{Type: Blue},
			{
				Type:          Invisible,
				PassingEvent:  esPassStarSpace(6),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: Happening, StoppingEvent: esSendToStart},
			{ //Warp J
				Type: Invisible,
				PassingEvent: esBranchWithWarp(
					esEntrance9,
					esEntrance9,
					esEntrance8,
				),
			},
			{Type: Blue},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{
				Type: Invisible,
				PassingEvent: esWarpSpace(
					esEntrance8,
					esEntrance1,
					esEntrance4,
				),
			},
		},
		{ //10: Entrance 7 to Entrance 6 convergence
			{Type: Invisible}, //Tmp space for movement
			{Type: Blue},
			{Type: Blue},
		},
		{ //11: Entrance 8
			{Type: Invisible}, //Tmp space for movement
			{Type: Boo},
		},
		{ //12: Entrance 9
			{Type: Invisible}, //Tmp space for movement
			{Type: Invisible, PassingEvent: esVisitBowser}, //Bowser
		},
	},
	Links: &map[int]*[]ChainSpace{
		0:  {{1, 0}, {2, 0}},
		1:  {{3, 0}, {4, 0}},
		2:  {{5, 0}, {6, 0}},
		8:  {{0, 1}},
		10: {{9, 3}},
		11: {{0, 1}},
		12: {{0, 1}},
	},
	Data: esBoardData{},
}
