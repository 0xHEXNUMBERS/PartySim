package mp1

type Response interface{}

type Event interface {
	Responses() []Response
	Handle(Response, *Game) Movement
	ControllingPlayer() int
}

type BranchEvent struct {
	Player int
	Chain  int
	Moves  int
	Links  []ChainSpace
}

func (b BranchEvent) Responses() []Response {
	ret := []Response{nil}
	for _, l := range b.Links {
		ret = append(ret, l)
	}
	return ret
}

func (b BranchEvent) Handle(r Response, g *Game) Movement {
	if r == nil {
		return Movement{b.Player, b.Moves, false, nil}
	}
	newPlayerPos := r.(ChainSpace)
	g.Players[b.Player].CurrentSpace = newPlayerPos
	return Movement{b.Player, b.Moves, false, nil}

}

func (b BranchEvent) ControllingPlayer() int {
	return b.Player
}

type PayRangeEvent struct {
	Player int
	Min    int
	Max    int
	Moves  int
}

func (p PayRangeEvent) Responses() []Response {
	ret := make([]Response, (p.Max-p.Min)+1)
	for i := p.Min; i <= p.Max; i++ {
		ret[i-p.Min] = i
	}
	return ret
}

func (p PayRangeEvent) Handle(r Response, g *Game) Movement {
	cost := r.(int)
	g.AwardCoins(p.Player, -cost, false)
	return Movement{p.Player, p.Moves, false, nil}
}

func (p PayRangeEvent) ControllingPlayer() int {
	return p.Player
}

type MushroomEvent struct {
	Player int
}

func (m MushroomEvent) Responses() []Response {
	return []Response{false, true}
}

func (m MushroomEvent) Handle(r Response, g *Game) Movement {
	red := r.(bool)
	if red {
		return Movement{Skip: true}
	}
	g.Players[m.Player].SkipTurn = true
	return Movement{m.Player, 0, false, nil}
}

func (m MushroomEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

type BooCoinsEvent struct {
	PayRangeEvent
	RecvPlayer int
}

func (b BooCoinsEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BooCoinsEvent) Handle(r Response, g *Game) Movement {
	b.PayRangeEvent.Handle(r, g)
	g.AwardCoins(b.RecvPlayer, r.(int), false)
	return Movement{b.RecvPlayer, b.PayRangeEvent.Moves, false, nil}
}

type BooEvent struct {
	Player  int
	Players [4]Player
	Moves   int
	Coins   int
}

type BooStealAction struct {
	RecvPlayer   int
	GivingPlayer int
	Star         bool
}

func (b BooEvent) Responses() []Response {
	res := make([]Response, 0)
	if b.Coins >= 50 {
		for i := 0; i < 4; i++ {
			if i == b.Player {
				continue
			}
			if b.Players[i].Stars > 0 {
				res = append(res, BooStealAction{b.Player, i, true})
			}
		}
	}
	for i := 0; i < 4; i++ {
		if i == b.Player {
			continue
		}
		if b.Players[i].Coins > 0 {
			res = append(res, BooStealAction{b.Player, i, false})
		}
	}
	return res
}

func (b BooEvent) Handle(r Response, g *Game) Movement {
	steal := r.(BooStealAction)
	if steal.Star {
		g.AwardCoins(steal.RecvPlayer, -50, false)
		g.Players[steal.GivingPlayer].Stars--
	} else {
		return Movement{
			ExtraEvent: BooCoinsEvent{
				PayRangeEvent{steal.GivingPlayer, 1, 15, b.Moves},
				steal.RecvPlayer,
			},
		}
	}
	return Movement{b.Player, b.Moves, false, nil}
}

func (b BooEvent) ControllingPlayer() int {
	return b.Player
}
