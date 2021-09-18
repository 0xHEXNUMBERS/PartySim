package board

import "github.com/0xhexnumbers/partysim/mp1"

//ytiBoardData holds all of the board specific data related to YTI.
type ytiBoardData struct {
	Thwomps         [2]int
	AcceptThwompPos [2]mp1.ChainSpace
	RejectThwompPos [2]mp1.ChainSpace
	StarPosition    mp1.ChainSpace
}

//ytiCheckThwomp checks to see if the current player can pass this thwomp.
func ytiCheckThwomp(thwomp int) func(*mp1.Game, int, int) int {
	return func(g *mp1.Game, player, moves int) int {
		bd := g.Board.Data.(ytiBoardData)
		if g.Players[player].Coins >= bd.Thwomps[thwomp] {
			g.NextEvent = YTIThwompBranchEvent{
				Player: player,
				Moves:  moves,
				Thwomp: thwomp,
			}
		} else {
			pos := bd.RejectThwompPos[thwomp]
			g.Players[player].CurrentSpace = pos
		}
		return moves - 1
	}
}

var ytiLeftIslandStar = mp1.NewChainSpace(0, 19)
var ytiRightIslandStar = mp1.NewChainSpace(1, 18)

//ytiSwapStarPosition will swap the current star position.
func ytiSwapStarPosition(g *mp1.Game, player int) {
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
func ytiGainStar(g *mp1.Game, player, moves int) int {
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
var YTI = mp1.Board{
	Chains: &[]mp1.Chain{
		{ //Left island
			{Type: mp1.Blue}, //Branch #1 Dir A
			{Type: mp1.Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue}, //Branch #2 Dir B, links from before
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Chance},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: mp1.Invisible, PassingEvent: ytiGainStar},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Start},
			{Type: mp1.Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: ytiCheckThwomp(0)},
		},
		{ //Right island part 1
			{Type: mp1.Blue}, //Branch #2 Dir A
			{Type: mp1.Chance},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: mp1.Blue},
			{Type: mp1.Blue}, //Branch #1 Dir B
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: ytiGainStar},
			{Type: mp1.Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Boo},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Happening, StoppingEvent: ytiSwapStarPosition},
			{Type: mp1.Blue},
			{Type: mp1.Invisible, PassingEvent: ytiCheckThwomp(1)},
		},
	},
	Links: nil,
	Data: ytiBoardData{
		[2]int{1, 1},
		[2]mp1.ChainSpace{mp1.NewChainSpace(1, 6), mp1.NewChainSpace(0, 7)},
		[2]mp1.ChainSpace{mp1.NewChainSpace(0, 0), mp1.NewChainSpace(1, 0)},
		ytiLeftIslandStar,
	},
}
