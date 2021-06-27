package mp1

type SpaceType int

const (
	Invisible SpaceType = iota
	Blue
	Red
	Minigame
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
	Type         SpaceType
	Event        func(game *Game)
	PassingEvent func(game *Game, player, moves int) Event
}

type Chain []Space

type ChainSpace struct {
	Chain int
	Space int
}

type Board struct {
	Chains []Chain
	Links  map[int][]ChainSpace
}

func (b Board) Copy() Board {
	chains := make([]Chain, 0)
	for _, c := range b.Chains {
		chain := make([]Space, 0)
		for _, s := range c {
			chain = append(chain, s)
		}
		chains = append(chains, chain)
	}
	links := make(map[int][]ChainSpace)
	for i, s := range b.Links {
		slice := make([]ChainSpace, 0)
		for _, j := range s {
			slice = append(slice, j)
		}
		links[i] = s
	}
	return Board{Chains: chains, Links: links}
}
