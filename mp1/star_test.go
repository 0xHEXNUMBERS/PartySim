package mp1

import (
	"reflect"
	"testing"
)

var multipleStarBoard = MakeRepeatedBoard(Star, 5)

var singleStarBoard = Board{
	Chains: &[]Chain{
		{
			{Type: Start},
			{Type: Star},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
			{Type: Blue},
		},
	},
}

func TestMultipleStarSpaces(t *testing.T) {
	g := InitializeGame(multipleStarBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Char = "Daisy"
	g.Players[1].Char = "Luigi"
	g.Players[2].Char = "Donkey Kong"
	g.Players[3].Char = "Mario"
	g.Players[0].Coins = 20
	g.Players[1].Coins = 20
	g.Players[2].Coins = 20
	g.Players[3].Coins = 20

	expectedRes := []Response{uint8(0), uint8(1), uint8(2), uint8(3), uint8(4)}
	gotRes := g.NextEvent.Responses()
	if !reflect.DeepEqual(expectedRes, gotRes) {
		t.Errorf("Expected responses: %#v, got: %#v",
			expectedRes, gotRes)
	}

	g.NextEvent.Handle(uint8(1), g) //Set {0, 2} as current star space
	g.NextEvent.Handle(2, g)        //Move player to {0, 3}, gaining star
	g.NextEvent.Handle(uint8(3), g) //Set {0, 4} as current star space
	//Player 0 should have 1 star, and 3 coins
	expectedStars := 1
	gotStars := g.Players[0].Stars
	expectedCoins := 3
	gotCoins := g.Players[0].Coins
	if expectedStars != gotStars || expectedCoins != gotCoins {
		t.Errorf("Expected {star:coin} count: {%d:%d}, got: {%d:%d}",
			expectedStars, expectedCoins, gotStars, gotCoins,
		)
	}

	g.NextEvent.Handle(2, g) //Move player to {0, 2} landing on chance time
	//Chance time should be happening
	expectedEvent := ChanceTime{Player: 1}
	gotEvent := g.NextEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected event: %#v, got: %#v",
			expectedEvent, gotEvent,
		)
	}
}

func TestSingleStarSpace(t *testing.T) {
	g := InitializeGame(singleStarBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Char = "Daisy"
	g.Players[1].Char = "Luigi"
	g.Players[2].Char = "Donkey Kong"
	g.Players[3].Char = "Mario"
	g.Players[0].Coins = 20
	g.Players[1].Coins = 20
	g.Players[2].Coins = 20
	g.Players[3].Coins = 20

	g.NextEvent.Handle(2, g) //Move player to {0, 3}, gaining star
	//Player 0 should have 1 star, and 3 coins
	expectedStars := 1
	gotStars := g.Players[0].Stars
	expectedCoins := 3
	gotCoins := g.Players[0].Coins
	if expectedStars != gotStars || expectedCoins != gotCoins {
		t.Errorf("Expected Daisy {star:coin} count: {%d:%d}, got: {%d:%d}",
			expectedStars, expectedCoins, gotStars, gotCoins,
		)
	}

	g.NextEvent.Handle(2, g) //Move player to {0, 3}, gaining star
	gotStars = g.Players[1].Stars
	gotCoins = g.Players[1].Coins
	if expectedStars != gotStars || expectedCoins != gotCoins {
		t.Errorf("Expected Luigi {star:coin} count: {%d:%d}, got: {%d:%d}",
			expectedStars, expectedCoins, gotStars, gotCoins,
		)
	}
}
