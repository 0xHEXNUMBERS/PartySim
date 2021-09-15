package board

import (
	"testing"

	"github.com/0xhexnumbers/partysim/mp1"
)

func TestMRCFork(t *testing.T) {
	g := *mp1.InitializeGame(MRC, mp1.GameConfig{MaxTurns: 25})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 15)

	g.NextEvent.Handle(1, &g) //Move
	expectedRes := []mp1.Response{mp1.NewChainSpace(1, 0), mp1.NewChainSpace(2, 0)}
	ResIs(expectedRes, g, "mp1.ChainSpaces", t)
}

func TestSwapCastleDirViaHappening(t *testing.T) {
	g := *mp1.InitializeGame(MRC, mp1.GameConfig{MaxTurns: 25})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 6)

	g.NextEvent.Handle(1, &g) //Move
	bd := g.Board.Data.(mrcBoardData)
	if !bd.IsBowser {
		t.Errorf("Bowser is not swapped in")
	}
}

func TestSwapCastleDirViaStar(t *testing.T) {
	g := *mp1.InitializeGame(MRC, mp1.GameConfig{MaxTurns: 25})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(4, 7)

	g.NextEvent.Handle(1, &g) //Move
	bd := g.Board.Data.(mrcBoardData)
	if !bd.IsBowser {
		t.Errorf("Bowser is not swapped in")
	}

	CoinsIs(23, 0, g, "", t)
	SpaceIs(mp1.NewChainSpace(0, 2), 0, g, "", t)
}
