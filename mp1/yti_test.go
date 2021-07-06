package mp1

import (
	"testing"
)

func TestMove(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g = MovePlayer(g, 0, 4)
	expected := ChainSpace{1, 27}
	got := g.Players[0].CurrentSpace
	if expected != got {
		t.Errorf("Position expected: %#v, got: %#v", expected, got)
	}

}

func TestCanPayThwomp(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g = MovePlayer(g, 0, 10)
	expected := BranchEvent{0, 1, 6, YTI.Links[1]}
	got := g.ExtraEvent
	if expected != got {
		t.Errorf("Event expected: %#v, got: %#v", expected, got)
	}
}

func TestCanNotPayThwomp(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 0, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g = MovePlayer(g, 0, 10)
	if g.ExtraEvent != nil {
		t.Error("Could not pay thwomp, yet recieved a branch event")
	}
}

func TestGainCoins(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g = MovePlayer(g, 0, 1)
	expected := 13
	got := g.Players[0].Coins
	if expected != got {
		t.Errorf("Coins expected: %d, got: %d", expected, got)
	}
}

func TestPayThwompAndGainCoins(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	//Move player to invisible space
	g = MovePlayer(g, 0, 10)
	//Move player to Chain 3 to pay thwomp 1
	g = g.ExtraEvent.Handle(ChainSpace{3, 0}, g)
	//Instantiate thwomp pay
	g = MovePlayer(g, g.ExtraMovement.Player, g.ExtraMovement.Moves)
	//Pay thwomp 3 coins
	g = g.ExtraEvent.Handle(3, g)
	//Move remaining spaces and gain 3 coins
	g = MovePlayer(g, g.ExtraMovement.Player, g.ExtraMovement.Moves)

	if g.ExtraEvent != nil {
		t.Errorf("Recieved unexpected event: %#v", g.ExtraEvent)
	}

	expectedSquare := ChainSpace{0, 12}
	gotSquare := g.Players[0].CurrentSpace
	if expectedSquare != gotSquare {
		t.Errorf("Space expected: %#v, got: %#v", expectedSquare, gotSquare)
	}

	expectedCoins := 10
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %d, got: %d", expectedCoins, gotCoins)
	}
}

func TestIgnoreThwomp(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g = MovePlayer(g, 0, 10)
	g = g.ExtraEvent.Handle(nil, g)
	g = MovePlayer(g, g.ExtraMovement.Player, g.ExtraMovement.Moves)

	if g.ExtraEvent != nil {
		t.Errorf("Recieved unexpected event: %#v", g.ExtraEvent)
	}

	expectedSquare := ChainSpace{1, 5}
	gotSquare := g.Players[0].CurrentSpace
	if expectedSquare != gotSquare {
		t.Errorf("Space expected: %#v, got: %#v", expectedSquare, gotSquare)
	}

	expectedCoins := 13
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %d, got: %d", expectedCoins, gotCoins)
	}
}

func TestStarSwapViaHappening(t *testing.T) {
	t.SkipNow()
	/*g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g = MovePlayer(g, 0, 3)
	if g.ExtraEvent != nil {
		t.Errorf("Unexpected event: %#v", g.ExtraEvent)
	}

	//starSpace := ChainSpace{1, 18}
	if g.Board.Chains[1][18].Type == BlackStar || g.Board.Chains[0][19].Type == Star {
		t.Errorf("Star spot did not swap, 1-18: %#v, 0-19: %#v",
			g.Board.Chains[1][18],
			g.Board.Chains[0][19],
		)
	}*/
}

func TestCoinsOnStart(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 22}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		CoinsOnStart: true,
	}

	g = MovePlayer(g, 0, 1)
	expectedCoins := 20
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Coins expected: %d, got: %d", expectedCoins, gotCoins)
	}
}

func TestMushroomSpace(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g = MovePlayer(g, 0, 4)
	expected := MushroomEvent{0}
	got := g.ExtraEvent
	if expected != got {
		t.Errorf("Expected event: %#v, got: %#v", expected, got)
	}

	//Received red mushroom
	g = got.Handle(true, g)
	expectedMvmnt := Movement{Skip: true}
	gotMvmnt := g.ExtraMovement
	if expectedMvmnt != gotMvmnt {
		t.Errorf("Expected Red Movement: %#v, got: %#v",
			expectedMvmnt,
			gotMvmnt,
		)
	}

	//Received poison mushroom
	g = got.Handle(false, g)
	expectedMvmnt = Movement{0, 0, false}
	gotMvmnt = g.ExtraMovement
	if expectedMvmnt != gotMvmnt {
		t.Errorf("Expected Poison Movement: %#v, got: %#v",
			expectedMvmnt,
			gotMvmnt,
		)
	}

	if g.Players[0].SkipTurn != true {
		t.Errorf("SkipTurn not set")
	}
}

func TestStealCoinsViaBoo(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g = MovePlayer(g, 0, 1)
	g = g.ExtraEvent.Handle(BooStealAction{0, 1, false}, g)
	expectedEvent := BooCoinsEvent{
		PayRangeEvent: PayRangeEvent{
			Player: 1,
			Min:    1,
			Max:    10, //Max of 10 coins
			Moves:  1,
		},
		RecvPlayer: 0,
	}
	gotEvent := g.ExtraEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected movement: %#v, got: %#v",
			expectedEvent,
			gotEvent,
		)
	}

	g = expectedEvent.Handle(5, g)
	expectedDaisyCoins := 15
	gotDaisyCoins := g.Players[0].Coins

	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Daisy expected: %d coins, got: %d coins",
			expectedDaisyCoins,
			gotDaisyCoins,
		)
	}

	expectedLuigiCoins := 5
	gotLuigiCoins := g.Players[1].Coins
	if expectedLuigiCoins != gotLuigiCoins {
		t.Errorf("Luigi expected: %d coins, got: %d coins",
			expectedLuigiCoins,
			gotLuigiCoins,
		)
	}
}

func TestStealTooManyCoinsViaBoo(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 4, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g = MovePlayer(g, 0, 1)
	g = g.ExtraEvent.Handle(BooStealAction{0, 1, false}, g)
	expectedEvent := BooCoinsEvent{
		PayRangeEvent: PayRangeEvent{
			Player: 1,
			Min:    1,
			Max:    4, //Max of 4 coins
			Moves:  1,
		},
		RecvPlayer: 0,
	}
	gotEvent := g.ExtraEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected movement: %#v, got: %#v",
			expectedEvent,
			gotEvent,
		)
	}

	g = expectedEvent.Handle(4, g)
	expectedDaisyCoins := 14
	gotDaisyCoins := g.Players[0].Coins

	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Daisy expected: %d coins, got: %d coins",
			expectedDaisyCoins,
			gotDaisyCoins,
		)
	}

	expectedLuigiCoins := 0
	gotLuigiCoins := g.Players[1].Coins
	if expectedLuigiCoins != gotLuigiCoins {
		t.Errorf("Luigi expected: %d coins, got: %d coins",
			expectedLuigiCoins,
			gotLuigiCoins,
		)
	}
}
