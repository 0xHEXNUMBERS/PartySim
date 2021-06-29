package mp1

import (
	"testing"
)

func TestMove(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			{"Daisy", 0, 10, ChainSpace{1, 23}, false},
			{"Luigi", 0, 10, ChainSpace{0, 0}, false},
			{"Donkey Kong", 0, 10, ChainSpace{0, 0}, false},
			{"Mario", 0, 10, ChainSpace{0, 0}, false},
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
			{"Daisy", 0, 10, ChainSpace{1, 23}, false},
			{"Luigi", 0, 10, ChainSpace{0, 0}, false},
			{"Donkey Kong", 0, 10, ChainSpace{0, 0}, false},
			{"Mario", 0, 10, ChainSpace{0, 0}, false},
		},
	}
	expected := BranchEvent{0, 1, 6}
	got := g.MovePlayer(0, 10)
	if expected != got {
		t.Errorf("Event type expected: %#v, got: %#v", expected, got)
	}
}

func TestCanNotPayThwomp(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			{"Daisy", 0, 0, ChainSpace{1, 23}, false},
			{"Luigi", 0, 10, ChainSpace{0, 0}, false},
			{"Donkey Kong", 0, 10, ChainSpace{0, 0}, false},
			{"Mario", 0, 10, ChainSpace{0, 0}, false},
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
			{"Daisy", 0, 10, ChainSpace{1, 23}, false},
			{"Luigi", 0, 10, ChainSpace{0, 0}, false},
			{"Donkey Kong", 0, 10, ChainSpace{0, 0}, false},
			{"Mario", 0, 10, ChainSpace{0, 0}, false},
		},
	}

	g.MovePlayer(0, 1)
	expected := 13
	got := g.Players[0].Coins
	if expected != got {
		t.Errorf("Coins expected: %d, got: %d", expected, got)
	}
}
