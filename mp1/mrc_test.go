package mp1

import (
	"reflect"
	"testing"
)

func TestMRCFork(t *testing.T) {
	g := *InitializeGame(MRC, GameConfig{MaxTurns: 25})
	g.Players[0].CurrentSpace = ChainSpace{0, 15}

	g.ExtraEvent.Handle(1, &g) //Move
	expectedRes := []Response{ChainSpace{1, 0}, ChainSpace{2, 0}}
	gotRes := g.ExtraEvent.Responses()
	if !reflect.DeepEqual(expectedRes, gotRes) {
		t.Errorf("Expected ChainSpaces: %#v, got: %#v",
			expectedRes, gotRes)
	}
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

	expectedCoins := 23
	gotCoins := g.Players[0].Coins
	if expectedCoins != gotCoins {
		t.Errorf("Expected coins: %d, got: %d",
			expectedCoins, gotCoins)
	}
}
