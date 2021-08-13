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
		gotSpace := g.Players[i].CurrentSpace
		if expectedSpaces[i] != gotSpace {
			t.Errorf("Expected %d space: %#v, got: %#v",
				i, expectedSpaces[i], gotSpace)
		}
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

	expectedCoins := 13
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected coins: %d, got: %d",
			expectedCoins, gotCoins)
	}

	expectedSpace := ChainSpace{0, 1}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
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
	expectedSpace := ChainSpace{3, 6}
	gotSpace := gNothing.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected Nothing space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}

	gBowser := g
	gBowser.ExtraEvent.Handle(wbcShyGuyResponse{
		wbcFlyToBowser, 0,
	}, &gBowser)
	gBowser.ExtraEvent.Handle(2, &gBowser)
	expectedSpace = ChainSpace{4, 3}
	gotSpace = gBowser.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected Bowser space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}

	gBringPlayer := g
	gBringPlayer.ExtraEvent.Handle(wbcShyGuyResponse{
		wbcBringPlayer, 1,
	}, &gBringPlayer)
	expectedSpace = ChainSpace{3, 6}
	gotSpace = gBringPlayer.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected BringPlayer0 space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}

	expectedSpace = ChainSpace{3, 4}
	gotSpace = gBringPlayer.Players[1].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected BringPlayer1 space: %#v, got: %#v",
			expectedSpace, gotSpace)
	}
}
