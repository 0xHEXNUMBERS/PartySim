package mp1

type Response interface{}

type Event interface {
	Responses() []Response
	AffectedPlayer() int
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
