package board

import (
	"reflect"
	"testing"

	"github.com/0xhexnumbers/partysim/mp1"
)

func SpaceIs(expected mp1.ChainSpace, player int, g mp1.Game, flavour string, t *testing.T) {
	got := g.Players[player].CurrentSpace
	if expected != got {
		t.Errorf("Expected %s %d Space: %#v, got: %#v",
			flavour, player, expected, got)
	}
}

func StarsIs(expected, player int, g mp1.Game, flavour string, t *testing.T) {
	got := g.Players[player].Stars
	if expected != got {
		t.Errorf("Expected Player %d %s Stars: %d, got: %d",
			player, flavour, expected, got)
	}
}

func CoinsIs(expected, player int, g mp1.Game, flavour string, t *testing.T) {
	got := g.Players[player].Coins
	if expected != got {
		t.Errorf("Expected Player %d %s Coins: %d, got: %d",
			player, flavour, expected, got)
	}
}

func MinigameCoinsIs(expected, player int, g mp1.Game, flavour string, t *testing.T) {
	got := g.Players[player].MinigameCoins
	if expected != got {
		t.Errorf("Expected Player %d %s Coins: %d, got: %d",
			player, flavour, expected, got)
	}
}

func SpaceTypeIs(expected, got mp1.SpaceType, flavour string, t *testing.T) {
	if expected != got {
		t.Errorf("Expected %s Space Type: %d, got: %d",
			flavour, expected, got)
	}
}

func IntIs(expected, got int, flavour string, t *testing.T) {
	if expected != got {
		t.Errorf("Expected %s: %d, got: %d",
			flavour, expected, got)
	}
}

func EventIs(expected, got mp1.Event, flavour string, t *testing.T) {
	if expected != got {
		t.Errorf("Expected %s Event: %#v, got: %#v",
			flavour, expected, got)
	}
}

func ResIs(expected []mp1.Response, g mp1.Game, flavour string, t *testing.T) {
	got := g.NextEvent.Responses()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected %s Res: %#v, got: %#v",
			flavour, expected, got)
	}
}
