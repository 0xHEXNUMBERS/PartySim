package mp1

import "testing"

func TestWhompPayment(t *testing.T) {
	g := *InitializeGame(DKJA, GameConfig{MaxTurns: 20})
	g.NextEvent.Handle(uint8(0), &g) //Star
	g.NextEvent.Handle(5, &g)        //Move

	gPay := g
	gPay.NextEvent.Handle(true, &gPay) //Pay Whomp
	SpaceIs(ChainSpace{4, 0}, 0, gPay, "Pay", t)
	CoinsIs(3, 0, gPay, "Pay", t)

	gSkip := g
	gSkip.NextEvent.Handle(false, &gSkip) //Ignore Whomp
	SpaceIs(ChainSpace{1, 0}, 0, gSkip, "Skip", t)
	CoinsIs(13, 0, gSkip, "Skip", t)

	gIgnore := gSkip
	gIgnore.Players[1].Coins = 0
	gIgnore.NextEvent.Handle(5, &gIgnore) //Move
	SpaceIs(ChainSpace{1, 0}, 1, gIgnore, "Ignore", t)
	CoinsIs(3, 1, gIgnore, "Ignore", t)
}

func TestCoinBlockade(t *testing.T) {
	g := *InitializeGame(DKJA, GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = ChainSpace{1, 3}
	g.Players[0].Coins = 20
	g.NextEvent.Handle(uint8(0), &g) //Star

	gPass := g
	gPass.NextEvent.Handle(1, &gPass)    //Move
	gPass.NextEvent.Handle(true, &gPass) //Pass
	SpaceIs(ChainSpace{2, 0}, 0, gPass, "Pass", t)
	CoinsIs(23, 0, gPass, "Pass", t)

	gSkip := g
	gSkip.NextEvent.Handle(1, &gSkip)     //Move
	gSkip.NextEvent.Handle(false, &gSkip) //Skip
	SpaceIs(ChainSpace{3, 0}, 0, gSkip, "Skip", t)
	CoinsIs(23, 0, gSkip, "Skip", t)

	gIgnore := g
	gIgnore.Players[0].Coins = 0
	gIgnore.NextEvent.Handle(1, &gIgnore) //Move
	SpaceIs(ChainSpace{3, 0}, 0, gIgnore, "Ignore", t)
	CoinsIs(3, 0, gIgnore, "Ignore", t)
}

func TestBoulder(t *testing.T) {
	g := *InitializeGame(DKJA, GameConfig{MaxTurns: 20})
	g.NextEvent.Handle(uint8(0), &g) //Star
	g.Players[0].CurrentSpace = ChainSpace{5, 1}
	g.Players[1].CurrentSpace = ChainSpace{5, 0}
	g.Players[2].CurrentSpace = ChainSpace{7, 0}
	g.Players[3].CurrentSpace = ChainSpace{5, 5}

	g.NextEvent.Handle(2, &g)
	SpaceIs(ChainSpace{0, 16}, 0, g, "", t)
	SpaceIs(ChainSpace{5, 0}, 1, g, "", t)
	SpaceIs(ChainSpace{0, 16}, 2, g, "", t)
	SpaceIs(ChainSpace{0, 16}, 3, g, "", t)
}
