package mp1

type SpaceType int

const (
	Invisible SpaceType = iota
	Blue
	Red
	MinigameSpace
	Happening
	Star
	Chance
	Start
	Mushroom
	Bowser
	BlackStar
	Boo
)

type Space struct {
	Type          SpaceType
	StoppingEvent func(game Game) Game
	PassingEvent  func(game Game, player, moves int) Game
}

type Chain []Space

type ChainSpace struct {
	Chain int
	Space int
}

type ExtraBoardData interface {
	Copy() ExtraBoardData
}

type Movement struct {
	Player int
	Moves  int
	Skip   bool
}

type Board struct {
	Chains *[]Chain
	Links  map[int]*[]ChainSpace
	Data   ExtraBoardData
}
