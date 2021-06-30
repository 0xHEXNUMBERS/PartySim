package mp1

import (
	"reflect"
	"testing"
)

func TestMove(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			{"Daisy", 0, 10, ChainSpace{1, 23}, false, 0, 0, 0},
			{"Luigi", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Donkey Kong", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Mario", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
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
			{"Daisy", 0, 10, ChainSpace{1, 23}, false, 0, 0, 0},
			{"Luigi", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Donkey Kong", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Mario", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
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
			{"Daisy", 0, 0, ChainSpace{1, 23}, false, 0, 0, 0},
			{"Luigi", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Donkey Kong", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Mario", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
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
			{"Daisy", 0, 10, ChainSpace{1, 23}, false, 0, 0, 0},
			{"Luigi", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Donkey Kong", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Mario", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
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
			{"Daisy", 0, 10, ChainSpace{1, 23}, false, 0, 0, 0},
			{"Luigi", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Donkey Kong", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
			{"Mario", 0, 10, ChainSpace{0, 0}, false, 0, 0, 0},
		},
	}

	//Move player to invisible space
	evt := g.MovePlayer(0, 10)
	//Move player to Chain 3 to pay thwomp 1
	mvmnt := g.Board.EventHandler(evt, ChainSpace{3, 0}, &g)
	//Instantiate thwomp pay
	evt = g.MovePlayer(mvmnt.Player, mvmnt.Moves)
	//Pay thwomp 3 coins
	mvmnt = g.Board.EventHandler(evt, 3, &g)
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
