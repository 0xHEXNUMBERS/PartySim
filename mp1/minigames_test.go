package mp1

import (
	"reflect"
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

	g = GetMinigame(g)
	minigame := g.ExtraEvent.(MinigameEvent)
	if minigame.Type != MinigameFFA {
		t.Fatalf("Expected Minigame Type: %d, got: %d",
			MinigameFFA, minigame.Type)
	}
	rewards := minigame.Responses()
	if !reflect.DeepEqual(rewards, MinigameRewardsFFA) {
		t.Fatal("Recieved incorrect minigame awards")
	}

	g = minigame.Handle(rewards[0], g) //Daisy wins
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

	g = GetMinigame(g)
	minigame := g.ExtraEvent.(MinigameEvent)
	if minigame.Type != Minigame1V3 {
		t.Fatalf("Expected Minigame Type: %d, got: %d",
			MinigameFFA, minigame.Type)
	}
	rewards := minigame.Responses()
	if !reflect.DeepEqual(rewards, MinigameRewards1V3) {
		t.Fatal("Recieved incorrect minigame awards")
	}

	g = minigame.Handle(rewards[0], g) //Mario wins
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

	g = GetMinigame(g)
	minigame := g.ExtraEvent.(MinigameEvent)
	if minigame.Type != Minigame2V2 {
		t.Fatalf("Expected Minigame Type: %d, got: %d",
			MinigameFFA, minigame.Type)
	}
	rewards := minigame.Responses()
	if !reflect.DeepEqual(rewards, MinigameRewards2V2) {
		t.Fatal("Recieved incorrect minigame awards")
	}

	g = minigame.Handle(rewards[0], g) //Daisy and DonkeyKong win
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

	g = FindGreenPlayer(g)
	expectedEvt := DeterminePlayerTeamEvent{2}
	gotEvt := g.ExtraEvent
	if expectedEvt != gotEvt {
		t.Errorf("Expected event: %#v, got: %#v",
			expectedEvt, gotEvt)
	}

	g = g.ExtraEvent.Handle(true, g)
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
	g = MovePlayer(g, 0, 1)
	gLose := g.ExtraEvent.Handle(MinigameRewards1P[0], g) //lose 5 coins
	expectedLoseCoins := 5
	gotLoseCoins := gLose.Players[0].Coins
	if expectedLoseCoins != gotLoseCoins {
		t.Errorf("Expected lose coins: %d, got: %d", expectedLoseCoins, gotLoseCoins)
	}

	gWin := g.ExtraEvent.Handle(MinigameRewards1P[11], g) //won WAP
	expectedWinCoins := 46
	gotWinCoins := gWin.Players[0].Coins
	if expectedWinCoins != gotWinCoins {
		t.Errorf("Expected win coins: %d, got: %d", expectedWinCoins, gotWinCoins)
	}
}
