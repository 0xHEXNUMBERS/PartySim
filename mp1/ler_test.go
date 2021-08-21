package mp1

import "testing"

func TestRBRFork(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{1, 4}
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	gBlueUpLeft := g
	gBlueUpLeft.Board.Data = lerBoardData{true}
	gBlueUpRight := gBlueUpLeft

	g.ExtraEvent.Handle(1, &g)
	SpaceIs(ChainSpace{5, 0}, 0, g, "", t)

	gBlueUpLeft.ExtraEvent.Handle(1, &gBlueUpLeft)
	gBlueUpLeft.ExtraEvent.Handle(
		ChainSpace{3, 0}, &gBlueUpLeft,
	)
	SpaceIs(ChainSpace{3, 0}, 0, gBlueUpLeft, "BlueUpLeft", t)

	gBlueUpRight.ExtraEvent.Handle(1, &gBlueUpRight)
	gBlueUpRight.ExtraEvent.Handle(
		ChainSpace{11, 0}, &gBlueUpRight,
	)
	SpaceIs(ChainSpace{11, 0}, 0, gBlueUpRight, "BlueUpRight", t)
}

func TestRBFork(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{3, 8}
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.ExtraEvent.Handle(1, &g)
	SpaceIs(ChainSpace{4, 0}, 0, g, "", t)

	gBlueUp.ExtraEvent.Handle(1, &gBlueUp)
	SpaceIs(ChainSpace{4, 4}, 0, gBlueUp, "BlueUp", t)
}

func TestBRFork1(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{5, 3}
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.ExtraEvent.Handle(1, &g)
	SpaceIs(ChainSpace{6, 10}, 0, g, "", t)

	gBlueUp.ExtraEvent.Handle(1, &gBlueUp)
	SpaceIs(ChainSpace{9, 0}, 0, gBlueUp, "BlueUp", t)
}

func TestBRFork2(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{6, 12}
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.ExtraEvent.Handle(1, &g)
	SpaceIs(ChainSpace{6, 0}, 0, g, "", t)

	gBlueUp.ExtraEvent.Handle(1, &gBlueUp)
	SpaceIs(ChainSpace{7, 0}, 0, gBlueUp, "BlueUp", t)
}

func TestBRFork3(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{9, 8}
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.ExtraEvent.Handle(1, &g)
	SpaceIs(ChainSpace{10, 0}, 0, g, "", t)

	gBlueUp.ExtraEvent.Handle(1, &gBlueUp)
	SpaceIs(ChainSpace{0, 0}, 0, gBlueUp, "BlueUp", t)
}

func TestSwapGatesViaHappening(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{3, 7}
	g.Players[0].Coins = 20
	g.ExtraEvent.Handle(uint8(0), &g) //Star

	g.ExtraEvent.Handle(1, &g) //Move
	bd := g.Board.Data.(lerBoardData)
	if !bd.BlueUp {
		t.Errorf("Gates did not swap via Happening")
	}
}

func TestSwapGatesTwice(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{11, 5}
	g.Players[0].Coins = 20
	g.ExtraEvent.Handle(uint8(0), &g) //Star

	g.ExtraEvent.Handle(1, &g)    //Move
	g.ExtraEvent.Handle(true, &g) //Swap Gates
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		t.Errorf("Gates did not swap twice")
	}
}

func TestSwapGatesInsufficientCoins(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{2, 3}
	g.Players[0].Coins = 0
	g.ExtraEvent.Handle(uint8(0), &g) //Star

	g.ExtraEvent.Handle(1, &g)

	SpaceIs(ChainSpace{2, 5}, 0, g, "", t)
}

func TestNormalBranch(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{0, 10}
	g.ExtraEvent.Handle(uint8(0), &g)

	g.ExtraEvent.Handle(1, &g)                //Move
	g.ExtraEvent.Handle(ChainSpace{1, 0}, &g) //Branch

	SpaceIs(ChainSpace{1, 0}, 0, g, "", t)
}
