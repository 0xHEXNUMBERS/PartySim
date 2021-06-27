package mp1

type Game struct {
	Board
	Players [4]Player
}

func (g *Game) MovePlayer(playerIdx, moves int) (e Event) {
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
	}
	g.Players[playerIdx].CurrentSpace = playerPos
	return nil
}
