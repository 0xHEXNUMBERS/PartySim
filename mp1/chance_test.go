package mp1

import "testing"

func TestChanceTime(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 4, 10, ChainSpace{0, 14}),
			NewPlayer("Luigi", 9, 25, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g = MovePlayer(g, 0, 1)
	//Steal from Luigi
	g = g.ExtraEvent.Handle(ChanceTimeResponse{CTBSide, 1}, g)
	//Daisy will recieve
	g = g.ExtraEvent.Handle(ChanceTimeResponse{CTBSide, 0}, g)

	g10Coins := g.ExtraEvent.Handle(ChanceTimeResponse{CTBMiddle, int(LTR10)}, g)
	expectedDaisyCoins := 20
	gotDaisyCoins := g10Coins.Players[0].Coins
	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Expected LTR10 Daisy Coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins,
		)
	}
	expectedLuigiCoins := 15
	gotLuigiCoins := g10Coins.Players[1].Coins
	if expectedLuigiCoins != gotLuigiCoins {
		t.Errorf("Expected LTR10 Luigi Coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins,
		)
	}

	g20Coins := g.ExtraEvent.Handle(ChanceTimeResponse{CTBMiddle, int(LTR20)}, g)
	expectedDaisyCoins = 30
	gotDaisyCoins = g20Coins.Players[0].Coins
	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Expected LTR20 Daisy Coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins,
		)
	}
	expectedLuigiCoins = 5
	gotLuigiCoins = g20Coins.Players[1].Coins
	if expectedLuigiCoins != gotLuigiCoins {
		t.Errorf("Expected LTR20 Luigi Coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins,
		)
	}

	g30Coins := g.ExtraEvent.Handle(ChanceTimeResponse{CTBMiddle, int(LTR30)}, g)
	expectedDaisyCoins = 35
	gotDaisyCoins = g30Coins.Players[0].Coins
	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Expected LTR30 Daisy Coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins,
		)
	}
	expectedLuigiCoins = 0
	gotLuigiCoins = g30Coins.Players[1].Coins
	if expectedLuigiCoins != gotLuigiCoins {
		t.Errorf("Expected LTR30 Luigi Coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins,
		)
	}

	gStar := g.ExtraEvent.Handle(ChanceTimeResponse{CTBMiddle, int(LTRStar)}, g)
	expectedDaisyStars := 5
	gotDaisyStars := gStar.Players[0].Stars
	if expectedDaisyStars != gotDaisyStars {
		t.Errorf("Expected LTRStar Daisy Stars: %d, got: %d",
			expectedDaisyStars, gotDaisyStars,
		)
	}
	expectedLuigiStars := 8
	gotLuigiStars := gStar.Players[1].Stars
	if expectedLuigiStars != gotLuigiStars {
		t.Errorf("Expected LTRStar Luigi Stars: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins,
		)
	}

	gSwapCoins := g.ExtraEvent.Handle(ChanceTimeResponse{CTBMiddle, int(SwapCoins)}, g)
	expectedDaisyCoins = 25
	gotDaisyCoins = gSwapCoins.Players[0].Coins
	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Expected SwapCoins Daisy Coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins,
		)
	}
	expectedLuigiCoins = 10
	gotLuigiCoins = gSwapCoins.Players[1].Coins
	if expectedLuigiCoins != gotLuigiCoins {
		t.Errorf("Expected SwapCoins Luigi Coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins,
		)
	}

	gSwapStars := g.ExtraEvent.Handle(ChanceTimeResponse{CTBMiddle, int(SwapStars)}, g)
	expectedDaisyStars = 9
	gotDaisyStars = gSwapStars.Players[0].Stars
	if expectedDaisyStars != gotDaisyStars {
		t.Errorf("Expected SwapStars Daisy Stars: %d, got: %d",
			expectedDaisyStars, gotDaisyStars,
		)
	}
	expectedLuigiStars = 4
	gotLuigiStars = gSwapStars.Players[1].Stars
	if expectedLuigiStars != gotLuigiStars {
		t.Errorf("Expected SwapStars Luigi Stars: %d, got: %d",
			expectedDaisyStars, gotDaisyStars,
		)
	}

	if gSwapStars.ExtraEvent != nil {
		t.Errorf("Unexpected event: %#v", gSwapStars.ExtraEvent)
	}
}
