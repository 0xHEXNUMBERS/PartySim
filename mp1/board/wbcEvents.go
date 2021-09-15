package board

import "github.com/0xhexnumbers/partysim/mp1"

//WBCCannon sets the player's new mp1.ChainSpace.
type WBCCannon struct {
	Player int
	Moves  int
	Chain  int
}

var wbcCannonDestinations = [5][]mp1.Response{
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13},
	mp1.NewRange(0, 18),
	mp1.NewRange(0, 17),
	{0, 1, 2, 3, 4, 6, 7, 8, 9, 10, 11, 12},
	mp1.NewRange(0, 6),
}

//Responses returns a slice of possible positions the player can land on.
func (w WBCCannon) Responses() []mp1.Response {
	//TODO: Handle star spaces
	return wbcCannonDestinations[w.Chain]
}

func (w WBCCannon) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle sets the player's new mp1.ChainSpace position.
func (w WBCCannon) Handle(r mp1.Response, g *mp1.Game) {
	space := r.(int)
	g.Players[w.Player].CurrentSpace = mp1.NewChainSpace(w.Chain, space)
	g.MovePlayer(w.Player, w.Moves)
}

//WBCBowserCannon set's the player's new chain.
type WBCBowserCannon struct {
	Player int
	Moves  int
}

//Responses returns a slice of ints from [0, 4].
func (w WBCBowserCannon) Responses() []mp1.Response {
	return mp1.NewRange(0, 3)
}

func (w WBCBowserCannon) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle sets the player's chain to r, and sets the next event to
//selecting the player's new space.
func (w WBCBowserCannon) Handle(r mp1.Response, g *mp1.Game) {
	chain := r.(int)
	g.NextEvent = WBCCannon{w.Player, w.Moves, chain}
}

//WBCShyGuyResponse is a possible response to the shyguy action.
type WBCShyGuyResponse struct {
	Action WBCShyGuyAction
	Player int
}

//WBCShyGuyAction is an enumeration of possible shyguy actions.
type WBCShyGuyAction int

const (
	//WBCNothing represents ignoring the Shy Guy.
	WBCNothing WBCShyGuyAction = iota
	//WBCFlyToBowser takes the player to Bowser.
	WBCFlyToBowser
	//WBCBringPlayer represents the action of moving a player to Shy Guy.
	WBCBringPlayer
)

//WBCShyGuyEvent let's the player decide on what to do when passing by
//shyguy.
type WBCShyGuyEvent struct {
	Player int
	Moves  int
}

var wbcShyGuyResponses = [4][]mp1.Response{
	{
		WBCShyGuyResponse{WBCNothing, 0},
		WBCShyGuyResponse{WBCFlyToBowser, 0},
		WBCShyGuyResponse{WBCBringPlayer, 1},
		WBCShyGuyResponse{WBCBringPlayer, 2},
		WBCShyGuyResponse{WBCBringPlayer, 3},
	},
	{
		WBCShyGuyResponse{WBCNothing, 0},
		WBCShyGuyResponse{WBCFlyToBowser, 0},
		WBCShyGuyResponse{WBCBringPlayer, 0},
		WBCShyGuyResponse{WBCBringPlayer, 2},
		WBCShyGuyResponse{WBCBringPlayer, 3},
	},
	{
		WBCShyGuyResponse{WBCNothing, 0},
		WBCShyGuyResponse{WBCFlyToBowser, 0},
		WBCShyGuyResponse{WBCBringPlayer, 0},
		WBCShyGuyResponse{WBCBringPlayer, 1},
		WBCShyGuyResponse{WBCBringPlayer, 3},
	},
	{
		WBCShyGuyResponse{WBCNothing, 0},
		WBCShyGuyResponse{WBCFlyToBowser, 0},
		WBCShyGuyResponse{WBCBringPlayer, 0},
		WBCShyGuyResponse{WBCBringPlayer, 1},
		WBCShyGuyResponse{WBCBringPlayer, 2},
	},
}

//Responses returns the available responses a player can take.
func (w WBCShyGuyEvent) Responses() []mp1.Response {
	return wbcShyGuyResponses[w.Player]
}

func (w WBCShyGuyEvent) ControllingPlayer() int {
	return w.Player
}

//Handle executes the response r.
func (w WBCShyGuyEvent) Handle(r mp1.Response, g *mp1.Game) {
	res := r.(WBCShyGuyResponse)
	switch res.Action {
	case WBCFlyToBowser:
		g.NextEvent = WBCCannon{w.Player, w.Moves, 4}
	case WBCBringPlayer:
		g.Players[res.Player].CurrentSpace = mp1.NewChainSpace(3, 4)
		g.MovePlayer(w.Player, w.Moves)
	default:
		g.MovePlayer(w.Player, w.Moves)
	}
}
