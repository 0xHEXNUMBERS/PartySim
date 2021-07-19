package mp1

type StarData struct {
	StarSpaceCount   int
	AbsoluteVisited  uint64 //For Chance vs. Blue
	RelativeVisited  uint64 //For determining next star space
	CurrentStarSpace uint8  //Index into RelativeVisited
	IndexToPosition  *[]ChainSpace
}

func (s StarData) Responses() []Response {
	res := []Response{}
	if s.RelativeVisited == (1<<s.StarSpaceCount)-1 {
		for i := 0; i < s.StarSpaceCount; i++ {
			if i != int(s.CurrentStarSpace) {
				res = append(res, i)
			}
		}
	} else {
		for i := 0; i < s.StarSpaceCount; i++ {
			if s.RelativeVisited|(1<<i) == 0 {
				res = append(res, i)
			}
		}
	}
	return res
}

func (s StarData) ControllingPlayer() int {
	return CPU_PLAYER
}

//TODO: Implement Handle(r, g)
