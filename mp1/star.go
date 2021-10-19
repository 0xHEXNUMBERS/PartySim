package mp1

//StarData holds the data for star locations, which stars have been
//collected at least once, and which stars have been collected recently.
type StarData struct {
	StarSpaceCount   uint8
	AbsoluteVisited  uint64 //For Chance vs. Blue
	RelativeVisited  uint64 //For determining next star space
	CurrentStarSpace ChainSpace
	IndexToPosition  *[]ChainSpace
}

//GetIndex is a mapping from ChainSpace to the internal s.IndexToPosition
//slice index. Returns -1 if c is not in s.IndexToPosition.
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

//StarLocationEvent holds the implementation for picking a new star space.
type StarLocationEvent struct {
	StarData
	Player int
	Moves  int
}

func (s StarLocationEvent) Type() EventType {
	return CHAINSPACE_EVT_TYPE
}

//Responses returns a slice of the available indexes of the available star
//spaces.
func (s StarLocationEvent) Responses() []Response {
	var i int
	res := []Response{}
	for i = 0; i < int(s.StarSpaceCount); i++ {
		if s.RelativeVisited&(1<<i) == 0 {
			res = append(res, (*s.IndexToPosition)[i])
		}
	}
	return res
}

func (s StarLocationEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle takes the index r and sets the new star space to that index. If
//r is the last available star space, then the list of star spaces already
//landed on is reset.
func (s StarLocationEvent) Handle(r Response, g *Game) {
	c := r.(ChainSpace)
	i := s.GetIndex(c)
	if i < 0 { //Error
		return
	}

	s.AbsoluteVisited |= (1 << i)
	s.RelativeVisited |= (1 << i)
	s.CurrentStarSpace = c

	if s.RelativeVisited == (1<<s.StarSpaceCount)-1 { //Clear relative count if all star spaces are visited
		s.RelativeVisited = 0
	}

	g.StarSpaces = s.StarData
	if s.Moves != 0 {
		g.MovePlayer(s.Player, s.Moves)
	} else {
		g.SetDiceBlock()
	}
}
