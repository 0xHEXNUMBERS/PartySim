package mp1

import "testing"

func TestHiddenBlock(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 21}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{MaxTurns: 20, EventsDice: true},
	}
	g.MovePlayer(0, 1) //Blue Space

	expectedEvent := HiddenBlockEvent{0}
	gotEvent := g.ExtraEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected event: %#v, got: %#v",
			expectedEvent, gotEvent)
	}

	gHidden := g
	gHidden.ExtraEvent.Handle(true, &gHidden)

	expectedDiceEvent := EventDiceBlock{0}
	gotDiceEvent := gHidden.ExtraEvent
	if expectedDiceEvent != gotDiceEvent {
		t.Errorf("Expected Dice event: %#v, got: %#v",
			expectedDiceEvent, gotDiceEvent)
	}

	gBlue := g
	gBlue.ExtraEvent.Handle(false, &gBlue)

	expectedPlayerEvent := PickDiceBlock{1, g.Config}
	gotPlayerEvent := gBlue.ExtraEvent
	if expectedPlayerEvent != gotPlayerEvent {
		t.Errorf("Expected Player event: %#v, got: %#v",
			expectedPlayerEvent, gotPlayerEvent)
	}

	expectedCoins := 13
	gotCoins := gBlue.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected coins: %d, got: %d",
			expectedCoins, gotCoins)
	}
}
