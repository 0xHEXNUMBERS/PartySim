package mp1

type mrcBoardData struct {
	IsBowser bool
}

func mrcSwapCastleDir(g *Game, player int) {
	bd := g.Board.Data.(mrcBoardData)
	bd.IsBowser = !bd.IsBowser
	g.Board.Data = bd
}

func mrcVisitCastle(g *Game, player int, moves int) int {
	bd := g.Board.Data.(mrcBoardData)
	if bd.IsBowser {
		g.AwardCoins(player, -40, false)
	} else {
		if g.Players[player].Coins >= 20 {
			g.Players[player].Stars++
			g.AwardCoins(player, -20, false)
		}
	}
	g.Players[player].CurrentSpace = ChainSpace{0, 0}
	mrcSwapCastleDir(g, player)
	return moves
}

var MRC = Board{
	Chains: &[]Chain{
		{ //Start to first fork
			{Type: Invisible}, //Temp space so players can walk on Start
			{Type: Start},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: Blue},
			{Type: Red},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
		},
		{ //First fork: left
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Chance},
			{Type: Blue},
			{Type: Blue},
			{Type: Red},
			{Type: Boo},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
		},
		{ //First fork: right to second fork
			{Type: Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: Red},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Red},
			{Type: Blue},
			{Type: Blue},
			{Type: Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: Blue},
			{Type: Blue},
		},
		{ //Second fork: right
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: Blue},
		},
		{ //Second fork: left to end
			{Type: Blue},
			{Type: Red},
			{Type: Blue},
			{Type: Blue},
			{Type: Happening, StoppingEvent: mrcSwapCastleDir},
			{Type: Blue},
			{Type: Chance},
			{Type: Bowser},
			{Type: Invisible, PassingEvent: mrcVisitCastle},
		},
	},
	Links: &map[int]*[]ChainSpace{
		0: {{1, 0}, {2, 0}},
		1: {{2, 2}},
		2: {{3, 0}, {4, 0}},
		3: {{4, 3}},
	},
	Data: mrcBoardData{},
}
