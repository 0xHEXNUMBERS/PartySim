package mp1

import "testing"

func TestWhompPayment(t *testing.T) {
	g := *InitializeGame(DKJA, GameConfig{MaxTurns: 20})
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	g.ExtraEvent.Handle(5, &g)        //Move

	gPay := g
	gPay.ExtraEvent.Handle(true, &gPay) //Pay Whomp

	expectedPos := ChainSpace{4, 0}
	gotPos := gPay.Players[0].CurrentSpace
	if expectedPos != gotPos {
		t.Errorf("Expected Pay pos: %#v, got: %#v",
			expectedPos, gotPos)
	}
	expectedCoins := 3
	gotCoins := gPay.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Pay Coins: %#v, got: %#v",
			expectedCoins, gotCoins)
	}

	gSkip := g
	gSkip.ExtraEvent.Handle(false, &gSkip) //Ignore Whomp

	expectedPos = ChainSpace{1, 0}
	gotPos = gSkip.Players[0].CurrentSpace
	if expectedPos != gotPos {
		t.Errorf("Expected Skip pos: %#v, got: %#v",
			expectedPos, gotPos)
	}
	expectedCoins = 13
	gotCoins = gSkip.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Skip Coins: %#v, got: %#v",
			expectedCoins, gotCoins)
	}

	gIgnore := gSkip
	gIgnore.Players[1].Coins = 0
	gIgnore.ExtraEvent.Handle(5, &gIgnore) //Move

	expectedPos = ChainSpace{1, 0}
	gotPos = gIgnore.Players[1].CurrentSpace
	if expectedPos != gotPos {
		t.Errorf("Expected Ignore pos: %#v, got: %#v",
			expectedPos, gotPos)
	}
	expectedCoins = 3
	gotCoins = gIgnore.Players[1].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Ignore Coins: %#v, got: %#v",
			expectedCoins, gotCoins)
	}
}

func TestCoinBlockade(t *testing.T) {
	g := *InitializeGame(DKJA, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{1, 3}
	g.Players[0].Coins = 20
	g.ExtraEvent.Handle(uint8(0), &g) //Star

	gPass := g
	gPass.ExtraEvent.Handle(1, &gPass)    //Move
	gPass.ExtraEvent.Handle(true, &gPass) //Pass

	expectedPos := ChainSpace{2, 0}
	gotPos := gPass.Players[0].CurrentSpace
	if expectedPos != gotPos {
		t.Errorf("Expected Pass pos: %#v, got: %#v",
			expectedPos, gotPos)
	}
	expectedCoins := 23
	gotCoins := gPass.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Pass Coins: %#v, got: %#v",
			expectedCoins, gotCoins)
	}

	gSkip := g
	gSkip.ExtraEvent.Handle(1, &gSkip)     //Move
	gSkip.ExtraEvent.Handle(false, &gSkip) //Skip

	expectedPos = ChainSpace{3, 0}
	gotPos = gSkip.Players[0].CurrentSpace
	if expectedPos != gotPos {
		t.Errorf("Expected Skip pos: %#v, got: %#v",
			expectedPos, gotPos)
	}
	expectedCoins = 23
	gotCoins = gSkip.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Skip Coins: %#v, got: %#v",
			expectedCoins, gotCoins)
	}

	gIgnore := g
	gIgnore.Players[0].Coins = 0
	gIgnore.ExtraEvent.Handle(1, &gIgnore) //Move

	expectedPos = ChainSpace{3, 0}
	gotPos = gIgnore.Players[0].CurrentSpace
	if expectedPos != gotPos {
		t.Errorf("Expected Ignore pos: %#v, got: %#v",
			expectedPos, gotPos)
	}
	expectedCoins = 3
	gotCoins = gIgnore.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Ignore Coins: %#v, got: %#v",
			expectedCoins, gotCoins)
	}
}

func TestBoulder(t *testing.T) {
	g := *InitializeGame(DKJA, GameConfig{MaxTurns: 20})
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	g.Players[0].CurrentSpace = ChainSpace{5, 1}
	g.Players[1].CurrentSpace = ChainSpace{5, 0}
	g.Players[2].CurrentSpace = ChainSpace{7, 0}
	g.Players[3].CurrentSpace = ChainSpace{5, 5}

	g.ExtraEvent.Handle(2, &g)
	playerPos := ChainSpace{0, 16}
	gotPos := g.Players[0].CurrentSpace
	if playerPos != gotPos {
		t.Errorf("Expected Player0 pos: %#v, got: %#v",
			playerPos, gotPos)
	}
	playerPos = ChainSpace{5, 0}
	gotPos = g.Players[1].CurrentSpace
	if playerPos != gotPos {
		t.Errorf("Expected Player1 pos: %#v, got: %#v",
			playerPos, gotPos)
	}
	playerPos = ChainSpace{0, 16}
	gotPos = g.Players[2].CurrentSpace
	if playerPos != gotPos {
		t.Errorf("Expected Player2 pos: %#v, got: %#v",
			playerPos, gotPos)
	}
	playerPos = ChainSpace{0, 16}
	gotPos = g.Players[3].CurrentSpace
	if playerPos != gotPos {
		t.Errorf("Expected Player3 pos: %#v, got: %#v",
			playerPos, gotPos)
	}
}
