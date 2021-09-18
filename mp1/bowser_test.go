package mp1

import (
	"testing"
)

var BowserBoard = MakeSimpleBoard(Bowser)

func TestBowser10CoinsForStar(t *testing.T) {
	g := *InitializeGame(BowserBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 0
	g.Players[0].Stars = 1

	g.MovePlayer(0, 1)
	StarsIs(0, 0, g, "", t)
	CoinsIs(10, 0, g, "", t)
}

func TestBowserGain20Coins(t *testing.T) {
	g := *InitializeGame(BowserBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 0

	g.MovePlayer(0, 1)
	CoinsIs(20, 0, g, "", t)
}

func TestCoinsForBowser(t *testing.T) {
	g := *InitializeGame(BowserBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 25

	g.MovePlayer(0, 1)
	g.NextEvent.Handle(CoinsForBowser, &g)
	EventIs(NormalDiceBlock{Range{1, 10}, 1}, g.NextEvent, "", t)
	CoinsIs(15, 0, g, "", t)
}

func TestBowserBalloonBurst(t *testing.T) {
	g := *InitializeGame(BowserBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 50
	g.Players[1].Coins = 50
	g.Players[2].Coins = 50
	g.Players[3].Coins = 50

	g.MovePlayer(0, 1)
	g.NextEvent.Handle(BowserBalloonBurst, &g)
	gDraw := g
	gDraw.NextEvent.Handle(4, &gDraw)
	for i := range gDraw.Players {
		CoinsIs(30, i, gDraw, "Draw", t)
	}

	gP1Win := g
	gP1Win.NextEvent.Handle(0, &gP1Win)
	CoinsIs(50, 0, gP1Win, "P1Win", t)
	for i := 1; i < 4; i++ {
		CoinsIs(40, i, gP1Win, "P1Win", t)
	}
}

func TestBowsersFaceLift(t *testing.T) {
	g := *InitializeGame(BowserBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 50
	g.Players[1].Coins = 50
	g.Players[2].Coins = 50
	g.Players[3].Coins = 50

	g.MovePlayer(0, 1)
	g.NextEvent.Handle(BowsersFaceLift, &g)
	gDraw := g
	gDraw.NextEvent.Handle(0b1111, &gDraw) //Draw
	CoinsIs(0, 0, gDraw, "Draw", t)

	gP1Loss := g
	gP1Loss.NextEvent.Handle(0b1110, &gP1Loss) //All players except Daisy
	CoinsIs(40, 0, gP1Loss, "P1Loss", t)
}

func TestBowsersTugoWar(t *testing.T) {
	g := *InitializeGame(BowserBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 50
	g.Players[1].Coins = 50
	g.Players[2].Coins = 50
	g.Players[3].Coins = 50
	g.MovePlayer(0, 1)
	g.NextEvent.Handle(BowsersTugoWar, &g)
	gDraw := g
	gDraw.NextEvent.Handle(BTWDraw, &gDraw)
	for i := range gDraw.Players {
		CoinsIs(20, i, gDraw, "Draw", t)
	}

	g1TWin := g
	g1TWin.NextEvent.Handle(BTW1TWin, &g1TWin)
	for i := 1; i < 4; i++ {
		CoinsIs(40, i, g1TWin, "1TWin", t)
	}

	g3TWin := g
	g3TWin.NextEvent.Handle(BTW3TWin, &g3TWin)
	CoinsIs(40, 0, g3TWin, "3TWin", t)
}

func TestBashnCash(t *testing.T) {
	g := *InitializeGame(BowserBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 54
	g.Players[1].Coins = 50
	g.Players[2].Coins = 50
	g.Players[3].Coins = 50
	g.MovePlayer(0, 1)
	g.NextEvent.Handle(BashnCash, &g)
	gGE5 := g
	gGE5.NextEvent.Handle(5, &gGE5) //Should lose 5 * 5 = 25 coins
	CoinsIs(29, 0, gGE5, "GE5", t)

	gE5 := g
	gE5.NextEvent.Handle(10, &gE5) //Should lose 5 * 10 = 50 coins
	CoinsIs(4, 0, gE5, "E5", t)

	gLT5 := g
	gLT5.NextEvent.Handle(13, &gLT5) //Should lose 5 * 10 + 3 = 53 coins
	CoinsIs(1, 0, gLT5, "LT5", t)
}

func TestBowserRevolution(t *testing.T) {
	g := *InitializeGame(BowserBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 75
	g.Players[1].Coins = 25
	g.Players[2].Coins = 99
	g.Players[3].Coins = 50
	g.MovePlayer(0, 1)
	g.NextEvent.Handle(BowserRevolution, &g)
	for i := range g.Players {
		CoinsIs(62, i, g, "", t)
	}
}

func TestBowsersChanceTime(t *testing.T) {
	g := *InitializeGame(BowserBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 50
	g.Players[1].Coins = 50
	g.Players[2].Coins = 50
	g.Players[3].Coins = 50
	g.MovePlayer(0, 1)
	g.NextEvent.Handle(BowsersChanceTime, &g)
	g.NextEvent.Handle(BCTResponse{0, 20}, &g) //Daisy
	CoinsIs(30, 0, g, "", t)
}
