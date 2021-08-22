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

func (m NormalDiceBlock) Handle(r Response, g *Game) {
	moves := r.(int)
	g.MovePlayer(m.Player, moves)
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

func (r RedDiceBlock) Handle(res Response, g *Game) {
	coinsLost := res.(int)
	g.AwardCoins(r.Player, -coinsLost, false)
	g.MovePlayer(r.Player, coinsLost)
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

func (b BlueDiceBlock) Handle(r Response, g *Game) {
	coinsWon := r.(int)
	g.AwardCoins(b.Player, coinsWon, false)
	g.MovePlayer(b.Player, coinsWon)
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

func (w WarpDiceBlock) Handle(r Response, g *Game) {
	selectedPlayer := r.(int)
	tmpSpace := g.Players[w.Player].CurrentSpace
	g.Players[w.Player].CurrentSpace = g.Players[selectedPlayer].CurrentSpace
	g.Players[selectedPlayer].CurrentSpace = tmpSpace
	g.ActivateSpace(w.Player)
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

func (e EventDiceBlock) Handle(r Response, g *Game) {
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
		//TODO: Typically bowser just takes 20 coins
		//Does anything happen if player has 0 coins?
		g.AwardCoins(e.Player, -20, false)
		g.EndCharacterTurn()
	case KoopaEventBlock:
		g.AwardCoins(e.Player, 10, false)
		g.EndCharacterTurn()
	}
}

type PickDiceBlock struct {
	Player int
	Config GameConfig
}

func (p PickDiceBlock) Responses() []Response {
	res := []Response{NormalDiceBlock{p.Player}}
	if p.Config.RedDice {
		res = append(res, RedDiceBlock{p.Player})
	}
	if p.Config.BlueDice {
		res = append(res, BlueDiceBlock{p.Player})
	}
	if p.Config.WarpDice {
		res = append(res, WarpDiceBlock{p.Player})
	}
	if p.Config.EventsDice {
		res = append(res, EventDiceBlock{p.Player})
	}
	return res
}

func (p PickDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

func (p PickDiceBlock) Handle(r Response, g *Game) {
	evt := r.(Event)
	g.ExtraEvent = evt
}
