package mp1

import "math/rand"

type Game struct {
	Board
	Players       [4]Player
	CurrentPlayer int
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
		default:
			moves--
		}
	}
	//Activate Space
	curSpace := g.Board.Chains[playerPos.Chain][playerPos.Space]
	switch curSpace.Type {
	case Blue:
		g.Players[playerIdx].Coins += 3
	case Red:
		g.Players[playerIdx].Coins -= 3
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
	}
	g.Players[playerIdx].CurrentSpace = playerPos
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	return nil
}
