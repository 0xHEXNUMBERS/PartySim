package mp1

type Response interface{}

type Event interface {
	Responses() []Response
	Handle(Response, Game) Game
	ControllingPlayer() int
}

type BranchEvent struct {
	Player int
	Chain  int
	Moves  int
	Links  *[]ChainSpace
}

func (b BranchEvent) Responses() []Response {
	ret := []Response{nil}
	links := *b.Links
	for _, l := range links {
		ret = append(ret, l)
	}
	return ret
}

func (b BranchEvent) Handle(r Response, g Game) Game {
	g = ResetGameExtras(g)
	if r == nil {
		g.ExtraMovement = Movement{b.Player, b.Moves, false}
		return g
	}
	newPlayerPos := r.(ChainSpace)
	g.Players[b.Player].CurrentSpace = newPlayerPos
	g.ExtraMovement = Movement{b.Player, b.Moves, false}
	return g
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

func (p PayRangeEvent) Handle(r Response, g Game) Game {
	g = ResetGameExtras(g)
	cost := r.(int)
	g = AwardCoins(g, p.Player, -cost, false)
	g.ExtraMovement = Movement{p.Player, p.Moves, false}
	return g
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

func (m MushroomEvent) Handle(r Response, g Game) Game {
	g = ResetGameExtras(g)
	red := r.(bool)
	if red {
		g.ExtraMovement = Movement{Skip: true}
		return g
	}
	g.Players[m.Player].SkipTurn = true
	g.ExtraMovement = Movement{m.Player, 0, false}
	return g
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

func (b BooCoinsEvent) Handle(r Response, g Game) Game {
	g = ResetGameExtras(g)
	g = b.PayRangeEvent.Handle(r, g)
	g = AwardCoins(g, b.RecvPlayer, r.(int), false)
	g.ExtraMovement = Movement{b.RecvPlayer, b.PayRangeEvent.Moves, false}
	return g
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

func (b BooEvent) Handle(r Response, g Game) Game {
	g = ResetGameExtras(g)
	steal := r.(BooStealAction)
	if steal.Star {
		g = AwardCoins(g, steal.RecvPlayer, -50, false)
		g.Players[steal.GivingPlayer].Stars--
	} else {
		maxCoins := 15
		if b.Players[steal.GivingPlayer].Coins <= maxCoins {
			maxCoins = b.Players[steal.GivingPlayer].Coins
		}
		g.ExtraEvent = BooCoinsEvent{
			PayRangeEvent{steal.GivingPlayer, 1, maxCoins, b.Moves},
			steal.RecvPlayer,
		}
		return g
	}
	g.ExtraMovement = Movement{b.Player, b.Moves, false}
	return g
}

func (b BooEvent) ControllingPlayer() int {
	return b.Player
}

type DeterminePlayerTeamEvent struct {
	Player int
}

func (d DeterminePlayerTeamEvent) Responses() []Response {
	return []Response{true, false}
}

func (d DeterminePlayerTeamEvent) Handle(r Response, g Game) Game {
	g = ResetGameExtras(g)
	isBlue := r.(bool)

	if isBlue {
		g.Players[d.Player].LastSpaceType = Blue
	} else {
		g.Players[d.Player].LastSpaceType = Red
	}
	return g
}

func (d DeterminePlayerTeamEvent) ControllingPlayer() int {
	return CPU_PLAYER
}
