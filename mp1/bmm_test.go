package mp1

import (
	"testing"
)

func TestBMMMovement(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20})
	g.NextEvent.Handle(uint8(0), &g) //Star
	g.NextEvent.Handle(3, &g)        //Move

	SpaceIs(ChainSpace{0, 15}, 0, g, "", t)
	CoinsIs(13, 0, g, "", t)
}

func TestBMMFork(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20})
	g.NextEvent.Handle(uint8(0), &g) //Star

	g.Players[0].CurrentSpace = ChainSpace{0, 17}
	g.NextEvent.Handle(1, &g)

	gIgnore := g
	gIgnore.NextEvent.Handle(false, &gIgnore) //Do not pay
	SpaceIs(ChainSpace{1, 0}, 0, gIgnore, "Ignore", t)
	CoinsIs(13, 0, gIgnore, "Ignore", t)

	gPay := g
	gPay.NextEvent.Handle(true, &gPay)

	CoinsIs(0, 0, gPay, "Pay", t)
	expectedRes := []Response{ChainSpace{1, 0}, ChainSpace{2, 2}}
	ResIs(expectedRes, gPay, "Pay", t)

	gBowser := gPay
	gBowser.NextEvent.Handle(ChainSpace{1, 0}, &gBowser)

	SpaceIs(ChainSpace{1, 0}, 0, gBowser, "Bowser", t)
	CoinsIs(3, 0, gBowser, "Bowser", t)

	gStar := gPay
	gStar.NextEvent.Handle(ChainSpace{2, 2}, &gStar)

	SpaceIs(ChainSpace{2, 2}, 0, gStar, "Star", t)
	CoinsIs(3, 0, gStar, "Star", t)
}

func TestBMMVolcano(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20})
	g.NextEvent.Handle(uint8(0), &g) //Star

	g.Players[0].CurrentSpace = ChainSpace{1, 5}
	g.NextEvent.Handle(1, &g) //Move P0 to Happening

	bd := g.Board.Data.(bmmBoardData)
	if !bd.MagmaActive {
		t.Errorf("Magma is not set")
	}

	g.NextEvent.Handle(1, &g) //Move P1 to Red

	CoinsIs(7, 1, g, "", t)

	g.NextEvent.Handle(1, &g) //Move P2
	g.NextEvent.Handle(1, &g) //Move P3

	g.NextEvent.Handle(false, &g) //P0 is red
	g.NextEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.NextEvent.Handle(3, &g) //P3 wins

	g.NextEvent.Handle(1, &g) //Move P0
	g.NextEvent.Handle(1, &g) //Move P1
	g.NextEvent.Handle(1, &g) //Move P2

	bd = g.Board.Data.(bmmBoardData)
	IntIs(1, bd.MagmaTurnCount, "Turn Count", t)

	g.NextEvent.Handle(1, &g) //Move P3

	g.NextEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.NextEvent.Handle(3, &g) //P3 wins

	bd = g.Board.Data.(bmmBoardData)
	if bd.MagmaActive {
		t.Errorf("Magma has not been reset")
	}
}

func TestHiddenBlockOnInvisibleSpace(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20, EventsDice: true})
	g.NextEvent.Handle(uint8(0), &g) //Star

	g.NextEvent.Handle(NormalDiceBlock{0}, &g)
	g.NextEvent.Handle(1, &g)     //Move
	g.NextEvent.Handle(false, &g) //No hidden block here

	CoinsIs(13, 0, g, "Coins", t)
}
