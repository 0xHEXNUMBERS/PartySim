package mp1

import "testing"

var ChanceBoard = MakeSimpleBoard(Chance)

func TestChanceTime(t *testing.T) {
	g := *InitializeGame(ChanceBoard, GameConfig{MaxTurns: 20})
	g.Players[0].Coins = 10
	g.Players[1].Coins = 25
	g.Players[0].Stars = 4
	g.Players[1].Stars = 9

	//Land on Chance
	g.MovePlayer(0, 1)
	//Steal from Player 1
	g.NextEvent.Handle(ChanceTimeResponse{CTBLeft, 1}, &g)
	//Player 0 will will receive
	g.NextEvent.Handle(ChanceTimeResponse{CTBRight, 0}, &g)

	g10Coins := g
	g10Coins.NextEvent.Handle(ChanceTimeResponse{CTBMiddle, int(LTR10)}, &g10Coins)
	CoinsIs(20, 0, g10Coins, "LTR10", t)
	CoinsIs(15, 1, g10Coins, "LTR10", t)

	g20Coins := g
	g20Coins.NextEvent.Handle(ChanceTimeResponse{CTBMiddle, int(LTR20)}, &g20Coins)
	CoinsIs(30, 0, g20Coins, "LTR20", t)
	CoinsIs(5, 1, g20Coins, "LTR20", t)

	g30Coins := g
	g30Coins.NextEvent.Handle(ChanceTimeResponse{CTBMiddle, int(LTR30)}, &g30Coins)
	CoinsIs(35, 0, g30Coins, "LTR30", t)
	CoinsIs(0, 1, g30Coins, "LTR30", t)

	gStar := g
	gStar.NextEvent.Handle(ChanceTimeResponse{CTBMiddle, int(LTRStar)}, &gStar)
	StarsIs(5, 0, gStar, "LTR", t)
	StarsIs(8, 1, gStar, "LTR", t)

	gSwapCoins := g
	gSwapCoins.NextEvent.Handle(ChanceTimeResponse{CTBMiddle, int(SwapCoins)}, &gSwapCoins)
	CoinsIs(25, 0, gSwapCoins, "Swap", t)
	CoinsIs(10, 1, gSwapCoins, "Swap", t)

	gSwapStars := g
	gSwapStars.NextEvent.Handle(ChanceTimeResponse{CTBMiddle, int(SwapStars)}, &gSwapStars)
	StarsIs(9, 0, gSwapStars, "Swap", t)
	StarsIs(4, 1, gSwapStars, "Swap", t)
	EventIs(NormalDiceBlock{Range{1, 10}, 1}, gSwapStars.NextEvent, "", t)
}
