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
	if r != nil {
		newPlayerPos := r.(ChainSpace)
		g.Players[b.Player].CurrentSpace = newPlayerPos
	}
	g = MovePlayer(g, b.Player, b.Moves)
	return g
}

func (b BranchEvent) ControllingPlayer() int {
	return b.Player
}

type PayRangeEvent struct {
	Player int
	Min    int
	Max    int
}

func (p PayRangeEvent) Responses() []Response {
	ret := make([]Response, (p.Max-p.Min)+1)
	for i := p.Min; i <= p.Max; i++ {
		ret[i-p.Min] = i
	}
	return ret
}

func (p PayRangeEvent) Handle(r Response, g Game) Game {
	//TODO: Figure out if we need to add EndCharacterTurn()
	cost := r.(int)
	g = AwardCoins(g, p.Player, -cost, false)
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
	redMushroom := r.(bool)
	if redMushroom {
		g.ExtraEvent = NormalDiceBlock{m.Player}
		return g
	}
	g.Players[m.Player].SkipTurn = true
	g = EndCharacterTurn(g)
	return g
}

func (m MushroomEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

type BooCoinsEvent struct {
	PayRangeEvent
	RecvPlayer int
	Moves      int
}

func (b BooCoinsEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BooCoinsEvent) Handle(r Response, g Game) Game {
	g = b.PayRangeEvent.Handle(r, g)
	g = AwardCoins(g, b.RecvPlayer, r.(int), false)

	if b.Moves != 0 {
		g = MovePlayer(g, b.RecvPlayer, b.Moves)
	} else {
		g = EndCharacterTurn(g)
	}
	return g
}

type BooEvent struct {
	Player  int
	Players [4]Player
	Moves   int //No call to MovePlayer on 0
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
			PayRangeEvent{steal.GivingPlayer, 1, maxCoins},
			steal.RecvPlayer,
			b.Moves,
		}
		return g
	}
	if b.Moves != 0 {
		g = MovePlayer(g, b.Player, b.Moves)
	} else {
		g = EndCharacterTurn(g)
	}
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
	isBlue := r.(bool)

	if isBlue {
		g.Players[d.Player].LastSpaceType = Blue
	} else {
		g.Players[d.Player].LastSpaceType = Red
	}
	g = FindGreenPlayer(g)
	return g
}

func (d DeterminePlayerTeamEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

type CPURangeEvent struct {
	Min int
	Max int
}

func (c CPURangeEvent) Responses() []Response {
	var ret []Response
	for i := c.Min; i <= c.Max; i++ {
		ret = append(ret, i)
	}
	return ret
}

func (c CPURangeEvent) ControllingPlayer() int {
	return CPU_PLAYER
}
