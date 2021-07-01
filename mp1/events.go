package mp1

type Response interface{}

type Event interface {
	Responses() []Response
	AffectedPlayer() int
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

func (b BranchEvent) AffectedPlayer() int {
	return b.Player
}

func (b BranchEvent) Handle(r Response, g *Game) Movement {
	if r == nil {
		return Movement{b.Player, b.Moves}
	}
	newPlayerPos := r.(ChainSpace)
	g.Players[b.Player].CurrentSpace = newPlayerPos
	return Movement{b.Player, b.Moves}

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

func (p PayRangeEvent) AffectedPlayer() int {
	return p.Player
}

func (p PayRangeEvent) Handle(r Response, g *Game) Movement {
	cost := r.(int)
	g.AwardCoins(p.Player, -cost, false)
	return Movement{p.Player, p.Moves}
}
