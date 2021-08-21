package mp1

import "testing"

func TestCannon(t *testing.T) {
	g := *InitializeGame(WBC, GameConfig{MaxTurns: 20})
	g.ExtraEvent.Handle(uint8(0), &g) //Set Star Space
	g.ExtraEvent.Handle(5, &g)        //Move
	g.ExtraEvent.Handle(11, &g)       //Land on {1, 12}
	g.ExtraEvent.Handle(5, &g)        //Move
	g.ExtraEvent.Handle(8, &g)        //Land on {3, 9}

	g.Players[2].CurrentSpace = ChainSpace{3, 12}
	g.ExtraEvent.Handle(1, &g) //Move
	g.ExtraEvent.Handle(0, &g) //Land on {0, 1}

	expectedSpaces := []ChainSpace{{1, 12}, {3, 9}, {0, 1}}
	for i := 0; i < 3; i++ {
		SpaceIs(expectedSpaces[i], i, g, "", t)
	}
}

func TestBowserCannon(t *testing.T) {
	g := *InitializeGame(WBC, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{4, 6}
	g.Players[0].Coins = 30

	g.ExtraEvent.Handle(uint8(0), &g) //Set Star Space
	g.ExtraEvent.Handle(1, &g)        //Move
	g.ExtraEvent.Handle(0, &g)        //Goto chain 0
	g.ExtraEvent.Handle(0, &g)        //Goto space 0

	CoinsIs(13, 0, g, "", t)
	SpaceIs(ChainSpace{0, 1}, 0, g, "", t)
}

func TestShyGuy(t *testing.T) {
	g := *InitializeGame(WBC, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{3, 4}

	g.ExtraEvent.Handle(uint8(0), &g) //Set Star Space
	g.ExtraEvent.Handle(1, &g)        //Move
	gNothing := g
	gNothing.ExtraEvent.Handle(wbcShyGuyResponse{
		wbcNothing, 0,
	}, &gNothing)
	SpaceIs(ChainSpace{3, 6}, 0, gNothing, "Nothing", t)

	gBowser := g
	gBowser.ExtraEvent.Handle(wbcShyGuyResponse{
		wbcFlyToBowser, 0,
	}, &gBowser)
	gBowser.ExtraEvent.Handle(2, &gBowser)
	SpaceIs(ChainSpace{4, 3}, 0, gBowser, "Bowser", t)

	gBringPlayer := g
	gBringPlayer.ExtraEvent.Handle(wbcShyGuyResponse{
		wbcBringPlayer, 1,
	}, &gBringPlayer)
	SpaceIs(ChainSpace{3, 6}, 0, gBringPlayer, "BringPlayer0", t)
	SpaceIs(ChainSpace{3, 4}, 1, gBringPlayer, "BringPlayer1", t)
}
