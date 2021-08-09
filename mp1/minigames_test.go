package mp1

import (
	"testing"
)

func Test4V4Minigame(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.Players[0].LastSpaceType = Blue
	g.Players[1].LastSpaceType = Blue
	g.Players[2].LastSpaceType = Blue
	g.Players[3].LastSpaceType = Blue

	g.GetMinigame()
	_, ok := g.ExtraEvent.(MinigameFFASelector)
	if !ok {
		t.Fatalf("Expected Minigame Type: MinigameFFASelector, got: %T",
			MinigameFFASelector{})
	}

	g.ExtraEvent.Handle(MinigameFFAMusicalMushroom, &g)
	g.ExtraEvent.Handle(0, &g) //Daisy wins
	expectedDaisyCoins := 20
	gotDaisyCoins := g.Players[0].Coins
	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Expected Daisy coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins)
	}
	expectedDaisyMinigameCoins := 10
	gotDaisyMinigameCoins := g.Players[0].MinigameCoins
	if expectedDaisyMinigameCoins != gotDaisyMinigameCoins {
		t.Errorf("Expected Daisy coins: %d, got: %d",
			expectedDaisyMinigameCoins, gotDaisyMinigameCoins)
	}
}

func Test1V3Minigame(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.Players[0].LastSpaceType = Blue
	g.Players[1].LastSpaceType = Blue
	g.Players[2].LastSpaceType = Blue
	g.Players[3].LastSpaceType = Red

	g.GetMinigame()
	minigame, ok := g.ExtraEvent.(Minigame1V3Selector)
	if !ok {
		t.Fatalf("Expected Minigame Type: Minigame1V3Selector, got: %T",
			Minigame1V3Selector{})
	}
	expectedSoloPlayer := 3
	gotSoloPlayer := minigame.Player
	if expectedSoloPlayer != gotSoloPlayer {
		t.Errorf("Expected solo player: %d, got: %d",
			expectedSoloPlayer,
			gotSoloPlayer,
		)
	}

	g.ExtraEvent.Handle(Minigame1V3TightropeTreachery, &g)
	g.ExtraEvent.Handle(0, &g) //Mario wins
	expectedMarioCoins := 25
	gotDaisyCoins := g.Players[3].Coins
	if expectedMarioCoins != gotDaisyCoins {
		t.Errorf("Expected Mario coins: %d, got: %d",
			expectedMarioCoins, gotDaisyCoins)
	}
	expectedMarioMinigameCoins := 15
	gotMarioMinigameCoins := g.Players[3].MinigameCoins
	if expectedMarioMinigameCoins != gotMarioMinigameCoins {
		t.Errorf("Expected Mario coins: %d, got: %d",
			expectedMarioMinigameCoins, gotMarioMinigameCoins)
	}
}

func Test2V2Minigame(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.Players[0].LastSpaceType = Blue
	g.Players[1].LastSpaceType = Red
	g.Players[2].LastSpaceType = Blue
	g.Players[3].LastSpaceType = Red

	g.GetMinigame()
	_, ok := g.ExtraEvent.(Minigame2V2Selector)
	if !ok {
		t.Fatalf("Expected Minigame Type: Minigame2V2Selector, got: %T",
			Minigame2V2Selector{})
	}

	g.ExtraEvent.Handle(Minigame2V2HandcarHavoc, &g)
	g.ExtraEvent.Handle(0, &g) //Daisy and DonkeyKong win
	expectedDaisyCoins := 20
	gotDaisyCoins := g.Players[0].Coins
	if expectedDaisyCoins != gotDaisyCoins {
		t.Errorf("Expected Daisy coins: %d, got: %d",
			expectedDaisyCoins, gotDaisyCoins)
	}
	expectedDaisyMinigameCoins := 10
	gotDaisyMinigameCoins := g.Players[0].MinigameCoins
	if expectedDaisyMinigameCoins != gotDaisyMinigameCoins {
		t.Errorf("Expected Daisy coins: %d, got: %d",
			expectedDaisyMinigameCoins, gotDaisyMinigameCoins)
	}
	expectedDKCoins := 20
	gotDKCoins := g.Players[2].Coins
	if expectedDKCoins != gotDKCoins {
		t.Errorf("Expected DK coins: %d, got: %d",
			expectedDKCoins, gotDKCoins)
	}
	expectedDKMinigameCoins := 10
	gotDKMinigameCoins := g.Players[2].MinigameCoins
	if expectedDKMinigameCoins != gotDKMinigameCoins {
		t.Errorf("Expected DK coins: %d, got: %d",
			expectedDKMinigameCoins, gotDKMinigameCoins)
	}
}

func TestGreenToBlue(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 21}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.Players[0].LastSpaceType = Blue
	g.Players[1].LastSpaceType = Blue
	g.Players[2].LastSpaceType = Happening
	g.Players[3].LastSpaceType = Blue

	g.FindGreenPlayer()
	expectedEvt := DeterminePlayerTeamEvent{2}
	gotEvt := g.ExtraEvent
	if expectedEvt != gotEvt {
		t.Errorf("Expected event: %#v, got: %#v",
			expectedEvt, gotEvt)
	}

	g.ExtraEvent.Handle(true, &g)
	expectedSpace := Blue
	gotSpace := g.Players[2].LastSpaceType
	if expectedSpace != gotSpace {
		t.Errorf("Expected Space Type: %d, got: %d",
			expectedSpace, gotSpace)
	}
}

func TestLandOnMinigameSpace(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 20}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{0, 0}),
		},
	}
	g.MovePlayer(0, 1)
	g.ExtraEvent.Handle(Minigame1PShellGame, &g)
	gLose := g
	gLose.ExtraEvent.Handle(-5, &gLose) //lose 5 coins
	expectedLoseCoins := 5
	gotLoseCoins := gLose.Players[0].Coins
	if expectedLoseCoins != gotLoseCoins {
		t.Errorf("Expected lose coins: %d, got: %d", expectedLoseCoins, gotLoseCoins)
	}

	gWin := g
	gWin.ExtraEvent.Handle(10, &gWin) //won WAP
	expectedWinCoins := 20
	gotWinCoins := gWin.Players[0].Coins
	if expectedWinCoins != gotWinCoins {
		t.Errorf("Expected win coins: %d, got: %d", expectedWinCoins, gotWinCoins)
	}
}

func TestPlayer4MinigameSpace(t *testing.T) {
	g := Game{
		Board: YTI,
		Players: [4]Player{
			NewPlayer("Daisy", 0, 10, ChainSpace{1, 20}),
			NewPlayer("Luigi", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Donkey Kong", 0, 10, ChainSpace{0, 0}),
			NewPlayer("Mario", 0, 10, ChainSpace{1, 20}),
		},
	}
	g.Players[0].LastSpaceType = Blue
	g.Players[1].LastSpaceType = Blue
	g.Players[2].LastSpaceType = Blue
	g.Players[3].LastSpaceType = Blue
	g.CurrentPlayer = 3

	g.MovePlayer(3, 1)
	g.ExtraEvent.Handle(Minigame1PShellGame, &g)
	g.ExtraEvent.Handle(-5, &g)
	expectedEvent := MinigameFFASelector{}
	gotEvent := g.ExtraEvent
	if expectedEvent != gotEvent {
		t.Errorf("Expected event: %#v, got: %#v", expectedEvent, gotEvent)
	}
}
