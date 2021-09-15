package board

import "github.com/0xhexnumbers/partysim/mp1"

//esBoardData holds all of the board specific data related to ES.
type esBoardData struct {
	StarTaken [7]bool
	Gate      int //0 if unknown
	Gate2or3  bool
}

var esEntrance1 = mp1.NewChainSpace(0, 0)
var esStartingSpace = mp1.NewChainSpace(0, 1)
var esEntrance4 = mp1.NewChainSpace(7, 0)
var esEntrance5 = mp1.NewChainSpace(8, 0)
var esEntrance6 = mp1.NewChainSpace(9, 0)
var esEntrance7 = mp1.NewChainSpace(10, 0)
var esEntrance8 = mp1.NewChainSpace(11, 0)
var esEntrance9 = mp1.NewChainSpace(12, 0)

//esVisitBowser steals a star if the player has one, or steals 20 coins
//otherwise. The next event is set to change the gate the board will play
//under.
func esVisitBowser(g *mp1.Game, player, moves int) int {
	if g.Players[player].Stars > 0 {
		g.Players[player].Stars--
	} else {
		g.AwardCoins(player, -20, false)
	}
	bd := g.Board.Data.(esBoardData)
	g.NextEvent = ESChangeGates{player, moves, bd.Gate}
	return moves
}

