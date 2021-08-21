package mp1

type esBranchEvent struct {
	Player int
	Moves  int
	Warp1  ChainSpace
	Warp2  ChainSpace
	Warp3  ChainSpace
}

func (e esBranchEvent) Responses() []Response {
	return []Response{true, false}
}

func (e esBranchEvent) ControllingPlayer() int {
	return e.Player
}

func (e esBranchEvent) Handle(r Response, g *Game) {
	gotoWarp := r.(bool)
	bd := g.Board.Data.(esBoardData)
	if gotoWarp {
		switch bd.Gate {
		case 0:
			g.ExtraEvent = esWarpDest{
				e.Player,
				e.Moves,
				bd.Gate2or3,
				e.Warp1,
				e.Warp2,
				e.Warp3,
			}
			return
		case 1:
			g.Players[e.Player].CurrentSpace = e.Warp1
		case 2:
			g.Players[e.Player].CurrentSpace = e.Warp2
		case 3:
			g.Players[e.Player].CurrentSpace = e.Warp3
		}
		g.MovePlayer(e.Player, e.Moves)
	} else {
		g.MovePlayer(e.Player, e.Moves)
	}
}

type esVisitBabyBowser struct {
	Player int
	Moves  int
	Index  int
}

func (e esVisitBabyBowser) Responses() []Response {
	return []Response{true, false}
}

func (e esVisitBabyBowser) ControllingPlayer() int {
	return e.Player
}

func (e esVisitBabyBowser) Handle(r Response, g *Game) {
	battle := r.(bool)
	if battle {
		g.AwardCoins(e.Player, -20, false)
		g.ExtraEvent = esBattleBabyBowser{
			e.Player, e.Moves, e.Index,
		}
	} else {
		g.MovePlayer(e.Player, e.Moves)
	}
}

type esBattleBabyBowser struct {
	Player int
	Moves  int
	Index  int
}

func (e esBattleBabyBowser) Responses() []Response {
	return []Response{true, false}
}

func (e esBattleBabyBowser) ControllingPlayer() int {
	return CPU_PLAYER
}

func (e esBattleBabyBowser) Handle(r Response, g *Game) {
	star := r.(bool)
	bd := g.Board.Data.(esBoardData)
	if star {
		g.Players[e.Player].Stars++
		bd.StarTaken[e.Index] = true
		if esAllStarsCollected(bd) {
			bd.StarTaken = [7]bool{
				false, false, false, false, false, false, false,
			}
		}
		g.Board.Data = bd
	} else if g.Players[e.Player].Stars > 0 {
		g.Players[e.Player].Stars--
	}
	g.MovePlayer(e.Player, e.Moves)
}

type esWarpCDest struct {
	Player int
	Moves  int
}

func (e esWarpCDest) Responses() []Response {
	return []Response{esEntrance1, esEntrance7}
}

func (e esWarpCDest) ControllingPlayer() int {
	return CPU_PLAYER
}

func (e esWarpCDest) Handle(r Response, g *Game) {
	dest := r.(ChainSpace)
	g.Players[e.Player].CurrentSpace = dest

	bd := g.Board.Data.(esBoardData)
	if dest == esEntrance7 {
		bd.Gate = 1
	} else {
		bd.Gate2or3 = true
	}
	g.Board.Data = bd

	g.MovePlayer(e.Player, e.Moves)
}

type esWarpDestResponse struct {
	Dest ChainSpace
	Gate int
}

type esWarpDest struct {
	Player   int
	Moves    int
	Gate2or3 bool
	Island1  ChainSpace
	Island2  ChainSpace
	Island3  ChainSpace
}

func (e esWarpDest) Responses() []Response {
	ret := []Response{
		esWarpDestResponse{e.Island1, 1},
		esWarpDestResponse{e.Island2, 2},
		esWarpDestResponse{e.Island3, 3},
	}
	if e.Gate2or3 {
		ret = ret[1:]
	}
	return ret
}

func (e esWarpDest) ControllingPlayer() int {
	return CPU_PLAYER
}

func (e esWarpDest) Handle(r Response, g *Game) {
	dest := r.(esWarpDestResponse)
	bd := g.Board.Data.(esBoardData)
	bd.Gate = dest.Gate
	g.Board.Data = bd
	g.Players[e.Player].CurrentSpace = dest.Dest
	g.MovePlayer(e.Player, e.Moves)
}

type esChangeGates struct {
	Player  int
	Moves   int
	Current int
}

var esChangeGatesResponses = [3][]Response{
	{2, 3},
	{1, 3},
	{1, 2},
}

func (e esChangeGates) Responses() []Response {
	return esChangeGatesResponses[e.Current-1]
}

func (e esChangeGates) ControllingPlayer() int {
	return CPU_PLAYER
}

func (e esChangeGates) Handle(r Response, g *Game) {
	gate := r.(int)
	bd := g.Board.Data.(esBoardData)
	bd.Gate = gate
	bd.Gate2or3 = (gate != 1)
	g.Players[e.Player].CurrentSpace = ChainSpace{0, 0}
	g.MovePlayer(e.Player, e.Moves-1)
}
