package mp1

import "testing"

func TestChanceTime(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 4, 10, ChainSpace{0, 14}),
			NewPlayer("Luigi", 9, 25, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	//Steal from Luigi
	g.NextEvent.Handle(ChanceTimeResponse{CTBLeft, 1}, &g)
	//Daisy will recieve
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
	EventIs(NormalDiceBlock{1}, gSwapStars.NextEvent, "", t)
}
