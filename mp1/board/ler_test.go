package board

import (
	"testing"

	"github.com/0xhexnumbers/partysim/mp1"
)

func TestRBRFork(t *testing.T) {
	g := *mp1.InitializeGame(LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 4)
	g.NextEvent.Handle(uint8(0), &g) //Star
	gBlueUpLeft := g
	gBlueUpLeft.Board.Data = lerBoardData{true}
	gBlueUpRight := gBlueUpLeft

	g.NextEvent.Handle(1, &g)
	SpaceIs(mp1.NewChainSpace(5, 0), 0, g, "", t)

	gBlueUpLeft.NextEvent.Handle(1, &gBlueUpLeft)
	gBlueUpLeft.NextEvent.Handle(
		mp1.NewChainSpace(3, 0), &gBlueUpLeft,
	)
	SpaceIs(mp1.NewChainSpace(3, 0), 0, gBlueUpLeft, "BlueUpLeft", t)

	gBlueUpRight.NextEvent.Handle(1, &gBlueUpRight)
	gBlueUpRight.NextEvent.Handle(
		mp1.NewChainSpace(11, 0), &gBlueUpRight,
	)
	SpaceIs(mp1.NewChainSpace(11, 0), 0, gBlueUpRight, "BlueUpRight", t)
}

func TestRBFork(t *testing.T) {
	g := *mp1.InitializeGame(LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(3, 8)
	g.NextEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.NextEvent.Handle(1, &g)
	SpaceIs(mp1.NewChainSpace(4, 0), 0, g, "", t)

	gBlueUp.NextEvent.Handle(1, &gBlueUp)
	SpaceIs(mp1.NewChainSpace(4, 4), 0, gBlueUp, "BlueUp", t)
}

func TestBRFork1(t *testing.T) {
	g := *mp1.InitializeGame(LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(5, 3)
	g.NextEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.NextEvent.Handle(1, &g)
	SpaceIs(mp1.NewChainSpace(6, 10), 0, g, "", t)

	gBlueUp.NextEvent.Handle(1, &gBlueUp)
	SpaceIs(mp1.NewChainSpace(9, 0), 0, gBlueUp, "BlueUp", t)
}

func TestBRFork2(t *testing.T) {
	g := *mp1.InitializeGame(LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(6, 12)
	g.NextEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.NextEvent.Handle(1, &g)
	SpaceIs(mp1.NewChainSpace(6, 0), 0, g, "", t)

	gBlueUp.NextEvent.Handle(1, &gBlueUp)
	SpaceIs(mp1.NewChainSpace(7, 0), 0, gBlueUp, "BlueUp", t)
}

func TestBRFork3(t *testing.T) {
	g := *mp1.InitializeGame(LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(9, 8)
	g.NextEvent.Handle(uint8(0), &g) //Star
	gBlueUp := g
	gBlueUp.Board.Data = lerBoardData{true}

	g.NextEvent.Handle(1, &g)
	SpaceIs(mp1.NewChainSpace(10, 0), 0, g, "", t)

	gBlueUp.NextEvent.Handle(1, &gBlueUp)
	SpaceIs(mp1.NewChainSpace(0, 0), 0, gBlueUp, "BlueUp", t)
}

func TestSwapGatesViaHappening(t *testing.T) {
	g := *mp1.InitializeGame(LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(3, 7)
	g.Players[0].Coins = 20
	g.NextEvent.Handle(uint8(0), &g) //Star

	g.NextEvent.Handle(1, &g) //Move
	bd := g.Board.Data.(lerBoardData)
	if !bd.BlueUp {
		t.Errorf("Gates did not swap via Happening")
	}
}

func TestSwapGatesTwice(t *testing.T) {
	g := *mp1.InitializeGame(LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(11, 5)
	g.Players[0].Coins = 20
	g.NextEvent.Handle(uint8(0), &g) //Star

	g.NextEvent.Handle(1, &g)    //Move
	g.NextEvent.Handle(true, &g) //Swap Gates
	bd := g.Board.Data.(lerBoardData)
	if bd.BlueUp {
		t.Errorf("Gates did not swap twice")
	}
}

func TestSwapGatesInsufficientCoins(t *testing.T) {
	g := *mp1.InitializeGame(LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(2, 3)
	g.Players[0].Coins = 0
	g.NextEvent.Handle(uint8(0), &g) //Star

	g.NextEvent.Handle(1, &g)

	SpaceIs(mp1.NewChainSpace(2, 5), 0, g, "", t)
}

func TestNormalBranch(t *testing.T) {
	g := *mp1.InitializeGame(LER, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 10)
	g.NextEvent.Handle(uint8(0), &g)

	g.NextEvent.Handle(1, &g)                       //Move
	g.NextEvent.Handle(mp1.NewChainSpace(1, 0), &g) //Branch

	SpaceIs(mp1.NewChainSpace(1, 0), 0, g, "", t)
}
