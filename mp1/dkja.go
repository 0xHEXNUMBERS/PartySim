package mp1

type dkjaBoardData struct {
	WhompPos                 [3]bool
	WhompMainDestination     [3]ChainSpace
	WhompOffshootDestination [3]ChainSpace
	CoinAcceptDestination    [2]ChainSpace
	CoinRejectDestination    [2]ChainSpace
}

func dkjaGetWhompDestination(g *Game, whomp int) ChainSpace {
	data := g.Board.Data.(dkjaBoardData)
	var pos ChainSpace
	if data.WhompPos[whomp] {
		pos = data.WhompOffshootDestination[whomp]
	} else {
		pos = data.WhompMainDestination[whomp]
	}
	return pos
}

func dkjaCanPassWhomp(whomp int) func(*Game, int, int) {
	return func(g *Game, player, moves int) {
		if g.Players[player].Coins >= 10 {
			g.ExtraEvent = dkjaWhompEvent{
				player, moves, whomp,
			}
		} else {
			pos := dkjaGetWhompDestination(g, whomp)
			g.Players[player].CurrentSpace = pos
		}
	}
}

func dkjaCanPassCoinBlockade(blockade int) func(*Game, int, int) {
	return func(g *Game, player, moves int) {
		if g.Players[player].Coins >= 20 {
			g.ExtraEvent = dkjaCoinBranchEvent{
				player, moves, blockade,
			}
		} else {
			data := g.Board.Data.(dkjaBoardData)
			g.Players[player].CurrentSpace = data.CoinRejectDestination[blockade]
		}
	}
}

func dkjaBoulder(g *Game, player int) {
	for i := 0; i < 4; i++ {
		pos := g.Players[i].CurrentSpace
		if pos.Chain == 7 || (pos.Chain == 5 && pos.Space != 0) {
			g.Players[i].CurrentSpace = ChainSpace{0, 16}
		}
	}
}

var DKJA = Board{
	Chains: &[]Chain{
		{ //Last Offshoot to first thwomp fork
			{Type: Blue},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
			{Type: Star},
			{Type: Happening, StoppingEvent: dkjaBoulder},
			{Type: Happening, StoppingEvent: dkjaBoulder},
			{Type: Blue},
			{Type: Boo},
			{Type: Blue},
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Start},
			{Type: Blue},
			{Type: Red},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: dkjaCanPassWhomp(0)},
		},
		{ //First offshoot to coin blockade
			{Type: Blue},
			{Type: Blue},
			{Type: Red},
			{Type: Blue},
			{Type: Invisible, PassingEvent: dkjaCanPassCoinBlockade(0)},
		},
		{ //Through first coin blockade
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Star},
			{Type: Red},
		},
		{ //Around first coin blockade
			{Type: Blue},
			{Type: Blue},
			{Type: Happening, StoppingEvent: dkjaBoulder},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
		},
		{ //First main pathway
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: dkjaCanPassWhomp(1)},
		},
		{ //Second main pathway
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Blue},
			{Type: Happening, StoppingEvent: dkjaBoulder},
			{Type: Happening, StoppingEvent: dkjaBoulder},
			{Type: Blue},
			{Type: Star},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Happening, StoppingEvent: dkjaBoulder},
			{Type: Boo},
			{Type: Invisible, PassingEvent: dkjaCanPassCoinBlockade(1)},
		},
		{ //Second offshoot pathway
			{Type: Blue},
			{Type: Blue},
			{Type: Red},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: Bowser},
			{Type: Blue},
			{Type: Blue},
		},
		{ //Through second coin blockade
			{Type: MinigameSpace},
			{Type: Red},
			{Type: Bowser},
			{Type: Red},
			{Type: Blue},
			{Type: Happening, StoppingEvent: dkjaBoulder},
			{Type: Happening, StoppingEvent: dkjaBoulder},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
		},
		{ //Around second coin blockade
			{Type: Blue},
			{Type: Mushroom},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: Invisible, PassingEvent: dkjaCanPassWhomp(2)},
		},
		{ //Third Main Pathway
			{Type: Red},
			{Type: Blue},
			{Type: Star},
			{Type: Blue},
			{Type: MinigameSpace},
			{Type: Happening, StoppingEvent: dkjaBoulder},
			{Type: Blue},
		},
	},
	Links: &map[int]*[]ChainSpace{
		2: {{0, 18}},
		3: {{0, 16}},
		6: {{5, 9}},
		7: {{0, 15}},
		9: {{0, 13}},
	},
	Data: dkjaBoardData{
		WhompPos: [3]bool{false, false, false},
		WhompMainDestination: [3]ChainSpace{
			{4, 0},
			{5, 0},
			{9, 0},
		},
		WhompOffshootDestination: [3]ChainSpace{
			{1, 0},
			{6, 0},
			{0, 0},
		},
		CoinAcceptDestination: [2]ChainSpace{
			{2, 0},
			{7, 0},
		},
		CoinRejectDestination: [2]ChainSpace{
			{3, 0},
			{8, 0},
		},
	},
}
