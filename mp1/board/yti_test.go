package board

import (
	"testing"

	"github.com/0xhexnumbers/partysim/mp1"
)

func TestMove(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 23)
	g.MovePlayer(0, 4)
	SpaceIs(mp1.NewChainSpace(1, 27), 0, g, "", t)
}

func TestCanPayThwomp(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 23)
	g.MovePlayer(0, 10)
	EventIs(YTIThwompBranchEvent{0, 6, 1}, g.NextEvent, "", t)
}

func TestCanNotPayThwomp(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 23)
	g.Players[0].Coins = 0
	g.MovePlayer(0, 10)
	EventIs(mp1.NormalDiceBlock{
		Range:  mp1.Range{Min: 1, Max: 10},
		Player: 1,
	}, g.NextEvent, "", t)
}

func TestGainCoins(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 23)
	g.MovePlayer(0, 1)
	CoinsIs(13, 0, g, "", t)
}

func TestPayThwompAndGainCoins(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 23)
	//Move player to invisible space
	g.MovePlayer(0, 10)
	//Accept payment to thwomp 1
	g.NextEvent.Handle(YTIThwompBranchPay, &g)
	//Pay thwomp 3 coins, move and land on blue space
	g.NextEvent.Handle(3, &g)
	EventIs(mp1.NormalDiceBlock{
		Range:  mp1.Range{Min: 1, Max: 10},
		Player: 1,
	}, g.NextEvent, "", t)
	SpaceIs(mp1.NewChainSpace(0, 12), 0, g, "", t)
	CoinsIs(10, 0, g, "", t)
}

func TestIgnoreThwomp(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 23)
	g.MovePlayer(0, 10)
	g.NextEvent.Handle(YTIThwompBranchIgnore, &g)
	EventIs(mp1.NormalDiceBlock{Range: mp1.Range{Min: 1, Max: 10}, Player: 1}, g.NextEvent, "", t)
	SpaceIs(mp1.NewChainSpace(1, 5), 0, g, "", t)
	CoinsIs(13, 0, g, "", t)
}

func TestStarSwapViaHappening(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 23)
	g.MovePlayer(0, 3)
	EventIs(mp1.NormalDiceBlock{
		Range:  mp1.Range{Min: 1, Max: 10},
		Player: 1,
	}, g.NextEvent, "", t)

	bd := g.Board.Data.(ytiBoardData)
	if bd.StarPosition != ytiRightIslandStar {
		t.Errorf("Expected star position: %#v, got: %#v", ytiRightIslandStar, bd.StarPosition)
	}
}

func TestCoinsOnStart(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 22)
	g.MovePlayer(0, 1)
	CoinsIs(20, 0, g, "", t)
}

func TestMushroomSpace(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 7)
	g.MovePlayer(0, 4)
	EventIs(mp1.MushroomEvent{Player: 0}, g.NextEvent, "", t)

	//Received red mushroom
	gRed := g
	gRed.NextEvent.Handle(mp1.RedMushroom, &gRed)
	EventIs(mp1.NormalDiceBlock{
		Range:  mp1.Range{Min: 1, Max: 10},
		Player: 0,
	}, gRed.NextEvent, "", t)

	//Received poison mushroom
	gPoison := g
	gPoison.NextEvent.Handle(mp1.PoisonMushroom, &gPoison)
	EventIs(mp1.NormalDiceBlock{
		Range:  mp1.Range{Min: 1, Max: 10},
		Player: 1,
	}, gPoison.NextEvent, "", t)
	if !gPoison.Players[0].SkipTurn {
		t.Errorf("SkipTurn not set")
	}
}

func TestSkipTurnViaMinigame(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 7)
	g.Players[1].CurrentSpace = mp1.NewChainSpace(0, 7)
	g.Players[2].CurrentSpace = mp1.NewChainSpace(0, 7)
	g.Players[3].CurrentSpace = mp1.NewChainSpace(0, 7)

	//All players receive poison mushroom
	g.NextEvent.Handle(4, &g)
	g.NextEvent.Handle(mp1.PoisonMushroom, &g)
	g.NextEvent.Handle(4, &g)
	g.NextEvent.Handle(mp1.PoisonMushroom, &g)
	g.NextEvent.Handle(4, &g)
	g.NextEvent.Handle(mp1.PoisonMushroom, &g)
	g.NextEvent.Handle(4, &g)
	g.NextEvent.Handle(mp1.PoisonMushroom, &g)

	//Perform FFA Minigame
	g.NextEvent.Handle(mp1.MinigameFFAMusicalMushroom, &g)
	g.NextEvent.Handle(0, &g)

	//Should receive 2nd minigame as all players were poisoned
	EventIs(mp1.MinigameFFASelector{}, g.NextEvent, "Minigame", t)

	//Perform 2nd minigame
	g.NextEvent.Handle(mp1.MinigameFFAMusicalMushroom, &g)
	g.NextEvent.Handle(0, &g)
	IntIs(2, int(g.Turn), "Turn #", t)
	EventIs(mp1.NormalDiceBlock{
		Range:  mp1.Range{Min: 1, Max: 10},
		Player: 0,
	}, g.NextEvent, "", t)
}

