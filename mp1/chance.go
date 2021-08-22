package mp1

type ChanceMiddleBlock int

const (
	LTRStar ChanceMiddleBlock = iota
	LTR10
	LTR20
	LTR30
	RTLStar
	RTL10
	RTL20
	RTL30
	SwapCoins
	SwapStars

	CMBCount
)

type ChanceTimeBlock int

const (
	CTBLeft ChanceTimeBlock = iota
	CTBMiddle
	CTBRight
)

type ChanceTime struct {
	Player int

	LeftSideHit       bool
	LeftSidePosition  int
	MiddleHit         bool
	MiddlePosition    int
	RightSideHit      bool
	RightSidePosition int
}

type ChanceTimeResponse struct {
	Block    ChanceTimeBlock
	Position int
}

func (c ChanceTime) Responses() []Response {
	res := []Response{}
	if !c.LeftSideHit {
		for i := 0; i < 4; i++ {
			if c.RightSideHit && c.RightSidePosition == i {
				continue
			}
			res = append(res, ChanceTimeResponse{CTBLeft, i})
		}
	}
	if !c.MiddleHit {
		for i := 0; i < int(CMBCount); i++ {
			res = append(res, ChanceTimeResponse{CTBMiddle, i})
		}
	}
	if !c.RightSideHit {
		for i := 0; i < 4; i++ {
			if c.LeftSideHit && c.LeftSidePosition == i {
				continue
			}
			res = append(res, ChanceTimeResponse{CTBRight, i})
		}
	}
	return res
}

func (c ChanceTime) ControllingPlayer() int {
	count := 0
	if c.LeftSideHit {
		count++
	}
	if c.RightSideHit {
		count++
	}
	if c.MiddleHit {
		count++
	}

	if count < 2 { //Player can pick first 2 blocks
		return c.Player
	} else { //Block spins too fast for player, effectively random
		return CPU_PLAYER
	}
}

func (c ChanceTime) Handle(r Response, g *Game) {
	res := r.(ChanceTimeResponse)
	switch res.Block {
	case CTBLeft:
		c.LeftSideHit = true
		c.LeftSidePosition = res.Position
	case CTBMiddle:
		c.MiddleHit = true
		c.MiddlePosition = res.Position
	case CTBRight:
		c.RightSideHit = true
		c.RightSidePosition = res.Position
	}

	if c.LeftSideHit && c.MiddleHit && c.RightSideHit {
		middlePos := ChanceMiddleBlock(c.MiddlePosition)
		switch middlePos {
		case LTR10:
			g.GiveCoins(c.LeftSidePosition, c.RightSidePosition, 10, false)
		case LTR20:
			g.GiveCoins(c.LeftSidePosition, c.RightSidePosition, 20, false)
		case LTR30:
			g.GiveCoins(c.LeftSidePosition, c.RightSidePosition, 30, false)
		case LTRStar:
			if g.Players[c.LeftSidePosition].Stars > 0 {
				g.Players[c.LeftSidePosition].Stars--
				g.Players[c.RightSidePosition].Stars++
			}
		case RTL10:
			g.GiveCoins(c.RightSidePosition, c.LeftSidePosition, 10, false)
		case RTL20:
			g.GiveCoins(c.RightSidePosition, c.LeftSidePosition, 20, false)
		case RTL30:
			g.GiveCoins(c.RightSidePosition, c.LeftSidePosition, 30, false)
		case RTLStar:
			if g.Players[c.LeftSidePosition].Stars > 0 {
				g.Players[c.LeftSidePosition].Stars++
				g.Players[c.RightSidePosition].Stars--
			}
		case SwapCoins:
			tmp := g.Players[c.LeftSidePosition].Coins
			g.Players[c.LeftSidePosition].Coins = g.Players[c.RightSidePosition].Coins
			g.Players[c.RightSidePosition].Coins = tmp
		case SwapStars:
			tmp := g.Players[c.LeftSidePosition].Stars
			g.Players[c.LeftSidePosition].Stars = g.Players[c.RightSidePosition].Stars
			g.Players[c.RightSidePosition].Stars = tmp
		}
		g.EndCharacterTurn()
	} else {
		g.ExtraEvent = c
	}
}
