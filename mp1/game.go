package mp1

type Game struct {
	Board
	Turn          uint
	Players       [4]Player
	CurrentPlayer int
	CoinsOnStart  bool
	ExtraMovement Movement
	ExtraEvent    Event
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

func ResetGameExtras(g Game) Game {
	g.ExtraEvent = nil
	g.ExtraMovement.Skip = false
	return g
}

func MovePlayer(g Game, playerIdx, moves int) Game {
	if g.Players[playerIdx].SkipTurn {
		g.Players[playerIdx].SkipTurn = false
		return g
	}

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
			g = curSpace.PassingEvent(g, playerIdx, moves)
			if g.ExtraEvent != nil {
				g.Players[playerIdx].CurrentSpace = playerPos
				return g
			}
		case Start:
			if g.CoinsOnStart {
				g = AwardCoins(g, playerIdx, 10, false)
			}
		case Star:
			g = curSpace.PassingEvent(g, playerIdx, moves)
			if g.ExtraEvent != nil {
				return g
			}
		case Boo:
			g.ExtraEvent = BooEvent{playerIdx, g.Players, moves, g.Players[playerIdx].Coins}
			return g
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
	case Red:
		g = AwardCoins(g, playerIdx, -3, false)
	case Mushroom:
		g.ExtraEvent = MushroomEvent{playerIdx}
		return g
	case Happening:
		g.Players[playerIdx].HappeningCount++
		g := curSpace.StoppingEvent(g)
		if g.ExtraEvent != nil {
			return g
		}
	case Bowser:
		//Special events when player has 0 coins
		if g.Players[playerIdx].Coins == 0 {
			if g.Players[playerIdx].Stars > 0 {
				g = AwardCoins(g, playerIdx, 10, false)
				g.Players[playerIdx].Stars--
			} else {
				g = AwardCoins(g, playerIdx, 20, false)
			}
		} else {
			g.ExtraEvent = BowserEvent{playerIdx}
		}
	}
	//Switch Active Player
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	return g
}
