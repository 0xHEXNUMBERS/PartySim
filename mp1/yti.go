package mp1

//ytiBoardData holds all of the board specific data related to YTI.
type ytiBoardData struct {
	Thwomps         [2]int
	AcceptThwompPos [2]ChainSpace
	RejectThwompPos [2]ChainSpace
	StarPosition    ChainSpace
}

//ytiCheckThwomp checks to see if the current player can pass this thwomp.
func ytiCheckThwomp(thwomp int) func(*Game, int, int) int {
	return func(g *Game, player, moves int) int {
		bd := g.Board.Data.(ytiBoardData)
		if g.Players[player].Coins >= bd.Thwomps[thwomp] {
			g.NextEvent = ytiThwompBranchEvent{
				player,
				moves,
				thwomp,
			}
		} else {
			pos := bd.RejectThwompPos[thwomp]
			g.Players[player].CurrentSpace = pos
		}
		return moves - 1
	}
}

var ytiLeftIslandStar = ChainSpace{0, 19}
var ytiRightIslandStar = ChainSpace{1, 18}

//ytiSwapStarPosition will swap the current star position.
func ytiSwapStarPosition(g *Game, player int) {
	bd := g.Board.Data.(ytiBoardData)
	if bd.StarPosition == ytiLeftIslandStar {
		bd.StarPosition = ytiRightIslandStar
	} else {
		bd.StarPosition = ytiLeftIslandStar
	}
	g.Board.Data = bd
}

//ytiGainStar will increment the player's star count if they have >=20
//coins.
func ytiGainStar(g *Game, player, moves int) int {
	bd := g.Board.Data.(ytiBoardData)
	if bd.StarPosition == g.Players[player].CurrentSpace {
		if g.Players[player].Coins >= 20 {
			g.AwardCoins(player, -20, false)
			g.Players[player].Stars++
			ytiSwapStarPosition(g, 0)
		}
	} else { //Star at other island
		g.AwardCoins(player, -30, false)
	}
	return moves
}

//YTI holds the data for Yoshi's Tropical Island.
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
			{Type: Invisible, PassingEvent: ytiGainStar},
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
			{Type: Invisible, PassingEvent: ytiGainStar},
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
	},
	Links: nil,
	Data: ytiBoardData{
		[2]int{1, 1},
		[2]ChainSpace{{1, 6}, {0, 7}},
		[2]ChainSpace{{0, 0}, {1, 0}},
		ytiLeftIslandStar,
	},
}
