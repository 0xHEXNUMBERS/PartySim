package mp1

//NormalDiceBlock holds the implementation of a regular dice block.
type NormalDiceBlock struct {
	Player int
}

//Responses returns a slice of ints from [1, 10].
func (m NormalDiceBlock) Responses() []Response {
	return CPURangeEvent{1, 10}.Responses()
}

func (m NormalDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle moves the player r spaces.
func (m NormalDiceBlock) Handle(r Response, g *Game) {
	moves := r.(int)
	g.MovePlayer(m.Player, moves)
}

//RedDiceBlock holds the implementation of a red dice block.
type RedDiceBlock struct {
	Player int
}

//Responses returns a slice of ints from [1, 10].
func (r RedDiceBlock) Responses() []Response {
	return CPURangeEvent{1, 10}.Responses()
}

func (r RedDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle takes r coins from the player, and moves the player r spaces.
func (r RedDiceBlock) Handle(res Response, g *Game) {
	coinsLost := res.(int)
	g.AwardCoins(r.Player, -coinsLost, false)
	g.MovePlayer(r.Player, coinsLost)
}

//BlueDiceBlock holds the implementation of a blue dice block.
type BlueDiceBlock struct {
	Player int
}

//Responses returns a slice of ints from [1, 10].
func (b BlueDiceBlock) Responses() []Response {
	return CPURangeEvent{1, 10}.Responses()
}

func (b BlueDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives r coins from the player, and moves the player r spaces.
func (b BlueDiceBlock) Handle(r Response, g *Game) {
	coinsWon := r.(int)
	g.AwardCoins(b.Player, coinsWon, false)
	g.MovePlayer(b.Player, coinsWon)
}

//WarpDiceBlock holds the implementation of a warp dice block.
type WarpDiceBlock struct {
	Player int
}

//Responses returns a slice of ints containing the indexes of the other
//players.
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

//Handle warps player w.Player to player r's position, and player r to
//player w.Player's position. The player w.Player then lands on their new
//space.
func (w WarpDiceBlock) Handle(r Response, g *Game) {
	selectedPlayer := r.(int)
	tmpSpace := g.Players[w.Player].CurrentSpace
	g.Players[w.Player].CurrentSpace = g.Players[selectedPlayer].CurrentSpace
	g.Players[selectedPlayer].CurrentSpace = tmpSpace
	g.ActivateSpace(w.Player)
}

//EventDiceBlock hodls the implementation for an event dice block.
type EventDiceBlock struct {
	Player int
}

//EventBlockEvent is an enumeration of the possible event dice block
//actions.
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

//Responses returns a slice of the possible event dice block actions.
func (e EventDiceBlock) Responses() []Response {
	return EventBlockResponses
}

func (e EventDiceBlock) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle performs the EventBlockEvent r on the game, setting the next
//event if needed.
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

//PickDiceBlock holds the data for deciding the next dice block that
//appears.
type PickDiceBlock struct {
	Player int
	Config GameConfig
}

//Responses returns a slice of the available dice blocks that can appear
//based on the game's configuration.
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

//Handle sets the next event to the dice block r.
func (p PickDiceBlock) Handle(r Response, g *Game) {
	evt := r.(Event)
	g.ExtraEvent = evt
}
