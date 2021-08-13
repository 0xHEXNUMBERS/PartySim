package mp1

type wbcCannon struct {
	Player int
	Moves  int
	Chain  int
}

var wbcCannonDestinations = [5][]Response{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13},
	CPURangeEvent{0, 18}.Responses(),
	CPURangeEvent{0, 17}.Responses(),
	{0, 1, 2, 3, 4, 6, 7, 8, 9, 10, 11, 12},
	CPURangeEvent{0, 6}.Responses(),
}

func (w wbcCannon) Responses() []Response {
	//TODO: Handle star spaces
	return wbcCannonDestinations[w.Chain]
}

func (w wbcCannon) ControllingPlayer() int {
	return CPU_PLAYER
}

func (w wbcCannon) Handle(r Response, g *Game) {
	space := r.(int)
	g.Players[w.Player].CurrentSpace = ChainSpace{w.Chain, space}
	g.MovePlayer(w.Player, w.Moves)
}

type wbcBowserCannon struct {
	Player int
	Moves  int
}

func (w wbcBowserCannon) Responses() []Response {
	return CPURangeEvent{0, 3}.Responses()
}

func (w wbcBowserCannon) ControllingPlayer() int {
	return CPU_PLAYER
}

func (w wbcBowserCannon) Handle(r Response, g *Game) {
	chain := r.(int)
	g.ExtraEvent = wbcCannon{w.Player, w.Moves, chain}
}

type wbcShyGuyResponse struct {
	Action wbcShyGuyAction
	Player int
}

type wbcShyGuyAction int

const (
	wbcNothing wbcShyGuyAction = iota
	wbcFlyToBowser
	wbcBringPlayer
)

type wbcShyGuyEvent struct {
	Player int
	Moves  int
}

var wbcShyGuyResponses = [4][]Response{
	{
		wbcShyGuyResponse{wbcNothing, 0},
		wbcShyGuyResponse{wbcFlyToBowser, 0},
		wbcShyGuyResponse{wbcBringPlayer, 1},
		wbcShyGuyResponse{wbcBringPlayer, 2},
		wbcShyGuyResponse{wbcBringPlayer, 3},
	},
	{
		wbcShyGuyResponse{wbcNothing, 0},
		wbcShyGuyResponse{wbcFlyToBowser, 0},
		wbcShyGuyResponse{wbcBringPlayer, 0},
		wbcShyGuyResponse{wbcBringPlayer, 2},
		wbcShyGuyResponse{wbcBringPlayer, 3},
	},
	{
		wbcShyGuyResponse{wbcNothing, 0},
		wbcShyGuyResponse{wbcFlyToBowser, 0},
		wbcShyGuyResponse{wbcBringPlayer, 0},
		wbcShyGuyResponse{wbcBringPlayer, 1},
		wbcShyGuyResponse{wbcBringPlayer, 3},
	},
	{
		wbcShyGuyResponse{wbcNothing, 0},
		wbcShyGuyResponse{wbcFlyToBowser, 0},
		wbcShyGuyResponse{wbcBringPlayer, 0},
		wbcShyGuyResponse{wbcBringPlayer, 1},
		wbcShyGuyResponse{wbcBringPlayer, 2},
	},
}

func (w wbcShyGuyEvent) Responses() []Response {
	return wbcShyGuyResponses[w.Player]
}

func (w wbcShyGuyEvent) ControllingPlayer() int {
	return w.Player
}

func (w wbcShyGuyEvent) Handle(r Response, g *Game) {
	res := r.(wbcShyGuyResponse)
	switch res.Action {
	case wbcFlyToBowser:
		g.ExtraEvent = wbcCannon{w.Player, w.Moves, 4}
	case wbcBringPlayer:
		g.Players[res.Player].CurrentSpace = ChainSpace{3, 4}
		g.MovePlayer(w.Player, w.Moves)
	default:
		g.MovePlayer(w.Player, w.Moves)
	}
}
