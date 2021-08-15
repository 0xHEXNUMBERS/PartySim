package mp1

import (
	"reflect"
	"testing"
)

func TestBMMMovement(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20})
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	g.ExtraEvent.Handle(3, &g)        //Move

	expectedSpace := ChainSpace{0, 15}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected Space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
	expectedCoins := 13
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Coins: %d, got: %d",
			expectedCoins, gotCoins)
	}
}

func TestBMMFork(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20})
	g.ExtraEvent.Handle(uint8(0), &g) //Star

	g.Players[0].CurrentSpace = ChainSpace{0, 17}
	g.ExtraEvent.Handle(1, &g)

	gIgnore := g
	gIgnore.ExtraEvent.Handle(false, &gIgnore) //Do not pay
	expectedSpace := ChainSpace{1, 0}
	gotSpace := gIgnore.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected Ignore Space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
	expectedCoins := 13
	gotCoins := gIgnore.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Ignore Coins: %d, got: %d",
			expectedCoins, gotCoins)
	}

	gPay := g
	gPay.ExtraEvent.Handle(true, &gPay)

	expectedCoins = 0
	gotCoins = gPay.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Pay Coins: %d, got: %d",
			expectedCoins, gotCoins)
	}
	expectedRes := []Response{ChainSpace{1, 0}, ChainSpace{2, 2}}
	gotRes := gPay.ExtraEvent.Responses()
	if !reflect.DeepEqual(expectedRes, gotRes) {
		t.Errorf("Expected Pay Responses: %#v, got: %#v",
			expectedRes, gotRes)
	}

	gBowser := gPay
	gBowser.ExtraEvent.Handle(ChainSpace{1, 0}, &gBowser)

	expectedSpace = ChainSpace{1, 0}
	gotSpace = gBowser.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected Bowser Space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
	expectedCoins = 3
	gotCoins = gBowser.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Bowser Coins: %d, got: %d",
			expectedCoins, gotCoins)
	}

	gStar := gPay
	gStar.ExtraEvent.Handle(ChainSpace{2, 2}, &gStar)

	expectedSpace = ChainSpace{2, 2}
	gotSpace = gStar.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected Star Space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
	expectedCoins = 3
	gotCoins = gStar.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Star Coins: %d, got: %d",
			expectedCoins, gotCoins)
	}
}

func TestBMMVolcano(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20})
	g.ExtraEvent.Handle(uint8(0), &g) //Star

	g.Players[0].CurrentSpace = ChainSpace{1, 5}
	g.ExtraEvent.Handle(1, &g) //Move P0 to Happening

	bd := g.Board.Data.(bmmBoardData)
	if !bd.MagmaActive {
		t.Errorf("Magma is not set")
	}

	g.ExtraEvent.Handle(1, &g) //Move P1 to Red

	expectedCoins := 7
	gotCoins := g.Players[1].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Coins: %d, got: %d",
			expectedCoins, gotCoins)
	}

	g.ExtraEvent.Handle(1, &g) //Move P2
	g.ExtraEvent.Handle(1, &g) //Move P3

	g.ExtraEvent.Handle(false, &g) //P0 is red
	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(3, &g) //P3 wins

	g.ExtraEvent.Handle(1, &g) //Move P0
	g.ExtraEvent.Handle(1, &g) //Move P1
	g.ExtraEvent.Handle(1, &g) //Move P2

	bd = g.Board.Data.(bmmBoardData)
	expectedTurnCount := 1
	gotTurnCount := bd.MagmaTurnCount
	if bd.MagmaTurnCount != 1 {
		t.Errorf("Expected Turn Count: %d, got: %d",
			expectedTurnCount, gotTurnCount)
	}

	g.ExtraEvent.Handle(1, &g) //Move P3

	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(3, &g) //P3 wins

	bd = g.Board.Data.(bmmBoardData)
	if bd.MagmaActive {
		t.Errorf("Magma has not been reset")
	}
}

func TestHiddenBlockOnInvisibleSpace(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20, EventsDice: true})
	g.ExtraEvent.Handle(uint8(0), &g) //Star

	g.ExtraEvent.Handle(NormalDiceBlock{0}, &g)
	g.ExtraEvent.Handle(1, &g)     //Move
	g.ExtraEvent.Handle(false, &g) //No hidden block here

	expectedCoins := 13
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected Coins: %d, got: %d",
			expectedCoins, gotCoins)
	}
}
