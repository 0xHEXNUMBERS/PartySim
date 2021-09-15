package board

import (
	"testing"

	"github.com/0xhexnumbers/partysim/mp1"
)

func TestSeedCheckAutoStar(t *testing.T) {
	g := *mp1.InitializeGame(PBC, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 46)
	g.Players[1].CurrentSpace = mp1.NewChainSpace(0, 46)
	g.Players[2].CurrentSpace = mp1.NewChainSpace(0, 46)
	g.Players[3].CurrentSpace = mp1.NewChainSpace(0, 46)

	g.NextEvent.Handle(2, &g)    //Move
	g.NextEvent.Handle(true, &g) //Picked Bowser
	g.NextEvent.Handle(2, &g)    //Move: Toad
	g.NextEvent.Handle(2, &g)    //Move: Toad
	g.NextEvent.Handle(2, &g)    //Move: Toad

	SpaceIs(mp1.NewChainSpace(1, 1), 0, g, "0", t)
	for i := 1; i < 4; i++ {
		expectedPos := mp1.NewChainSpace(0, 1)
		gotPos := g.Players[i].CurrentSpace
		if expectedPos != gotPos {
			t.Errorf("Expected %d pos: %#v, got: %#v",
				i, expectedPos, gotPos)
		}
	}
	for i := 0; i < 4; i++ {
		expectedCoins := 3
		gotCoins := g.Players[i].Coins
		if expectedCoins != gotCoins {
			t.Errorf("Expected %d coins: %d, got: %d",
				i, expectedCoins, gotCoins)
		}
	}
}

func TestSeedCheckAutoBowser(t *testing.T) {
	g := *mp1.InitializeGame(PBC, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 46)
	g.Players[1].CurrentSpace = mp1.NewChainSpace(0, 46)
	g.Players[2].CurrentSpace = mp1.NewChainSpace(0, 46)
	g.Players[3].CurrentSpace = mp1.NewChainSpace(0, 46)

	g.NextEvent.Handle(2, &g)     //Move
	g.NextEvent.Handle(false, &g) //Picked Toad
	g.NextEvent.Handle(2, &g)     //Move
	g.NextEvent.Handle(false, &g) //Picked Toad
	g.NextEvent.Handle(2, &g)     //Move
	g.NextEvent.Handle(false, &g) //Picked Toad
	g.NextEvent.Handle(2, &g)     //Move: Bowser

	for i := 0; i < 3; i++ {
		SpaceIs(mp1.NewChainSpace(0, 1), i, g, "", t)
	}
	SpaceIs(mp1.NewChainSpace(1, 1), 3, g, "", t)
	for i := 0; i < 4; i++ {
		CoinsIs(3, i, g, "", t)
	}
}

func TestPiranha(t *testing.T) {
	g := *mp1.InitializeGame(PBC, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 18)
	g.Players[1].CurrentSpace = mp1.NewChainSpace(0, 18)
	g.Players[2].CurrentSpace = mp1.NewChainSpace(0, 18)
	g.Players[3].CurrentSpace = mp1.NewChainSpace(0, 18)
	g.Players[0].Coins = 30
	g.Players[3].Stars = 1

	g.NextEvent.Handle(1, &g)    //Move to unoccupied space
	g.NextEvent.Handle(true, &g) //Pay 30 coins
	g.NextEvent.Handle(2, &g)    //Move to unoccupied space with <30 coins
	g.NextEvent.Handle(1, &g)    //Move to occupied space with no stars
	g.NextEvent.Handle(1, &g)    //Move to occupied space with stars
	StarsIs(1, 0, g, "", t)
	CoinsIs(0, 0, g, "0Coins", t)
	StarsIs(0, 3, g, "3Stars", t)
}