//esSendToStart sends each player to the starting space, and the landing
//player gains coins from a blue space.
func esSendToStart(g *mp1.Game, player int) {
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

//esWarpC handles Warp C.
func esWarpC(g *mp1.Game, player, moves int) int {
	bd := g.Board.Data.(esBoardData)
	if bd.Gate2or3 {
		g.Players[player].CurrentSpace = esStartingSpace
	} else {
		switch bd.Gate {
		case 0:
			g.NextEvent = ESWarpCDest{player, moves}
		case 1:
			g.Players[player].CurrentSpace = esEntrance7
		default:
			g.Players[player].CurrentSpace = esEntrance1
		}
	}
	return moves
}

//esWarpSpace handles warp spaces D-K.
func esWarpSpace(dest1, dest2, dest3 mp1.ChainSpace) func(*mp1.Game, int, int) int {
	return func(g *mp1.Game, player, moves int) int {
		bd := g.Board.Data.(esBoardData)
		switch bd.Gate {
		case 0:
			g.NextEvent = ESWarpDest{
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

//esBranchWithWarp handles warps that are on a chain by themselves.
func esBranchWithWarp(dest1, dest2, dest3 mp1.ChainSpace) func(*mp1.Game, int, int) int {
	return func(g *mp1.Game, player, moves int) int {
		g.NextEvent = ESBranchEvent{
			player,
			moves,
			dest1,
			dest2,
			dest3,
		}
		return moves
	}
}

//esAllStarsCollected returns true if all stars on the board have been
//collected.
func esAllStarsCollected(e esBoardData) bool {
	for _, star := range e.StarTaken {
		if !star {
			return false
		}
	}
	return true
}

//esPassStarSpace sets the next event to decide if the player wants to play
//the baby bowser minigame if they have >=20 coins and the star space is
//available. If the player does not have 20 coins, then they pass by the
//space. If the star has been collected, then the space is assumed to be
//landable.
func esPassStarSpace(i int) func(*mp1.Game, int, int) int {
	return func(g *mp1.Game, player, moves int) int {
		bd := g.Board.Data.(esBoardData)
		if !bd.StarTaken[i] {
			if g.Players[player].Coins >= 20 {
				g.NextEvent = ESVisitBabyBowser{
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

//esLandOnChanceTime sets the player's LastSaceType to Chance.
func esLandOnChanceTime(g *mp1.Game, player int) {
	g.Players[player].LastSpaceType = mp1.Chance
}

//ES holds the data for Eternal Star.
var ES = mp1.Board{
	Chains: &[]mp1.Chain{
		{ //0: Entrance 1
			{Type: mp1.Invisible}, //Warp 1 Entrance
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
		},
		{ //1: Entrance 1 Fork: Right Exit Through warp A to Entrance 2 Fork
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //2: Entrance 1 Fork: Left Exit Through warp B to Entrance 3 Fork
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //3: Entrance 2 Fork: Left Exit to warp C
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Bowser},
			{Type: mp1.Invisible, PassingEvent: esWarpC},
		},
		{ //4: Entrance 2 Fork: Right Exit to warp D
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{
				Type:          mp1.Invisible,
				PassingEvent:  esPassStarSpace(0),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: esSendToStart},
			{
				Type: mp1.Invisible,
				PassingEvent: esWarpSpace(
					esEntrance1,
					esEntrance7,
					esEntrance6,
				),
			},
		},
		{ //5: Entrance 3 Fork: Right exit to warp E
			{Type: mp1.Red},
			{
				Type:          mp1.Invisible,
				PassingEvent:  esPassStarSpace(1),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{
				Type: mp1.Invisible,
				PassingEvent: esWarpSpace(
					esEntrance6,
					esEntrance4,
					esEntrance1,
				),
			},
		},
		{ //6: Entrance 3 Fork: Left exit to warp F
			{Type: mp1.Blue},
			{Type: mp1.Happening, StoppingEvent: esSendToStart},
			{
				Type: mp1.Invisible,
				PassingEvent: esWarpSpace(
					esEntrance1,
					esEntrance6,
					esEntrance7,
				),
			},
		},
		{ //7: Entrance 4 to Warp H with Warp G branch
			{Type: mp1.Invisible}, //Movement space
			{Type: mp1.Blue},
			{
				Type:          mp1.Invisible,
				PassingEvent:  esPassStarSpace(2),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{ //Warp G
				Type: mp1.Invisible,
				PassingEvent: esBranchWithWarp(
					esEntrance9,
					esEntrance9,
					esEntrance8,
				),
			},
			{Type: mp1.Blue},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{
				Type:          mp1.Invisible,
				PassingEvent:  esPassStarSpace(3),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: mp1.Blue},
			{
				Type: mp1.Invisible,
				PassingEvent: esWarpSpace(
					esEntrance5,
					esEntrance8,
					esEntrance5,
				),
			},
		},
		{ //8: Entrance 5 to Warp I (always goes to start
			{Type: mp1.Invisible}, //Movement
			{Type: mp1.Red},
			{Type: mp1.Happening, StoppingEvent: esSendToStart},
			{
				Type:          mp1.Invisible,
				PassingEvent:  esPassStarSpace(4),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Bowser},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //9: Entrance 6 to warp K with Warp J Branch
			{Type: mp1.Invisible}, //Tmp space for warp
			{Type: mp1.Blue},
			{
				Type:          mp1.Invisible,
				PassingEvent:  esPassStarSpace(5),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Red},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{
				Type:          mp1.Invisible,
				PassingEvent:  esPassStarSpace(6),
				StoppingEvent: esLandOnChanceTime,
			},
			{Type: mp1.Happening, StoppingEvent: esSendToStart},
			{ //Warp J
				Type: mp1.Invisible,
				PassingEvent: esBranchWithWarp(
					esEntrance9,
					esEntrance9,
					esEntrance8,
				),
			},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Mushroom},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{Type: mp1.Blue},
			{
				Type: mp1.Invisible,
				PassingEvent: esWarpSpace(
					esEntrance8,
					esEntrance1,
					esEntrance4,
				),
			},
		},
		{ //10: Entrance 7 to Entrance 6 convergence
			{Type: mp1.Invisible}, //Tmp space for movement
			{Type: mp1.Blue},
			{Type: mp1.Blue},
		},
		{ //11: Entrance 8
			{Type: mp1.Invisible}, //Tmp space for movement
			{Type: mp1.Boo},
		},
		{ //12: Entrance 9
			{Type: mp1.Invisible},                              //Tmp space for movement
			{Type: mp1.Invisible, PassingEvent: esVisitBowser}, //Bowser
		},
	},
	Links: &map[int]*[]mp1.ChainSpace{
		0:  {mp1.NewChainSpace(1, 0), mp1.NewChainSpace(2, 0)},
		1:  {mp1.NewChainSpace(3, 0), mp1.NewChainSpace(4, 0)},
		2:  {mp1.NewChainSpace(5, 0), mp1.NewChainSpace(6, 0)},
		8:  {mp1.NewChainSpace(0, 1)},
		10: {mp1.NewChainSpace(9, 3)},
		11: {mp1.NewChainSpace(0, 1)},
		12: {mp1.NewChainSpace(0, 1)},
	},
	Data: esBoardData{},
}
