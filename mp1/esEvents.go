package mp1

//esBranchEvent let's the player decide if they want to take the warp.
type esBranchEvent struct {
	Player int
	Moves  int
	Warp1  ChainSpace
	Warp2  ChainSpace
	Warp3  ChainSpace
}

//Responses returns a slice of bools (true/false).
func (e esBranchEvent) Responses() []Response {
	return []Response{true, false}
}

func (e esBranchEvent) ControllingPlayer() int {
	return e.Player
}

//Handle executes based on r. If r is true, the player's new position is set
//based on the current gate (setting the next event to set the gate if
//the gate is unknown). If r is false, the player continues down their
//current chain.
func (e esBranchEvent) Handle(r Response, g *Game) {
	gotoWarp := r.(bool)
	bd := g.Board.Data.(esBoardData)
	if gotoWarp {
		switch bd.Gate {
		case 0:
			g.NextEvent = esWarpDest{
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

//esVisitBabyBowser let's the player decide if they want to play baby
//bowser's minigame to win a star.
type esVisitBabyBowser struct {
	Player int
	Moves  int
	Index  int
}

//Responses return a slice of bools (true/false).
func (e esVisitBabyBowser) Responses() []Response {
	return []Response{true, false}
}

func (e esVisitBabyBowser) ControllingPlayer() int {
	return e.Player
}

//Handle sets the next event to the baby bowser minigame if r is true. If r
//is false, then nothing happens.
func (e esVisitBabyBowser) Handle(r Response, g *Game) {
	battle := r.(bool)
	if battle {
		g.AwardCoins(e.Player, -20, false)
		g.NextEvent = esBattleBabyBowser{
			e.Player, e.Moves, e.Index,
		}
	} else {
		g.MovePlayer(e.Player, e.Moves)
	}
}

//esBattleBabyBowser decides if the player wins the minigame.
type esBattleBabyBowser struct {
	Player int
	Moves  int
	Index  int
}

//Responses return a slice of bools (true/false).
func (e esBattleBabyBowser) Responses() []Response {
	return []Response{true, false}
}

func (e esBattleBabyBowser) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives the player a star and sets the baby bowser's StarTaken flag
//to true if r is true. If r is false, a star is taken from the plaeyr.
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

//esWarpCDest decides which Warp C destination the player goes to.
type esWarpCDest struct {
	Player int
	Moves  int
}

//Resopnses returns a slice of the 2 possible spaces the player can warp
//to.
func (e esWarpCDest) Responses() []Response {
	return []Response{esEntrance1, esEntrance7}
}

func (e esWarpCDest) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle moves the player to the ChainSpace r and sets various flags if
//needed.
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

//esWarpDestResponse is a response that can be made to a esWarpDest Event.
type esWarpDestResponse struct {
	Dest ChainSpace
	Gate int
}

//esWarpDest decides which gate the board is playing under currently.
type esWarpDest struct {
	Player   int
	Moves    int
	Gate2or3 bool
	Island1  ChainSpace
	Island2  ChainSpace
	Island3  ChainSpace
}

//Responses returns the list of possible ChainSpaces that the player can
//warp to.
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

//Handle moves the player to the ChainSpace in r and set's the current
//gate the board is under in r.
func (e esWarpDest) Handle(r Response, g *Game) {
	dest := r.(esWarpDestResponse)
	bd := g.Board.Data.(esBoardData)
	bd.Gate = dest.Gate
	g.Board.Data = bd
	g.Players[e.Player].CurrentSpace = dest.Dest
	g.MovePlayer(e.Player, e.Moves)
}

//esChangeGates decides which Gate the board will change to.
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

//Responses returns the gates that can be switched to.
func (e esChangeGates) Responses() []Response {
	return esChangeGatesResponses[e.Current-1]
}

func (e esChangeGates) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle switches the current gate configuration to r, moves the player to
//the starting space, and moves the player their remaining spaces.
func (e esChangeGates) Handle(r Response, g *Game) {
	gate := r.(int)
	bd := g.Board.Data.(esBoardData)
	bd.Gate = gate
	bd.Gate2or3 = (gate != 1)
	g.Board.Data = bd
	g.Players[e.Player].CurrentSpace = ChainSpace{0, 0}
	g.MovePlayer(e.Player, e.Moves)
}