func TestSkipTurnViaCharacterTurn(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 7)
	g.Players[1].CurrentSpace = mp1.NewChainSpace(0, 7)
	g.Players[2].CurrentSpace = mp1.NewChainSpace(0, 7)
	g.Players[3].CurrentSpace = mp1.NewChainSpace(0, 7)
	g.CurrentPlayer = 2
	g.NextEvent = mp1.NormalDiceBlock{
		Range:  mp1.Range{Min: 1, Max: 10},
		Player: 2,
	}
	g.Players[0].LastSpaceType = mp1.Blue
	g.Players[1].LastSpaceType = mp1.Blue

	//Player 2 fails mushroom check
	g.NextEvent.Handle(4, &g)
	g.NextEvent.Handle(mp1.PoisonMushroom, &g)

	//Player 3 moves to blue space
	g.NextEvent.Handle(1, &g)

	//Handle Minigame
	g.NextEvent.Handle(mp1.MinigameFFAMusicalMushroom, &g)
	g.NextEvent.Handle(0, &g)

	//Players 0 & 1 move to blue space
	g.NextEvent.Handle(1, &g)
	g.NextEvent.Handle(1, &g)

	EventIs(mp1.NormalDiceBlock{
		Range:  mp1.Range{Min: 1, Max: 10},
		Player: 3,
	}, g.NextEvent, "", t)
}

func TestStealCoinsViaBoo(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 21)

	g.MovePlayer(0, 4) //Land on happening
	g.NextEvent.Handle(mp1.BooStealAction{
		RecvPlayer: 0, GivingPlayer: 1, Star: false,
	}, &g)
	expectedEvent := mp1.BooCoinsEvent{
		PayRangeEvent: mp1.PayRangeEvent{
			Range:  mp1.Range{Min: 1, Max: 10},
			Player: 1,
		},
		RecvPlayer: 0,
		Moves:      4,
	}
	EventIs(expectedEvent, g.NextEvent, "", t)

	expectedEvent.Handle(5, &g)
	CoinsIs(15, 0, g, "", t)
	CoinsIs(5, 1, g, "", t)
}

func TestStealTooManyCoinsViaBoo(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 21)
	g.Players[1].Coins = 4
	g.MovePlayer(0, 4) //Want to land on happening
	g.NextEvent.Handle(mp1.BooStealAction{
		RecvPlayer: 0, GivingPlayer: 1, Star: false,
	}, &g)
	expectedEvent := mp1.BooCoinsEvent{
		PayRangeEvent: mp1.PayRangeEvent{
			Range:  mp1.Range{Min: 1, Max: 4},
			Player: 1,
		},
		RecvPlayer: 0,
		Moves:      4,
	}
	EventIs(expectedEvent, g.NextEvent, "", t)

	expectedEvent.Handle(4, &g)
	CoinsIs(14, 0, g, "", t)
	CoinsIs(0, 1, g, "", t)
}

func TestPassEmptyBooSpace(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20, NoBoo: true})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(1, 21)
	g.MovePlayer(0, 4)
	EventIs(mp1.NormalDiceBlock{
		Range:  mp1.Range{Min: 1, Max: 10},
		Player: 1,
	}, g.NextEvent, "", t)
	SpaceIs(mp1.NewChainSpace(1, 26), 0, g, "", t)
}

func TestBuyStar(t *testing.T) {
	g := *mp1.InitializeGame(YTI, mp1.GameConfig{MaxTurns: 20})
	g.Players[0].CurrentSpace = mp1.NewChainSpace(0, 18)
	g.Players[0].Coins = 20
	g.MovePlayer(0, 1) //Land on blue space
	SpaceIs(mp1.NewChainSpace(0, 20), 0, g, "", t)
	CoinsIs(3, 0, g, "", t)
	StarsIs(1, 0, g, "", t)
}
