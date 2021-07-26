package mp1

import "testing"

func TestBowser10CoinsForStar(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 1, 0, ChainSpace{1, 14}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	expectedStars := 0
	gotStars := g.Players[0].Stars
	if expectedStars != gotStars {
		t.Errorf("Event expected: %#v, got: %#v", expectedStars, gotStars)
	}

	expectedCoins := 10
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %#v, got: %#v", expectedCoins, gotCoins)
	}
}

func TestBowserGain20Coins(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 0, ChainSpace{1, 14}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	expectedCoins := 20
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %#v, got: %#v", expectedCoins, gotCoins)
	}
}

func TestCoinsForBowser(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 25, ChainSpace{1, 14}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	g.ExtraEvent.Handle(CoinsForBowser, &g)
	expected := PickDiceBlock{1, g.Config}
	got := g.ExtraEvent
	if expected != got {
		t.Errorf("Event expected: %#v, got: %#v", expected, got)
	}

	expectedCoins := 15
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %#v, got: %#v", expectedCoins, gotCoins)
	}
}

func TestBowserBalloonBurst(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 50, ChainSpace{1, 14}),
			NewPlayer("Luigi", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 50, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	g.ExtraEvent.Handle(BowserBalloonBurst, &g)
	gDraw := g
	gDraw.ExtraEvent.Handle(BBBDraw, &gDraw)
	expectedDrawCoins := 30
	for _, p := range gDraw.Players {
		gotDrawCoins := p.Coins
		if expectedDrawCoins != gotDrawCoins {
			t.Errorf("Player: %s\tDraw Coins expected: %#v, got: %#v", p.Char, expectedDrawCoins, gotDrawCoins)
		}
	}

	gP1Win := g
	gP1Win.ExtraEvent.Handle(BBBP1Win, &gP1Win)
	expectedP1WinCoins := 50
	gotP1WinCoins := gP1Win.Players[0].Coins
	if expectedP1WinCoins != gotP1WinCoins {
		t.Errorf("Win Coins expected: %#v, got: %#v", expectedP1WinCoins, gotP1WinCoins)
	}
	expectedLossCoins := 40 //50 - 10 coins at turn 0
	for _, p := range gP1Win.Players {
		if p.Char == "Daisy" {
			continue
		}
		gotLossCoins := p.Coins
		if expectedLossCoins != gotLossCoins {
			t.Errorf("Player: %s\tLoss Coins expected: %#v, got: %#v", p.Char, expectedLossCoins, gotLossCoins)
		}
	}
}

func TestBowsersFaceList(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 50, ChainSpace{1, 14}),
			NewPlayer("Luigi", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 50, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	g.ExtraEvent.Handle(BowsersFaceLift, &g)
	gDraw := g
	gDraw.ExtraEvent.Handle(0b1111, &gDraw) //Draw
	expectedCoins := 0
	gotCoins := gDraw.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Draw Coins expected: %d, got: %d", expectedCoins, gotCoins)
	}

	gP1Loss := g
	gP1Loss.ExtraEvent.Handle(0b1110, &gP1Loss) //All players except Daisy
	expectedLossCoins := 40                     //50 - 10 coins at turn 0
	gotLossCoins := gP1Loss.Players[0].Coins
	if expectedLossCoins != gotLossCoins {
		t.Errorf("Loss Coins expected: %d, got: %d", expectedLossCoins, gotLossCoins)
	}
}

func TestBowsersTugoWar(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 50, ChainSpace{1, 14}),
			NewPlayer("Luigi", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 50, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	g.ExtraEvent.Handle(BowsersTugoWar, &g)
	gDraw := g
	gDraw.ExtraEvent.Handle(BTWDraw, &gDraw)
	expectedDrawCoins := 20
	for _, p := range gDraw.Players {
		gotDrawCoins := p.Coins
		if expectedDrawCoins != gotDrawCoins {
			t.Errorf("Player: %s\tDraw Coins expected: %d, got: %d", p.Char, expectedDrawCoins, gotDrawCoins)
		}
	}

	g1TWin := g
	g1TWin.ExtraEvent.Handle(BTW1TWin, &g1TWin)
	expected1TCoins := 40 //50 - 10 coins at turn 0
	for i := 1; i < len(g1TWin.Players); i++ {
		got1TCoins := g1TWin.Players[i].Coins
		if expected1TCoins != got1TCoins {
			t.Errorf("Player: %s\t1T Coins expected: %d, got: %d", g1TWin.Players[i].Char, expected1TCoins, got1TCoins)
		}
	}

	g3TWin := g
	g3TWin.ExtraEvent.Handle(BTW3TWin, &g3TWin)
	expected3TCoins := 40 //50 - 10 coins at turn 0
	got3TCoins := g3TWin.Players[0].Coins
	if expected3TCoins != got3TCoins {
		t.Errorf("3T Coins expected: %d, got: %d", expected3TCoins, got3TCoins)
	}

}

func TestBashnCash(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 54, ChainSpace{1, 14}),
			NewPlayer("Luigi", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 50, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	g.ExtraEvent.Handle(BashnCash, &g)
	gGE5 := g
	gGE5.ExtraEvent.Handle(5, &gGE5) //Should lose 5 * 5 = 25 coins
	expectedGE5Coins := 29           //54 - 25 = 29
	gotGE5Coins := gGE5.Players[0].Coins
	if expectedGE5Coins != gotGE5Coins {
		t.Errorf("GE5Coins expected: %d, got: %d", expectedGE5Coins, gotGE5Coins)
	}

	gE5 := g
	gE5.ExtraEvent.Handle(10, &gE5) //Should lose 5 * 10 = 50 coins
	expectedE5Coins := 4            //54 - 50 = 4
	gotE5Coins := gE5.Players[0].Coins
	if expectedE5Coins != gotE5Coins {
		t.Errorf("E5Coins expected: %d, got: %d", expectedE5Coins, gotE5Coins)
	}

	gLT5 := g
	gLT5.ExtraEvent.Handle(13, &gLT5) //Should lose 5 * 10 + 3 = 53 coins
	expectedLT5Coins := 1             //54 - 53 = 1
	gotLT5Coins := gLT5.Players[0].Coins
	if expectedLT5Coins != gotLT5Coins {
		t.Errorf("LT5Coins expected: %d, got: %d", expectedLT5Coins, gotLT5Coins)
	}
}

func TestBowserRevolution(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 75, ChainSpace{1, 14}),
			NewPlayer("Luigi", 0, 25, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 99, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 50, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	g.ExtraEvent.Handle(BowserRevolution, &g)
	expectedCoins := 62 //25 + 50 + 75 + 99 = 249 // 4 = 62
	for _, p := range g.Players {
		gotCoins := p.Coins
		if expectedCoins != gotCoins {
			t.Errorf("Player: %s\tCoins expected: %d, got: %d", p.Char, expectedCoins, gotCoins)
		}
	}
}

func TestBowsersChanceTime(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 50, ChainSpace{1, 14}),
			NewPlayer("Luigi", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 50, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 50, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	g.ExtraEvent.Handle(BowsersChanceTime, &g)
	g.ExtraEvent.Handle(BCTResponse{0, 20}, &g) //Daisy
	expectedCoins := 30
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %d, got: %d", expectedCoins, gotCoins)
	}
}
