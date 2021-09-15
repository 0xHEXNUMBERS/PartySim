package board

import "github.com/0xhexnumbers/partysim/mp1"

//bmmBoardData holds all of the board specific data related to BMM.
type bmmBoardData struct {
	MagmaActive    bool
	MagmaTurnCount int
}

//bmmEruptVolcano turns all *blue* spaces to *red* spaces for 2 full game
//turns.
func bmmEruptVolcano(g *mp1.Game, player int) {
	bd := g.Board.Data.(bmmBoardData)
	if !bd.MagmaActive {
		bd.MagmaActive = true
		bd.MagmaTurnCount = 8
		g.Board.Data = bd
	}
}

//bmmCharacterEndTurn decrements the turn counter at the end of each
//character's turn, swapping all *red* spaces back to *blue* spaces.
func bmmCharacterEndTurn(g *mp1.Game, player int) {
	bd := g.Board.Data.(bmmBoardData)
	if bd.MagmaActive {
		bd.MagmaTurnCount--
		if bd.MagmaTurnCount == 0 {
			bd.MagmaActive = false
		}
		g.Board.Data = bd
	}
}

//bmmLandOnRegularSpace sets the player's LastSpaceType to Blue/Red
//depending if the volcano is active or not.
func bmmLandOnRegularSpace(g *mp1.Game, player int) {
	bd := g.Board.Data.(bmmBoardData)
	if bd.MagmaActive {
		g.Players[player].LastSpaceType = mp1.Red
	} else {
		g.Players[player].LastSpaceType = mp1.Blue
	}
}

//bmmReachFork sets the next event to a custom branch event if the player
//has >=10 coins. Otherwise, they continue down the bowser path.
func bmmReachFork(bowserPath, starPath mp1.ChainSpace) func(*mp1.Game, int, int) int {
	return func(g *mp1.Game, player, moves int) int {
		if g.Players[player].Coins >= 10 {
			g.NextEvent = BMMBranchPay{player, moves, bowserPath, starPath}
		} else {
			g.Players[player].CurrentSpace = bowserPath
		}
		return moves - 1
	}
}

//bmmFinalFork sets the next event to the custom branch event.
func bmmFinalFork(g *mp1.Game, player, moves int) int {
	g.NextEvent = BMMBranchDecision{
		player, moves, mp1.NewChainSpace(4, 0), mp1.NewChainSpace(5, 0),
	}
	return moves
}

//bmmVisitBowser rolls a roulette if the player has one. Otherwise, bowser
//steals 20 coins.
func bmmVisitBowser(g *mp1.Game, player, moves int) int {
	if g.Players[player].Stars > 0 {
		g.NextEvent = BMMBowserRoulette{player, moves}
	} else {
		g.AwardCoins(player, -20, false)
	}
	return moves
}

//BMM holds the data for Bowser's Magma Mountain.
var BMM = mp1.Board{
	Chains: &[]mp1.Chain{
		{ //After last fork to first fork
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Star},
			{Type: mp1.Mushroom},
			{Type: mp1.Happening, StoppingEvent: bmmEruptVolcano},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Bowser},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Start},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Star},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, PassingEvent: bmmReachFork(
				mp1.NewChainSpace(1, 0), mp1.NewChainSpace(2, 2),
			)},
		},
		{ //Fork 1: Bowser Path to Fork 2
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Red},
			{Type: mp1.Star},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Happening, StoppingEvent: bmmEruptVolcano},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, PassingEvent: bmmReachFork(
				mp1.NewChainSpace(2, 0), mp1.NewChainSpace(3, 5),
			)},
		},
		{ //Fork 2: Bowser Path to Fork 3
			{Type: mp1.Mushroom},
			{Type: mp1.Red},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Star},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Happening, StoppingEvent: bmmEruptVolcano},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, PassingEvent: bmmReachFork(
				mp1.NewChainSpace(3, 0), mp1.NewChainSpace(0, 1),
			)},
		},
		{ //Fork 3: BowserPath to Fork 4
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Bowser},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Star},
			{Type: mp1.MinigameSpace},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Happening, StoppingEvent: bmmEruptVolcano},
			{Type: mp1.Mushroom},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Star},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, PassingEvent: bmmFinalFork},
		},
		{ //Fork 4: Bowser Path
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Red},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, PassingEvent: bmmVisitBowser},
			{Type: mp1.Red},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Happening, StoppingEvent: bmmEruptVolcano},
		},
		{ //Fork 4: Star Path
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Boo},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: mp1.Star},
			{Type: mp1.Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
		},
	},
	Links: &map[int]*[]mp1.ChainSpace{
		4: {mp1.NewChainSpace(0, 0)},
		5: {mp1.NewChainSpace(0, 0)},
	},
	BowserCoins:           0,
	Data:                  bmmBoardData{},
	EndCharacterTurnEvent: bmmCharacterEndTurn,
}
