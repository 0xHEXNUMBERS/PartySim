package mp1

import (
	"testing"
)

var MinigameBoard = MakeRepeatedBoard(MinigameSpace, 15)

func TestRedDiceBlock(t *testing.T) {
	g := *InitializeGame(MinigameBoard, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{}
	g.NextEvent = RedDiceBlock{Range{1, 10}, 0}
	g.NextEvent.Handle(9, &g) //Land on minigame space
	SpaceIs(ChainSpace{0, 9}, 0, g, "", t)
	CoinsIs(1, 0, g, "", t)
}

func TestBlueDiceBlock(t *testing.T) {
	g := *InitializeGame(MinigameBoard, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{}
	g.NextEvent = BlueDiceBlock{Range{1, 10}, 0}
	g.NextEvent.Handle(9, &g) //Land on minigame space
	SpaceIs(ChainSpace{0, 9}, 0, g, "", t)
	CoinsIs(19, 0, g, "", t)
}

func TestWarpDiceBlock(t *testing.T) {
	g := *InitializeGame(MinigameBoard, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{0, 5}
	g.Players[1].CurrentSpace = ChainSpace{}
	g.NextEvent = WarpDiceBlock{0}
	g.NextEvent.Handle(1, &g) //Swap with Luigi
	SpaceIs(ChainSpace{0, 0}, 0, g, "", t)
	SpaceIs(ChainSpace{0, 5}, 1, g, "", t)
}

func TestEventDiceBlock(t *testing.T) {
	g := *InitializeGame(MinigameBoard, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{0, 5}
	g.NextEvent = EventDiceBlock{0}
	gBoo := g
	gBoo.NextEvent.Handle(BooEventBlock, &gBoo)
	expectedBooEvent := BooEvent{0, gBoo.Players, 0, gBoo.Players[0].Coins}
	EventIs(expectedBooEvent, gBoo.NextEvent, "Boo", t)

	gBoo.NextEvent.Handle(BooStealAction{0, 1, false}, &gBoo)
	gBoo.NextEvent.Handle(10, &gBoo)
	SpaceIs(ChainSpace{0, 5}, 0, gBoo, "Boo", t)

	gBowser := g
	gBowser.NextEvent.Handle(BowserEventBlock, &gBowser)
	EventIs(NormalDiceBlock{Range{1, 10}, 1}, gBowser.NextEvent, "Bowser", t)
	CoinsIs(0, 0, gBowser, "Bowser", t)

	gKoopa := g
	gKoopa.NextEvent.Handle(KoopaEventBlock, &gKoopa)
	EventIs(NormalDiceBlock{Range{1, 10}, 1}, gKoopa.NextEvent, "Koopa", t)
	CoinsIs(20, 0, gKoopa, "Koopa", t)
}

func TestPickDiceBlock(t *testing.T) {
	g := *InitializeGame(MinigameBoard, GameConfig{MaxTurns: 20, RedDice: true, BlueDice: true})

	//Engine always sets first diceblock for all players to normal dice
	//block. We must tell the engine we're not on the first turn to test
	//hidden block behavior.
	g.Turn = 1
	g.SetDiceBlock()

	expected := []Response{
		NormalDiceBlock{Range{1, 10}, 0},
		RedDiceBlock{Range{1, 10}, 0},
		BlueDiceBlock{Range{1, 10}, 0},
	}
	ResIs(expected, g, "", t)
}
