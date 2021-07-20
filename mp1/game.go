package mp1

type GameConfig struct {
	MaxTurns     uint8
	NoBonusStars bool
	NoKoopa      bool
	NoBoo        bool
	RedDice      bool
	BlueDice     bool
	WarpDice     bool
	EventsDice   bool
}

type Game struct {
	Board
	Players       [4]Player
	Turn          uint8
	CurrentPlayer int
	ExtraEvent    Event
	Config        GameConfig
}

func LastFiveTurns(g Game) bool {
	return g.Config.MaxTurns-g.Turn <= 5
}

func AwardCoins(g Game, player, coins int, minigame bool) Game {
	g.Players[player].Coins += coins
	if g.Players[player].Coins < 0 {
		g.Players[player].Coins = 0
	}
	if minigame {
		g.Players[player].MinigameCoins += coins
		if g.Players[player].MinigameCoins < 0 {
			g.Players[player].MinigameCoins = 0
		}
	}
	if g.Players[player].Coins > g.Players[player].MaxCoins {
		g.Players[player].MaxCoins = g.Players[player].Coins
	}
	return g
}

func MovePlayer(g Game, playerIdx, moves int) Game {
	chains := *g.Board.Chains
	playerPos := g.Players[playerIdx].CurrentSpace
	for moves > 0 {
		playerPos.Space++
		g.Players[playerIdx].CurrentSpace = playerPos
		if playerPos.Space >= len(chains[playerPos.Chain]) {
			playerPos.Space = 0
		}
		curSpace := chains[playerPos.Chain][playerPos.Space]
		switch curSpace.Type {
		case Invisible:
			g.ExtraEvent = nil
			g = curSpace.PassingEvent(g, playerIdx, moves)
			if g.ExtraEvent != nil {
				g.Players[playerIdx].CurrentSpace = playerPos
				return g
			}
		case Start:
			if !g.Config.NoKoopa {
				g = AwardCoins(g, playerIdx, 10, false)
			}
		case Star:
			g.ExtraEvent = nil
			g = curSpace.PassingEvent(g, playerIdx, moves)
			if g.ExtraEvent != nil {
				return g
			}
		case Boo:
			if !g.Config.NoBoo {
				booEvt := BooEvent{
					playerIdx,
					g.Players,
					moves,
					g.Players[playerIdx].Coins,
				}
				if len(booEvt.Responses()) != 0 {
					g.ExtraEvent = booEvt
					return g
				}
			}
		default:
			moves--
		}
	}
	//Stop on Space
	g.Players[playerIdx].CurrentSpace = playerPos
	//Activate Space
	curSpace := chains[playerPos.Chain][playerPos.Space]
	g.Players[playerIdx].LastSpaceType = curSpace.Type
	switch curSpace.Type {
	case Blue:
		g = AwardCoins(g, playerIdx, 3, false)
		if LastFiveTurns(g) {
			g = AwardCoins(g, playerIdx, 3, false)
		}
		g = EndCharacterTurn(g)
	case Red:
		g = AwardCoins(g, playerIdx, -3, false)
		if LastFiveTurns(g) {
			g = AwardCoins(g, playerIdx, -3, false)
		}
		g = EndCharacterTurn(g)
	case Mushroom:
		g.ExtraEvent = MushroomEvent{playerIdx}
	case Happening:
		g.Players[playerIdx].HappeningCount++
		g.ExtraEvent = nil
		g = curSpace.StoppingEvent(g)
		if g.ExtraEvent == nil {
			g = EndCharacterTurn(g)
		}
	case Bowser:
		g = PreBowserCheck(g, playerIdx)
	case MinigameSpace:
		g.ExtraEvent = MinigameEvent{[4]int{playerIdx, 0, 0, 0}, Minigame1P}
	case Chance:
		g.ExtraEvent = ChanceTime{Player: playerIdx}
	}
	return g
}

func AwardBonusStars(g Game) Game {
	if g.Config.NoBonusStars {
		return g
	}

	maxCoins := g.Players[0].MaxCoins
	for i := 1; i < 4; i++ {
		maxCoins = max(maxCoins, g.Players[i].MaxCoins)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].MaxCoins == maxCoins {
			g.Players[i].Stars++
		}
	}

	maxMinigameCoins := g.Players[0].MinigameCoins
	for i := 1; i < 4; i++ {
		maxMinigameCoins = max(maxMinigameCoins, g.Players[i].MinigameCoins)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].MinigameCoins == maxMinigameCoins {
			g.Players[i].Stars++
		}
	}

	maxHappening := g.Players[0].HappeningCount
	for i := 1; i < 4; i++ {
		maxHappening = max(maxHappening, g.Players[i].HappeningCount)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].HappeningCount == maxHappening {
			g.Players[i].Stars++
		}
	}
	return g
}

func Winners(g Game) []int {
	maxStarHolders := []int{}
	maxStars := g.Players[0].Stars
	for i := 1; i < 4; i++ {
		maxStars = max(maxStars, g.Players[i].Stars)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].Stars == maxStars {
			maxStarHolders = append(maxStarHolders, i)
		}
	}

	winners := []int{}
	maxCoins := g.Players[maxStarHolders[0]].Coins
	for i := 1; i < len(maxStarHolders); i++ {
		maxCoins = max(maxCoins, g.Players[maxStarHolders[i]].Coins)
	}
	for i := 0; i < len(maxStarHolders); i++ {
		if g.Players[maxStarHolders[i]].Coins == maxCoins {
			winners = append(winners, maxStarHolders[i])
		}
	}
	return winners
}

func EndGameTurn(g Game) Game {
	g.Turn++
	if g.Turn == g.Config.MaxTurns {
		g = AwardBonusStars(g)
		//Game is over, no more events
		g.ExtraEvent = nil
		return g
	}
	g.ExtraEvent = PickDiceBlock{g.CurrentPlayer, g.Config}
	return g
}

func StartMinigamePrep(g Game) Game {
	g = FindGreenPlayer(g)
	return g
}

func EndCharacterTurn(g Game) Game {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	if g.CurrentPlayer == 0 {
		return StartMinigamePrep(g)
	}
	g.ExtraEvent = PickDiceBlock{g.CurrentPlayer, g.Config}
	return g
}
