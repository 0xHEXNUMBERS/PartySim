package mp1

type ChanceMiddleBlock int

const (
	LTRStar ChanceMiddleBlock = iota
	LTR10
	LTR20
	LTR30
	SwapCoins
	SwapStars

	CMBCount
)

type ChanceTimeBlock int

const (
	CTBSide ChanceTimeBlock = iota
	CTBMiddle
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
			res = append(res, ChanceTimeResponse{CTBSide, i})
		}
	}
	if !c.MiddleHit {
		for i := 0; i < int(CMBCount); i++ {
			res = append(res, ChanceTimeResponse{CTBMiddle, i})
		}
	}
	if !c.RightSideHit && c.LeftSideHit {
		for i := 0; i < 4; i++ {
			if c.LeftSidePosition != i {
				res = append(res, ChanceTimeResponse{CTBSide, i})
			}
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

func (c ChanceTime) Handle(r Response, g Game) Game {
	res := r.(ChanceTimeResponse)
	if res.Block == CTBSide {
		if c.LeftSideHit {
			c.RightSideHit = true
			c.RightSidePosition = res.Position
		} else {
			c.LeftSideHit = true
			c.LeftSidePosition = res.Position
		}
	} else {
		c.MiddleHit = true
		c.MiddlePosition = res.Position
	}

	if c.LeftSideHit && c.MiddleHit && c.RightSideHit {
		middlePos := ChanceMiddleBlock(c.MiddlePosition)
		switch middlePos {
		case LTR10:
			coinsTaken := min(g.Players[c.LeftSidePosition].Coins, 10)
			g = AwardCoins(g, c.LeftSidePosition, -coinsTaken, false)
			g = AwardCoins(g, c.RightSidePosition, coinsTaken, false)
		case LTR20:
			coinsTaken := min(g.Players[c.LeftSidePosition].Coins, 20)
			g = AwardCoins(g, c.LeftSidePosition, -coinsTaken, false)
			g = AwardCoins(g, c.RightSidePosition, coinsTaken, false)
		case LTR30:
			coinsTaken := min(g.Players[c.LeftSidePosition].Coins, 30)
			g = AwardCoins(g, c.LeftSidePosition, -coinsTaken, false)
			g = AwardCoins(g, c.RightSidePosition, coinsTaken, false)
		case LTRStar:
			if g.Players[c.LeftSidePosition].Stars > 0 {
				g.Players[c.LeftSidePosition].Stars--
				g.Players[c.RightSidePosition].Stars++
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
		g.ExtraEvent = nil
	} else {
		g.ExtraEvent = c
	}
	return g
}
