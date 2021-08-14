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
	expectedSpace := ChainSpace{5, 0}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}

	gBlueUpLeft.ExtraEvent.Handle(1, &gBlueUpLeft)
	gBlueUpLeft.ExtraEvent.Handle(
		ChainSpace{3, 0}, &gBlueUpLeft,
	)
	expectedSpace = ChainSpace{3, 0}
	gotSpace = gBlueUpLeft.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}

	gBlueUpRight.ExtraEvent.Handle(1, &gBlueUpRight)
	gBlueUpRight.ExtraEvent.Handle(
		ChainSpace{11, 0}, &gBlueUpRight,
	)
	expectedSpace = ChainSpace{11, 0}
	gotSpace = gBlueUpRight.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
}

func TestRBFork(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{3, 8}
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.ExtraEvent.Handle(1, &g)
	expectedSpace := ChainSpace{4, 0}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}

	gBlueUp.ExtraEvent.Handle(1, &gBlueUp)
	expectedSpace = ChainSpace{4, 4}
	gotSpace = gBlueUp.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
}

func TestBRFork1(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{5, 3}
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.ExtraEvent.Handle(1, &g)
	expectedSpace := ChainSpace{6, 10}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}

	gBlueUp.ExtraEvent.Handle(1, &gBlueUp)
	expectedSpace = ChainSpace{9, 0}
	gotSpace = gBlueUp.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
}

func TestBRFork2(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{6, 12}
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.ExtraEvent.Handle(1, &g)
	expectedSpace := ChainSpace{6, 0}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}

	gBlueUp.ExtraEvent.Handle(1, &gBlueUp)
	expectedSpace = ChainSpace{7, 0}
	gotSpace = gBlueUp.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
}

func TestBRFork3(t *testing.T) {
	g := *InitializeGame(LER, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{9, 8}
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.ExtraEvent.Handle(1, &g)
	expectedSpace := ChainSpace{10, 0}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}

	gBlueUp.ExtraEvent.Handle(1, &gBlueUp)
	expectedSpace = ChainSpace{0, 0}
	gotSpace = gBlueUp.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
}

func TestSwapGates(t *testing.T) {
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
