package mp1

type NormalDiceBlock struct {
	Player int
}

func (m NormalDiceBlock) Responses() []Response {
	return CPURangeEvent{1, 10}.Responses()
}

func (m NormalDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

func (m NormalDiceBlock) Handle(r Response, g Game) Game {
	moves := r.(int)
	return MovePlayer(g, m.Player, moves)
}

type RedDiceBlock struct {
	Player int
}

func (r RedDiceBlock) Responses() []Response {
	return CPURangeEvent{1, 10}.Responses()
}

func (r RedDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

func (r RedDiceBlock) Handle(res Response, g Game) Game {
	coinsLost := res.(int)
	g = AwardCoins(g, r.Player, -coinsLost, false)
	return MovePlayer(g, r.Player, coinsLost)
}

type BlueDiceBlock struct {
	Player int
}

func (b BlueDiceBlock) Responses() []Response {
	return CPURangeEvent{1, 10}.Responses()
}

func (b BlueDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

func (b BlueDiceBlock) Handle(r Response, g Game) Game {
	coinsWon := r.(int)
	g = AwardCoins(g, b.Player, coinsWon, false)
	return MovePlayer(g, b.Player, coinsWon)
}

type WarpDiceBlock struct {
	Player int
}

func (w WarpDiceBlock) Responses() []Response {
	var res [3]Response
	i := 0
	for player := 0; player < 4; player++ {
		if player != w.Player {
			res[i] = player
			i++
		}
	}
	return res[:]
}

func (w WarpDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

func (w WarpDiceBlock) Handle(r Response, g Game) Game {
	selectedPlayer := r.(int)
	tmpSpace := g.Players[w.Player].CurrentSpace
	g.Players[w.Player].CurrentSpace = g.Players[selectedPlayer].CurrentSpace
	g.Players[selectedPlayer].CurrentSpace = tmpSpace
	return g
}

type EventDiceBlock struct {
	Player int
}

type EventBlockEvent int

const (
	BooEventBlock EventBlockEvent = iota
	BowserEventBlock
	KoopaEventBlock
)

var EventBlockResponses = []Response{
	BooEventBlock,
	BowserEventBlock,
	KoopaEventBlock,
}

func (e EventDiceBlock) Responses() []Response {
	return EventBlockResponses
}

func (e EventDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

func (e EventDiceBlock) Handle(r Response, g Game) Game {
	event := r.(EventBlockEvent)
	switch event {
	case BooEventBlock:
		g.ExtraEvent = BooEvent{
			e.Player,
			g.Players,
			0,
			g.Players[e.Player].Coins,
		}
	case BowserEventBlock:
		g = PreBowserCheck(g, e.Player)
	case KoopaEventBlock:
		g = AwardCoins(g, e.Player, 10, false)
		g.ExtraEvent = nil
	}
	return g
}
