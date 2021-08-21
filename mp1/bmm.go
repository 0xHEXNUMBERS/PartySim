package mp1

type bmmBoardData struct {
	MagmaActive    bool
	MagmaTurnCount int
}

func bmmEruptVolcano(g *Game, player int) {
	bd := g.Board.Data.(bmmBoardData)
	if !bd.MagmaActive {
		bd.MagmaActive = true
		bd.MagmaTurnCount = 8
		g.Board.Data = bd
	}
}

func bmmCharacterEndTurn(g *Game, player int) {
	bd := g.Board.Data.(bmmBoardData)
	if bd.MagmaActive {
		bd.MagmaTurnCount--
		if bd.MagmaTurnCount == 0 {
			bd.MagmaActive = false
		}
		g.Board.Data = bd
	}
}

func bmmLandOnRegularSpace(g *Game, player int) {
	bd := g.Board.Data.(bmmBoardData)
	if bd.MagmaActive {
		g.Players[player].LastSpaceType = Red
	} else {
		g.Players[player].LastSpaceType = Blue
	}
}

func bmmReachFork(bowserPath, starPath ChainSpace) func(*Game, int, int) int {
	return func(g *Game, player, moves int) int {
		if g.Players[player].Coins >= 10 {
			g.ExtraEvent = bmmBranchPay{player, moves, bowserPath, starPath}
		} else {
			g.Players[player].CurrentSpace = bowserPath
		}
		return moves - 1
	}
}

func bmmFinalFork(g *Game, player, moves int) int {
	g.ExtraEvent = bmmBranchDecision{
		player, moves, ChainSpace{4, 0}, ChainSpace{5, 0},
	}
	return moves
}

func bmmVisitBowser(g *Game, player, moves int) int {
	if g.Players[player].Stars > 0 {
		g.ExtraEvent = bmmBowserRoulette{player, moves}
	} else {
		g.AwardCoins(player, -20, false)
	}
	return moves
}

var BMM = Board{
	Chains: &[]Chain{
		{ //After last fork to first fork
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: MinigameSpace},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Star},
			{Type: Mushroom},
			{Type: Happening, StoppingEvent: bmmEruptVolcano},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Bowser},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Start},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Star},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, PassingEvent: bmmReachFork(
				ChainSpace{1, 0}, ChainSpace{2, 2},
			)},
		},
		{ //Fork 1: Bowser Path to Fork 2
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: MinigameSpace},
			{Type: Red},
			{Type: Star},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Happening, StoppingEvent: bmmEruptVolcano},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, PassingEvent: bmmReachFork(
				ChainSpace{2, 0}, ChainSpace{3, 5},
			)},
		},
		{ //Fork 2: Bowser Path to Fork 3
			{Type: Mushroom},
			{Type: Red},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Star},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Happening, StoppingEvent: bmmEruptVolcano},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, PassingEvent: bmmReachFork(
				ChainSpace{3, 0}, ChainSpace{0, 1},
			)},
		},
		{ //Fork 3: BowserPath to Fork 4
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Bowser},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Star},
			{Type: MinigameSpace},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Happening, StoppingEvent: bmmEruptVolcano},
			{Type: Mushroom},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Star},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, PassingEvent: bmmFinalFork},
		},
		{ //Fork 4: Bowser Path
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Red},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, PassingEvent: bmmVisitBowser},
			{Type: Red},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Happening, StoppingEvent: bmmEruptVolcano},
		},
		{ //Fork 4: Star Path
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Boo},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
			{Type: Star},
			{Type: Invisible, StoppingEvent: bmmLandOnRegularSpace, HiddenBlock: true},
		},
	},
	Links: &map[int]*[]ChainSpace{
		4: {{0, 0}},
		5: {{0, 0}},
	},
	BowserCoins:           0,
	Data:                  bmmBoardData{},
	EndCharacterTurnEvent: bmmCharacterEndTurn,
}
