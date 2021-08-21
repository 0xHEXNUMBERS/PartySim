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
	CoinsIs(20, 0, g, "", t)
	MinigameCoinsIs(10, 0, g, "", t)
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
	IntIs(3, minigame.Player, "Solo Player", t)

	g.ExtraEvent.Handle(Minigame1V3TightropeTreachery, &g)
	g.ExtraEvent.Handle(0, &g) //Mario wins
	CoinsIs(25, 3, g, "", t)
	MinigameCoinsIs(15, 3, g, "", t)
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
	CoinsIs(20, 0, g, "", t)
	MinigameCoinsIs(10, 0, g, "", t)
	CoinsIs(20, 2, g, "", t)
	MinigameCoinsIs(10, 2, g, "", t)
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
	EventIs(DeterminePlayerTeamEvent{2}, g.ExtraEvent, "", t)

	g.ExtraEvent.Handle(true, &g)
	SpaceTypeIs(Blue, g.Players[2].LastSpaceType, "", t)
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
	CoinsIs(5, 0, gLose, "Lose", t)

	gWin := g
	gWin.ExtraEvent.Handle(10, &gWin) //won WAP
	CoinsIs(20, 0, gWin, "Win", t)
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
	EventIs(MinigameFFASelector{}, g.ExtraEvent, "", t)
}
