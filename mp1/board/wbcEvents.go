package board

import (
	"fmt"

	"github.com/0xhexnumbers/partysim/mp1"
)

//WBCCannon sets the player's new mp1.ChainSpace.
type WBCCannon struct {
	Player int
	Moves  int
	Chain  int
}

var wbcCannonDestinationsPerChain = [5][]mp1.Response{
	wbcCannonDestinations[:13],
	wbcCannonDestinations[13:32],
	wbcCannonDestinations[32:50],
	wbcCannonDestinations[50:62],
	wbcCannonDestinations[62:],
}

var wbcCannonDestinations = []mp1.Response{
	mp1.NewChainSpace(0, 0),
	mp1.NewChainSpace(0, 1),
	mp1.NewChainSpace(0, 2),
	mp1.NewChainSpace(0, 3),
	mp1.NewChainSpace(0, 4),
	mp1.NewChainSpace(0, 5),
	mp1.NewChainSpace(0, 6),
	mp1.NewChainSpace(0, 7),
	mp1.NewChainSpace(0, 8),
	mp1.NewChainSpace(0, 10),
	mp1.NewChainSpace(0, 11),
	mp1.NewChainSpace(0, 12),
	mp1.NewChainSpace(0, 13),
	mp1.NewChainSpace(1, 0),
	mp1.NewChainSpace(1, 1),
	mp1.NewChainSpace(1, 2),
	mp1.NewChainSpace(1, 3),
	mp1.NewChainSpace(1, 4),
	mp1.NewChainSpace(1, 5),
	mp1.NewChainSpace(1, 6),
	mp1.NewChainSpace(1, 7),
	mp1.NewChainSpace(1, 8),
	mp1.NewChainSpace(1, 9),
	mp1.NewChainSpace(1, 10),
	mp1.NewChainSpace(1, 11),
	mp1.NewChainSpace(1, 12),
	mp1.NewChainSpace(1, 13),
	mp1.NewChainSpace(1, 14),
	mp1.NewChainSpace(1, 15),
	mp1.NewChainSpace(1, 16),
	mp1.NewChainSpace(1, 17),
	mp1.NewChainSpace(1, 18),
	mp1.NewChainSpace(2, 0),
	mp1.NewChainSpace(2, 1),
	mp1.NewChainSpace(2, 2),
	mp1.NewChainSpace(2, 3),
	mp1.NewChainSpace(2, 4),
	mp1.NewChainSpace(2, 5),
	mp1.NewChainSpace(2, 6),
	mp1.NewChainSpace(2, 7),
	mp1.NewChainSpace(2, 8),
	mp1.NewChainSpace(2, 9),
	mp1.NewChainSpace(2, 10),
	mp1.NewChainSpace(2, 11),
	mp1.NewChainSpace(2, 12),
	mp1.NewChainSpace(2, 13),
	mp1.NewChainSpace(2, 14),
	mp1.NewChainSpace(2, 15),
	mp1.NewChainSpace(2, 16),
	mp1.NewChainSpace(2, 17),
	mp1.NewChainSpace(3, 0),
	mp1.NewChainSpace(3, 1),
	mp1.NewChainSpace(3, 2),
	mp1.NewChainSpace(3, 3),
	mp1.NewChainSpace(3, 4),
	mp1.NewChainSpace(3, 6),
	mp1.NewChainSpace(3, 7),
	mp1.NewChainSpace(3, 8),
	mp1.NewChainSpace(3, 9),
	mp1.NewChainSpace(3, 10),
	mp1.NewChainSpace(3, 11),
	mp1.NewChainSpace(3, 12),
	mp1.NewChainSpace(4, 0),
	mp1.NewChainSpace(4, 1),
	mp1.NewChainSpace(4, 2),
	mp1.NewChainSpace(4, 3),
	mp1.NewChainSpace(4, 4),
	mp1.NewChainSpace(4, 5),
	mp1.NewChainSpace(4, 6),
}

func (w WBCCannon) Question(g *mp1.Game) string {
	return fmt.Sprintf("What space did %s land on?",
		g.Players[w.Player].Char)
}

func (w WBCCannon) Type() mp1.EventType {
	return mp1.CHAINSPACE_EVT_TYPE
}

//Responses returns a slice of possible positions the player can land on.
func (w WBCCannon) Responses() []mp1.Response {
	//TODO: Handle star spaces
	return wbcCannonDestinationsPerChain[w.Chain]
}

func (w WBCCannon) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle sets the player's new mp1.ChainSpace position.
func (w WBCCannon) Handle(r mp1.Response, g *mp1.Game) {
	space := r.(mp1.ChainSpace)
	g.Players[w.Player].CurrentSpace = space
	g.MovePlayer(w.Player, w.Moves)
}

//WBCBowserCannon set's the player's new position after visiting Bowser.
type WBCBowserCannon struct {
	Player int
	Moves  int
}

func (w WBCBowserCannon) Question(g *mp1.Game) string {
	return fmt.Sprintf("Which space did %s land on?",
		g.Players[w.Player].Char)
}

func (w WBCBowserCannon) Type() mp1.EventType {
	return mp1.CHAINSPACE_EVT_TYPE
}

//Responses returns a slice of ints from [0, 4].
func (w WBCBowserCannon) Responses() []mp1.Response {
	return wbcCannonDestinations[:62]
}

func (w WBCBowserCannon) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle sets the player's chain to r, and sets the next event to
//selecting the player's new space.
func (w WBCBowserCannon) Handle(r mp1.Response, g *mp1.Game) {
	pos := r.(mp1.ChainSpace)
	g.Players[w.Player].CurrentSpace = pos
	g.MovePlayer(w.Player, w.Moves)
}

//WBCShyGuyResponse is a possible response to the shyguy action.
type WBCShyGuyResponse struct {
	Action WBCShyGuyAction
	Player int
}

func (w WBCShyGuyResponse) String() string {
	switch w.Action {
	case WBCNothing:
		return "Do Nothing"
	case WBCFlyToBowser:
		return "Fly To Bowser"
	case WBCBringPlayer:
		return fmt.Sprintf("Bring Player %d to Shy Guy", w.Player+1)
	}
	return ""
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

func (w WBCShyGuyEvent) Question(g *mp1.Game) string {
	return fmt.Sprintf("What does %s do with the Shy Guy?",
		g.Players[w.Player].Char)
}

func (w WBCShyGuyEvent) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
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
