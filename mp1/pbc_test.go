package mp1

import "testing"

func TestSeedCheckAutoStar(t *testing.T) {
	g := *InitializeGame(PBC, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{0, 46}
	g.Players[1].CurrentSpace = ChainSpace{0, 46}
	g.Players[2].CurrentSpace = ChainSpace{0, 46}
	g.Players[3].CurrentSpace = ChainSpace{0, 46}

	g.ExtraEvent.Handle(2, &g)    //Move
	g.ExtraEvent.Handle(true, &g) //Picked Bowser
	g.ExtraEvent.Handle(2, &g)    //Move: Toad
	g.ExtraEvent.Handle(2, &g)    //Move: Toad
	g.ExtraEvent.Handle(2, &g)    //Move: Toad

	expectedPos := ChainSpace{1, 1}
	gotPos := g.Players[0].CurrentSpace
	if expectedPos != gotPos {
		t.Errorf("Expected 0 pos: %#v, got: %#v",
			expectedPos, gotPos)
	}
	for i := 1; i < 4; i++ {
		expectedPos := ChainSpace{0, 1}
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
	g := *InitializeGame(PBC, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{0, 46}
	g.Players[1].CurrentSpace = ChainSpace{0, 46}
	g.Players[2].CurrentSpace = ChainSpace{0, 46}
	g.Players[3].CurrentSpace = ChainSpace{0, 46}

	g.ExtraEvent.Handle(2, &g)     //Move
	g.ExtraEvent.Handle(false, &g) //Picked Toad
	g.ExtraEvent.Handle(2, &g)     //Move
	g.ExtraEvent.Handle(false, &g) //Picked Toad
	g.ExtraEvent.Handle(2, &g)     //Move
	g.ExtraEvent.Handle(false, &g) //Picked Toad
	g.ExtraEvent.Handle(2, &g)     //Move: Bowser

	for i := 0; i < 3; i++ {
		expectedPos := ChainSpace{0, 1}
		gotPos := g.Players[i].CurrentSpace
		if expectedPos != gotPos {
			t.Errorf("Expected %d pos: %#v, got: %#v",
				i, expectedPos, gotPos)
		}
	}
	expectedPos := ChainSpace{1, 1}
	gotPos := g.Players[3].CurrentSpace
	if expectedPos != gotPos {
		t.Errorf("Expected 3 pos: %#v, got: %#v",
			expectedPos, gotPos)
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

func TestPiranha(t *testing.T) {
	g := *InitializeGame(PBC, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{0, 18}
	g.Players[1].CurrentSpace = ChainSpace{0, 18}
	g.Players[2].CurrentSpace = ChainSpace{0, 18}
	g.Players[3].CurrentSpace = ChainSpace{0, 18}
	g.Players[0].Coins = 30
	g.Players[3].Stars = 1

	g.ExtraEvent.Handle(1, &g)    //Move to unoccupied space
	g.ExtraEvent.Handle(true, &g) //Pay 30 coins
	g.ExtraEvent.Handle(2, &g)    //Move to unoccupied space with <30 coins
	g.ExtraEvent.Handle(1, &g)    //Move to occupied space with no stars
	g.ExtraEvent.Handle(1, &g)    //Move to occupied space with stars
	expectedStars0 := 1
	gotStars0 := g.Players[0].Stars
	if expectedStars0 != gotStars0 {
		t.Errorf("Expected 0 stars: %d, got: %d",
			expectedStars0, gotStars0)
	}
	expectedCoins0 := 0
	gotCoins0 := g.Players[0].Coins
	if expectedCoins0 != gotCoins0 {
		t.Errorf("Expected 0 coins: %d, got: %d",
			expectedCoins0, gotCoins0)
	}
	expectedStars3 := 0
	gotStars3 := g.Players[3].Stars
	if expectedStars3 != gotStars3 {
		t.Errorf("Expected 3 stars: %d, got: %d",
			expectedStars3, gotStars3)
	}
}
