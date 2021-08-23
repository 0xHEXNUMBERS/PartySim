package mp1

import (
	"testing"
)

func TestRedDiceBlock(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		NextEvent: RedDiceBlock{0},
	}
	g.NextEvent.Handle(9, &g) //Land on minigame space
	SpaceIs(ChainSpace{0, 9}, 0, g, "", t)
	CoinsIs(1, 0, g, "", t)
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
		NextEvent: BlueDiceBlock{0},
	}
	g.NextEvent.Handle(9, &g) //Land on minigame space
	SpaceIs(ChainSpace{0, 9}, 0, g, "", t)
	CoinsIs(19, 0, g, "", t)
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
		NextEvent: WarpDiceBlock{0},
	}
	g.NextEvent.Handle(1, &g) //Swap with Luigi
	SpaceIs(ChainSpace{0, 0}, 0, g, "", t)
	SpaceIs(ChainSpace{1, 23}, 1, g, "", t)
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
		NextEvent: EventDiceBlock{0},
	}
	gBoo := g
	gBoo.NextEvent.Handle(BooEventBlock, &gBoo)
	expectedBooEvent := BooEvent{0, gBoo.Players, 0, gBoo.Players[0].Coins}
	EventIs(expectedBooEvent, gBoo.NextEvent, "Boo", t)

	gBoo.NextEvent.Handle(BooStealAction{0, 1, false}, &gBoo)
	gBoo.NextEvent.Handle(10, &gBoo)
	SpaceIs(ChainSpace{1, 23}, 0, gBoo, "Boo", t)

	gBowser := g
	gBowser.NextEvent.Handle(BowserEventBlock, &gBowser)
	EventIs(NormalDiceBlock{1}, gBowser.NextEvent, "Bowser", t)
	CoinsIs(0, 0, gBowser, "Bowser", t)

	gKoopa := g
	gKoopa.NextEvent.Handle(KoopaEventBlock, &gKoopa)
	EventIs(NormalDiceBlock{1}, gKoopa.NextEvent, "Koopa", t)
	CoinsIs(20, 0, gKoopa, "Koopa", t)
}

func TestPickDiceBlock(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 23}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
		Config: GameConfig{
			RedDice:  true,
			BlueDice: true,
		},
	}
	g.NextEvent = PickDiceBlock{0, g.Config}
	expected := []Response{NormalDiceBlock{0}, RedDiceBlock{0}, BlueDiceBlock{0}}
	ResIs(expected, g, "", t)
}
