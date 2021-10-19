package board

import (
	"testing"

	"github.com/0xhexnumbers/partysim/mp1"
)

func TestCannon(t *testing.T) {
	g := *mp1.InitializeGame(WBC, mp1.GameConfig{MaxTurns: 20})
	g.NextEvent.Handle(mp1.NewChainSpace(0, 6), &g) //Set Star Space
	g.NextEvent.Handle(5, &g)                       //Move
	g.NextEvent.Handle(mp1.NewChainSpace(1, 11), &g)
	g.NextEvent.Handle(5, &g) //Move
	g.NextEvent.Handle(mp1.NewChainSpace(3, 8), &g)

	g.Players[2].CurrentSpace = mp1.NewChainSpace(3, 12)
	g.NextEvent.Handle(1, &g) //Move
	g.NextEvent.Handle(mp1.NewChainSpace(0, 0), &g)

	expectedSpaces := []mp1.ChainSpace{
		mp1.NewChainSpace(1, 12),
		mp1.NewChainSpace(3, 9),
		mp1.NewChainSpace(0, 1)}
	for i := 0; i < 3; i++ {
		SpaceIs(expectedSpaces[i], i, g, "", t)
	}
}

func TestBowserCannon(t *testing.T) {
	g := *mp1.InitializeGame(WBC, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(4, 6)
	g.Players[0].Coins = 30

	g.NextEvent.Handle(mp1.NewChainSpace(0, 6), &g) //Set Star Space
	g.NextEvent.Handle(1, &g)                       //Move
	g.NextEvent.Handle(mp1.NewChainSpace(0, 0), &g) //Goto {0, 0}

	CoinsIs(13, 0, g, "", t)
	SpaceIs(mp1.NewChainSpace(0, 1), 0, g, "", t)
}

func TestShyGuy(t *testing.T) {
	g := *mp1.InitializeGame(WBC, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(3, 4)

	g.NextEvent.Handle(mp1.NewChainSpace(0, 6), &g) //Set Star Space
	g.NextEvent.Handle(1, &g)                       //Move
	gNothing := g
	gNothing.NextEvent.Handle(WBCShyGuyResponse{
		WBCNothing, 0,
	}, &gNothing)
	SpaceIs(mp1.NewChainSpace(3, 6), 0, gNothing, "Nothing", t)

	gBowser := g
	gBowser.NextEvent.Handle(WBCShyGuyResponse{
		WBCFlyToBowser, 0,
	}, &gBowser)
	gBowser.NextEvent.Handle(mp1.NewChainSpace(4, 2), &gBowser)
	SpaceIs(mp1.NewChainSpace(4, 3), 0, gBowser, "Bowser", t)

	gBringPlayer := g
	gBringPlayer.NextEvent.Handle(WBCShyGuyResponse{
		WBCBringPlayer, 1,
	}, &gBringPlayer)
	SpaceIs(mp1.NewChainSpace(3, 6), 0, gBringPlayer, "BringPlayer0", t)
	SpaceIs(mp1.NewChainSpace(3, 4), 1, gBringPlayer, "BringPlayer1", t)
}
