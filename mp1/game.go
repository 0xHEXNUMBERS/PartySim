package mp1

import "math/rand"

type Game struct {
	Board
	Players [4]Player
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
			g.Players[playerIdx].Coins -= 20
			g.Players[playerIdx].Stars++
		}
	case BlackStar:
		if g.Players[playerIdx].Stars > 0 {
			g.Players[playerIdx].Stars--
		} else {
			if g.Players[playerIdx].Coins >= 20 {
				g.Players[playerIdx].Coins = 0
			} else {
				g.Players[playerIdx].Coins -= 20
			}
		}
	case Mushroom:
		if rand.Intn(2) == 0 {
			g.Players[playerIdx].SkipTurn = true
		} else {
			return MushroomEvent{playerIdx}
		}
	}
	g.Players[playerIdx].CurrentSpace = playerPos
	return nil
}
