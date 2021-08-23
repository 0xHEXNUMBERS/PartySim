package mp1

//wbcCannon sets the player's new ChainSpace.
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

//Responses returns a slice of possible positions the player can land on.
func (w wbcCannon) Responses() []Response {
	//TODO: Handle star spaces
	return wbcCannonDestinations[w.Chain]
}

func (w wbcCannon) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the player's new ChainSpace position.
func (w wbcCannon) Handle(r Response, g *Game) {
	space := r.(int)
	g.Players[w.Player].CurrentSpace = ChainSpace{w.Chain, space}
	g.MovePlayer(w.Player, w.Moves)
}

//wbcBowserCannon set's the player's new chain.
type wbcBowserCannon struct {
	Player int
	Moves  int
}

//Responses returns a slice of ints from [0, 4].
func (w wbcBowserCannon) Responses() []Response {
	return CPURangeEvent{0, 3}.Responses()
}

func (w wbcBowserCannon) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the player's chain to r, and sets the next event to
//selecting the player's new space.
func (w wbcBowserCannon) Handle(r Response, g *Game) {
	chain := r.(int)
	g.NextEvent = wbcCannon{w.Player, w.Moves, chain}
}

//wbcShyGuyResponse is a possible response to the shyguy action.
type wbcShyGuyResponse struct {
	Action wbcShyGuyAction
	Player int
}

//wbcShyGuyAction is an enumeration of possible shyguy actions.
type wbcShyGuyAction int

const (
	wbcNothing wbcShyGuyAction = iota
	wbcFlyToBowser
	wbcBringPlayer
)

//wbcShyGuyEvent let's the player decide on what to do when passing by
//shyguy.
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

//Responses returns the available responses a player can take.
func (w wbcShyGuyEvent) Responses() []Response {
	return wbcShyGuyResponses[w.Player]
}

func (w wbcShyGuyEvent) ControllingPlayer() int {
	return w.Player
}

//Handle executes the response r.
func (w wbcShyGuyEvent) Handle(r Response, g *Game) {
	res := r.(wbcShyGuyResponse)
	switch res.Action {
	case wbcFlyToBowser:
		g.NextEvent = wbcCannon{w.Player, w.Moves, 4}
	case wbcBringPlayer:
		g.Players[res.Player].CurrentSpace = ChainSpace{3, 4}
		g.MovePlayer(w.Player, w.Moves)
	default:
		g.MovePlayer(w.Player, w.Moves)
	}
}
