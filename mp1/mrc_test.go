package mp1

import (
	"testing"
)

func TestMRCFork(t *testing.T) {
	g := *InitializeGame(MRC, GameConfig{MaxTurns: 25})
	g.Players[0].CurrentSpace = ChainSpace{0, 15}

	g.ExtraEvent.Handle(1, &g) //Move
	expectedRes := []Response{ChainSpace{1, 0}, ChainSpace{2, 0}}
	ResIs(expectedRes, g, "ChainSpaces", t)
}

func TestSwapCastleDirViaHappening(t *testing.T) {
	g := *InitializeGame(MRC, GameConfig{MaxTurns: 25})
	g.Players[0].CurrentSpace = ChainSpace{0, 6}

	g.ExtraEvent.Handle(1, &g) //Move
	bd := g.Board.Data.(mrcBoardData)
	if !bd.IsBowser {
		t.Errorf("Bowser is not swapped in")
	}
}

func TestSwapCastleDirViaStar(t *testing.T) {
	g := *InitializeGame(MRC, GameConfig{MaxTurns: 25})
	g.Players[0].CurrentSpace = ChainSpace{4, 7}

	g.ExtraEvent.Handle(1, &g) //Move
	bd := g.Board.Data.(mrcBoardData)
	if !bd.IsBowser {
		t.Errorf("Bowser is not swapped in")
	}

	CoinsIs(23, 0, g, "", t)
	SpaceIs(ChainSpace{0, 2}, 0, g, "", t)
}
