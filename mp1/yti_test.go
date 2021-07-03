package mp1

import (
	"reflect"
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
	g.MovePlayer(0, 4)
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
	expected := BranchEvent{0, 1, 6, []ChainSpace{{1, 6}}}
	got := g.MovePlayer(0, 10)
	if reflect.DeepEqual(expected, got) {
		t.Errorf("Event type expected: %#v, got: %#v", expected, got)
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

	evt := g.MovePlayer(0, 10)
	if evt != nil {
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

	g.MovePlayer(0, 1)
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
	evt := g.MovePlayer(0, 10)
	//Move player to Chain 3 to pay thwomp 1
	mvmnt := evt.Handle(ChainSpace{3, 0}, &g)
	//Instantiate thwomp pay
	evt = g.MovePlayer(mvmnt.Player, mvmnt.Moves)
	//Pay thwomp 3 coins
	mvmnt = evt.Handle(3, &g)
	//Move remaining spaces and gain 3 coins
	evt = g.MovePlayer(0, mvmnt.Moves)

	if evt != nil {
		t.Errorf("Recieved unexpected event: %#v", evt)
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

func TestPassThwomp(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	evt := g.MovePlayer(0, 10)
	mvmnt := evt.Handle(nil, &g)
	evt = g.MovePlayer(mvmnt.Player, mvmnt.Moves)

	if evt != nil {
		t.Errorf("Recieved unexpected event: %#v", evt)
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
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	evt := g.MovePlayer(0, 3)
	if evt != nil {
		t.Errorf("Unexpected event: %#v", evt)
	}

	//starSpace := ChainSpace{1, 18}
	if g.Board.Chains[1][18].Type == BlackStar || g.Board.Chains[0][19].Type == Star {
		t.Errorf("Star spot did not swap, 1-18: %#v, 0-19: %#v",
			g.Board.Chains[1][18],
			g.Board.Chains[0][19],
		)
	}
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

	g.MovePlayer(0, 1)
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

	expected := MushroomEvent{0}
	got := g.MovePlayer(0, 4)
	if expected != got {
		t.Errorf("Expected event: %#v, got: %#v", expected, got)
	}

	//Received red mushroom
	expectedMvmnt := Movement{Skip: true}
	gotMvmnt := got.Handle(true, &g)
	if expectedMvmnt != gotMvmnt {
		t.Errorf("Expected Red Movement: %#v, got: %#v",
			expectedMvmnt,
			gotMvmnt,
		)
	}

	//Received poison mushroom
	expectedMvmnt = Movement{0, 0, false, nil}
	gotMvmnt = got.Handle(false, &g)
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
	evt := g.MovePlayer(0, 1)
	expectedMvmnt := Movement{
		ExtraEvent: BooCoinsEvent{
			PayRangeEvent: PayRangeEvent{
				Player: 1,
				Min:    1,
				Max:    15,
				Moves:  1,
			},
			RecvPlayer: 0,
		},
	}
	gotMvmnt := evt.Handle(BooStealAction{0, 1, false}, &g)
	if expectedMvmnt != gotMvmnt {
		t.Errorf("Expected movement: %#v, got: %#v",
			expectedMvmnt,
			gotMvmnt,
		)
	}

	expectedMvmnt.ExtraEvent.Handle(5, &g)
	expectedDaisyCoins := 15
	expectedLuigiCoins := 5
	gotDaisyCoins := g.Players[0].Coins
	gotLuigiCoins := g.Players[1].Coins

	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Daisy expected: %d coins, got: %d coins",
			expectedDaisyCoins,
			gotDaisyCoins,
		)
	}

	if expectedLuigiCoins != gotLuigiCoins {
		t.Errorf("Luigi expected: %d coins, got: %d coins",
			expectedLuigiCoins,
			gotLuigiCoins,
		)
	}
}
