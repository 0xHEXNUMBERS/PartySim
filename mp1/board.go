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
	BogusItem
	Boo
)

type Space struct {
	Type          SpaceType
	StoppingEvent func(game *Game, player int)
	PassingEvent  func(game *Game, player, moves int)
}

type Chain []Space

type ChainSpace struct {
	Chain int
	Space int
}

type ExtraBoardData interface{}

type Movement struct {
	Player int
	Moves  int
	Skip   bool
}

type Board struct {
	Chains      *[]Chain
	Links       *map[int]*[]ChainSpace
	BowserCoins int
	Data        ExtraBoardData
}
