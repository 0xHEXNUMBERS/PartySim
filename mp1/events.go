package mp1

type Response interface{}

type Event interface {
	Responses() []Response
	Handle(Response, *Game) Movement
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
		return Movement{b.Player, b.Moves, false}
	}
	newPlayerPos := r.(ChainSpace)
	g.Players[b.Player].CurrentSpace = newPlayerPos
	return Movement{b.Player, b.Moves, false}

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
	return Movement{p.Player, p.Moves, false}
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
	return Movement{m.Player, 0, false}
}
