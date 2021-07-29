package mp1

type StarData struct {
	StarSpaceCount   uint8
	AbsoluteVisited  uint64 //For Chance vs. Blue
	RelativeVisited  uint64 //For determining next star space
	CurrentStarSpace ChainSpace
	IndexToPosition  *[]ChainSpace
}

func (s StarData) GetIndex(c ChainSpace) int {
	for i := 0; i < len(*s.IndexToPosition); i++ {
		tmp := (*s.IndexToPosition)[i]
		if c == tmp {
			return i
		} /* else if c.Chain < tmp.Chain || (c.Chain == tmp.Chain && c.Space < tmp.Space) {
			return -1
		}*/
	}
	return -1
}

type StarLocationEvent struct {
	StarData
	Player int
	Moves  int
}

func (s StarLocationEvent) Responses() []Response {
	var i uint8
	res := []Response{}
	for i = 0; i < s.StarSpaceCount; i++ {
		if s.RelativeVisited|(1<<i) == 0 && i != s.StarSpaceCount {
			res = append(res, i)
		}
	}
	return res
}

func (s StarLocationEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (s StarLocationEvent) Handle(r Response, g *Game) {
	i := r.(uint8)
	s.AbsoluteVisited |= (1 << i)
	s.RelativeVisited |= (1 << i)
	s.CurrentStarSpace = (*s.IndexToPosition)[i]

	if s.RelativeVisited == (1<<s.StarSpaceCount)-1 { //Clear relative count if all star spaces are visited
		s.RelativeVisited = 0
	}

	g.StarSpaces = s.StarData
	if s.Moves != 0 {
		g.MovePlayer(s.Player, s.Moves)
	} else {
		g.ExtraEvent = PickDiceBlock{s.Player, g.Config}
	}
}
