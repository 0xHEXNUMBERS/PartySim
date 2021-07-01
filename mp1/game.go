package mp1

import "math/rand"

type Game struct {
	Board
	Players       [4]Player
	CurrentPlayer int
	CoinsOnStart  bool
}

func (g *Game) AwardCoins(player, coins int, minigame bool) {
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
}

func (g *Game) MovePlayer(playerIdx, moves int) (e Event) {
	if g.Players[playerIdx].SkipTurn {
		g.Players[playerIdx].SkipTurn = false
		return nil
	}
	playerPos := g.Players[playerIdx].CurrentSpace
	for moves > 0 {
		playerPos.Space++
		if playerPos.Space >= len(g.Board.Chains[playerPos.Chain]) {
			playerPos.Space = 0
		}
		curSpace := g.Board.Chains[playerPos.Chain][playerPos.Space]
		switch curSpace.Type {
		case Invisible:
			evt := curSpace.PassingEvent(g, playerIdx, moves)
			if evt != nil {
				g.Players[playerIdx].CurrentSpace = playerPos
				return evt
			}
		case Start:
			if g.CoinsOnStart {
				g.AwardCoins(playerIdx, 10, false)
			}
		default:
			moves--
		}
	}
	//Stop on Space
	g.Players[playerIdx].CurrentSpace = playerPos
	//Activate Space
	curSpace := g.Board.Chains[playerPos.Chain][playerPos.Space]
	switch curSpace.Type {
	case Blue:
		g.AwardCoins(playerIdx, 3, false)
	case Red:
		g.AwardCoins(playerIdx, -3, false)
	case Star:
		if g.Players[playerIdx].Coins >= 20 {
			g.AwardCoins(playerIdx, -20, false)
			g.Players[playerIdx].Stars++
		}
	case BlackStar:
		if g.Players[playerIdx].Stars > 0 {
			g.Players[playerIdx].Stars--
		} else {
			g.AwardCoins(playerIdx, -20, false)
		}
	case Mushroom:
		if rand.Intn(2) == 0 {
			g.Players[playerIdx].SkipTurn = true
		} else {
			return nil
		}
	case Happening:
		g.Players[playerIdx].HappeningCount++
		evt := curSpace.StoppingEvent(g)
		if evt != nil {
			return evt
		}
	}
	//Switch Active Player
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	return nil
}
