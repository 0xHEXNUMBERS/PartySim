package board

import "github.com/0xhexnumbers/partysim/mp1"

//DKJAWhompEvent let's the player decide to go and pay the whomp 10 coins
//or ignore the whomp.
type DKJAWhompEvent struct {
	Player int
	Moves  int
	Whomp  int
}

//mp1.Responses return a slice of bools (true/false).
func (d DKJAWhompEvent) Responses() []mp1.Response {
	return []mp1.Response{true, false}
}

func (d DKJAWhompEvent) ControllingPlayer() int {
	return d.Player
}

//Handle moves the player to the appropriate space and takes coins away
//based on r. If r is true, then the player loses 10 coins and moves past
//the whomp. If r is false, then the player moves away from the whomp.
func (d DKJAWhompEvent) Handle(r mp1.Response, g *mp1.Game) {
	pay := r.(bool)
	data := g.Board.Data.(dkjaBoardData)
	if pay {
		g.AwardCoins(d.Player, -10, false)
	} else {
		data.WhompPos[d.Whomp] = !data.WhompPos[d.Whomp]
		g.Board.Data = data
	}
	pos := dkjaGetWhompDestination(g, d.Whomp)
	g.Players[d.Player].CurrentSpace = pos
	g.MovePlayer(d.Player, d.Moves-1)
}

//DKJACoinBranchEvent let's the player decide to go through the coin
//blockade or not.
type DKJACoinBranchEvent struct {
	Player   int
	Moves    int
	Blockade int
}

//mp1.Responses return a slice of bools (true/false).
func (d DKJACoinBranchEvent) Responses() []mp1.Response {
	return []mp1.Response{true, false}
}

func (d DKJACoinBranchEvent) ControllingPlayer() int {
	return d.Player
}

//Handle moves the player based on r. If r is true, the player moves
//through the coin blockade. If r is false, the player continues down the main path.
func (d DKJACoinBranchEvent) Handle(r mp1.Response, g *mp1.Game) {
	accept := r.(bool)
	data := g.Board.Data.(dkjaBoardData)
	var pos mp1.ChainSpace
	if accept {
		pos = data.CoinAcceptDestination[d.Blockade]
	} else {
		pos = data.CoinRejectDestination[d.Blockade]
	}
	g.Players[d.Player].CurrentSpace = pos
	g.MovePlayer(d.Player, d.Moves-1)
}
