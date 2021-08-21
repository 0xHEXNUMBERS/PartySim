package mp1

import (
	"testing"
)

func TestBMMMovement(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20})
	g.ExtraEvent.Handle(uint8(0), &g) //Star
	g.ExtraEvent.Handle(3, &g)        //Move

	SpaceIs(ChainSpace{0, 15}, 0, g, "", t)
	CoinsIs(13, 0, g, "", t)
}

func TestBMMFork(t *testing.T) {
	g := *InitializeGame(BMM, GameConfig{MaxTurns: 20})
	g.ExtraEvent.Handle(uint8(0), &g) //Star

	g.Players[0].CurrentSpace = ChainSpace{0, 17}
	g.ExtraEvent.Handle(1, &g)

	gIgnore := g
	gIgnore.ExtraEvent.Handle(false, &gIgnore) //Do not pay
	SpaceIs(ChainSpace{1, 0}, 0, gIgnore, "Ignore", t)
	CoinsIs(13, 0, gIgnore, "Ignore", t)

	gPay := g
	gPay.ExtraEvent.Handle(true, &gPay)

	CoinsIs(0, 0, gPay, "Pay", t)
	expectedRes := []Response{ChainSpace{1, 0}, ChainSpace{2, 2}}
	ResIs(expectedRes, gPay, "Pay", t)

	gBowser := gPay
	gBowser.ExtraEvent.Handle(ChainSpace{1, 0}, &gBowser)

	SpaceIs(ChainSpace{1, 0}, 0, gBowser, "Bowser", t)
	CoinsIs(3, 0, gBowser, "Bowser", t)

	gStar := gPay
	gStar.ExtraEvent.Handle(ChainSpace{2, 2}, &gStar)

	SpaceIs(ChainSpace{2, 2}, 0, gStar, "Star", t)
	CoinsIs(3, 0, gStar, "Star", t)
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

	CoinsIs(7, 1, g, "", t)

	g.ExtraEvent.Handle(1, &g) //Move P2
	g.ExtraEvent.Handle(1, &g) //Move P3

	g.ExtraEvent.Handle(false, &g) //P0 is red
	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(3, &g) //P3 wins

	g.ExtraEvent.Handle(1, &g) //Move P0
	g.ExtraEvent.Handle(1, &g) //Move P1
	g.ExtraEvent.Handle(1, &g) //Move P2

	bd = g.Board.Data.(bmmBoardData)
	IntIs(1, bd.MagmaTurnCount, "Turn Count", t)

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

	CoinsIs(13, 0, g, "Coins", t)
}
