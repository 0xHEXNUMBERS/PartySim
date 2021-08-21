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
	g.MovePlayer(0, 4)
	SpaceIs(ChainSpace{1, 27}, 0, g, "", t)
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
	g.MovePlayer(0, 10)
	EventIs(ytiThwompBranchEvent{0, 6, 1}, g.ExtraEvent, "", t)
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

	g.MovePlayer(0, 10)
	EventIs(NormalDiceBlock{1}, g.ExtraEvent, "", t)
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
		Config: GameConfig{
			MaxTurns: 25,
		},
	}

	g.MovePlayer(0, 1)
	CoinsIs(13, 0, g, "", t)
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
		Config: GameConfig{
			MaxTurns: 25,
		},
	}

	//Move player to invisible space
	g.MovePlayer(0, 10)
	//Accept payment to thwomp 1
	g.ExtraEvent.Handle(true, &g)
	//Pay thwomp 3 coins, move and land on blue space
	g.ExtraEvent.Handle(3, &g)

	EventIs(NormalDiceBlock{1}, g.ExtraEvent, "", t)
	SpaceIs(ChainSpace{0, 12}, 0, g, "", t)
	CoinsIs(10, 0, g, "", t)
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
		Config: GameConfig{
			MaxTurns: 25,
		},
	}

	g.MovePlayer(0, 10)
	g.ExtraEvent.Handle(false, &g)

	EventIs(NormalDiceBlock{1}, g.ExtraEvent, "", t)
	SpaceIs(ChainSpace{1, 5}, 0, g, "", t)
	CoinsIs(13, 0, g, "", t)
}

func TestStarSwapViaHappening(t *testing.T) {
	t.SkipNow()
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}

	g.MovePlayer(0, 3)
	EventIs(NormalDiceBlock{1}, g.ExtraEvent, "", t)

	bd := g.Board.Data.(ytiBoardData)
	if bd.StarPosition == ytiRightIslandStar {
		t.Errorf("Expected star position: %#v, got: %#v", ytiRightIslandStar, bd.StarPosition)
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
	}

	g.MovePlayer(0, 1)
	CoinsIs(20, 0, g, "", t)
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

	g.MovePlayer(0, 4)
	EventIs(MushroomEvent{0}, g.ExtraEvent, "", t)

	//Received red mushroom
	gRed := g
	gRed.ExtraEvent.Handle(true, &gRed)
	EventIs(NormalDiceBlock{0}, gRed.ExtraEvent, "", t)

	//Received poison mushroom
	gPoison := g
	gPoison.ExtraEvent.Handle(false, &gPoison)
	EventIs(NormalDiceBlock{1}, gPoison.ExtraEvent, "", t)
	if !gPoison.Players[0].SkipTurn {
		t.Errorf("SkipTurn not set")
	}
}

func TestSkipTurnViaMinigame(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 7}),
		},
		ExtraEvent: NormalDiceBlock{0},
		Config:     GameConfig{MaxTurns: 20},
	}

	//All players recieve poison mushroom
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)

	//Perform FFA Minigame
	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(0, &g)

	//Should recieve 2nd minigame as all players were poisoned
	EventIs(MinigameFFASelector{}, g.ExtraEvent, "Minigame", t)

	//Perform 2nd minigame
	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(0, &g)
	IntIs(2, int(g.Turn), "Turn #", t)
	EventIs(NormalDiceBlock{0}, g.ExtraEvent, "", t)
}

func TestSkipTurnViaCharacterTurn(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 7}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 7}),
		},
		CurrentPlayer: 2,
		ExtraEvent:    NormalDiceBlock{2},
		Config:        GameConfig{MaxTurns: 20},
	}
	g.Players[0].LastSpaceType = Blue
	g.Players[1].LastSpaceType = Blue

	//Player 2 fails mushroom check
	g.ExtraEvent.Handle(4, &g)
	g.ExtraEvent.Handle(false, &g)

	//Player 3 moves to blue space
	g.ExtraEvent.Handle(1, &g)

	//Handle Minigame
	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(0, &g)

	//Players 0 & 1 move to blue space
	g.ExtraEvent.Handle(1, &g)
	g.ExtraEvent.Handle(1, &g)

	EventIs(NormalDiceBlock{3}, g.ExtraEvent, "", t)
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
	g.MovePlayer(0, 4) //Land on happening
	g.ExtraEvent.Handle(BooStealAction{0, 1, false}, &g)
	expectedEvent := BooCoinsEvent{
		PayRangeEvent: PayRangeEvent{
			Player: 1,
			Min:    1,
			Max:    10, //Max of 10 coins
		},
		RecvPlayer: 0,
		Moves:      4,
	}
	EventIs(expectedEvent, g.ExtraEvent, "", t)

	expectedEvent.Handle(5, &g)
	CoinsIs(15, 0, g, "", t)
	CoinsIs(5, 1, g, "", t)
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
	g.MovePlayer(0, 4) //Want to land on happening
	g.ExtraEvent.Handle(BooStealAction{0, 1, false}, &g)
	expectedEvent := BooCoinsEvent{
		PayRangeEvent: PayRangeEvent{
			Player: 1,
			Min:    1,
			Max:    4, //Max of 4 coins
		},
		RecvPlayer: 0,
		Moves:      4,
	}
	EventIs(expectedEvent, g.ExtraEvent, "", t)

	expectedEvent.Handle(4, &g)
	CoinsIs(14, 0, g, "", t)
	CoinsIs(0, 1, g, "", t)
}

func TestPassEmptyBooSpace(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 4, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{
			NoBoo: true,
		},
	}
	g.MovePlayer(0, 4)
	EventIs(NormalDiceBlock{1}, g.ExtraEvent, "", t)
	SpaceIs(ChainSpace{1, 26}, 0, g, "", t)
}

func TestBuyStar(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 20, ChainSpace{0, 18}),
			NewPlayer("Luigi", 0, 4, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{
			MaxTurns: 25,
		},
	}

	g.MovePlayer(0, 1) //Land on blue space

	SpaceIs(ChainSpace{0, 20}, 0, g, "", t)
	CoinsIs(3, 0, g, "", t)
	StarsIs(1, 0, g, "", t)
}
