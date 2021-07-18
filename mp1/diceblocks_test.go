package mp1

import "testing"

func TestRedDiceBlock(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		ExtraEvent: RedDiceBlock{0},
	}
	g = g.ExtraEvent.Handle(9, g) //Land on minigame space
	expectedSpace := ChainSpace{0, 9}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v", expectedSpace, gotSpace)
	}

	expectedCoins := 1
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected coins: %d, got: %d", expectedCoins, gotCoins)
	}
}

func TestBlueDiceBlock(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		ExtraEvent: BlueDiceBlock{0},
	}
	g = g.ExtraEvent.Handle(9, g) //Land on minigame space
	expectedSpace := ChainSpace{0, 9}
	gotSpace := g.Players[0].CurrentSpace
	if expectedSpace != gotSpace {
		t.Errorf("Expected space: %#v, got: %#v", expectedSpace, gotSpace)
	}

	expectedCoins := 19
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected coins: %d, got: %d", expectedCoins, gotCoins)
	}
}

func TestWarpDiceBlock(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		ExtraEvent: WarpDiceBlock{0},
	}
	g = g.ExtraEvent.Handle(1, g) //Swap with Luigi
	expectedDaisySpace := ChainSpace{0, 0}
	gotDaisySpace := g.Players[0].CurrentSpace
	if expectedDaisySpace != gotDaisySpace {
		t.Errorf("Expected Daisy space: %#v, got: %#v", expectedDaisySpace, gotDaisySpace)
	}

	expectedLuigiSpace := ChainSpace{1, 23}
	gotLuigiSpace := g.Players[1].CurrentSpace
	if expectedLuigiSpace != gotLuigiSpace {
		t.Errorf("Expected Luigi space: %#v, got: %#v", expectedLuigiSpace, gotLuigiSpace)
	}
}

func TestEventDiceBlock(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		ExtraEvent: EventDiceBlock{0},
	}
	gBoo := g.ExtraEvent.Handle(BooEventBlock, g)
	expectedBooEvent := BooEvent{0, gBoo.Players, 0, gBoo.Players[0].Coins}
	gotBooEvent := gBoo.ExtraEvent
	if expectedBooEvent != gotBooEvent {
		t.Errorf("Expected Boo event: %#v, got: %#v", expectedBooEvent, gotBooEvent)
	}

	gBowser := g.ExtraEvent.Handle(BowserEventBlock, g)
	expectedBowserEvent := BowserEvent{0}
	gotBowserEvent := gBowser.ExtraEvent
	if expectedBowserEvent != gotBowserEvent {
		t.Errorf("Expected Bowser event: %#v, got: %#v", expectedBowserEvent, gotBowserEvent)
	}

	gKoopa := g.ExtraEvent.Handle(KoopaEventBlock, g)
	if gKoopa.ExtraEvent != nil {
		t.Errorf("Unexpected Koopa event: %#v", gKoopa.ExtraEvent)
	}
}
