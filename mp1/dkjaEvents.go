package mp1

//dkjaWhompEvent let's the player decide to go and pay the whomp 10 coins
//or ignore the whomp.
type dkjaWhompEvent struct {
	Player int
	Moves  int
	Whomp  int
}

//Responses return a slice of bools (true/false).
func (d dkjaWhompEvent) Responses() []Response {
	return []Response{true, false}
}

func (d dkjaWhompEvent) ControllingPlayer() int {
	return d.Player
}

//Handle moves the player to the appropriate space and takes coins away
//based on r. If r is true, then the player loses 10 coins and moves past
//the whomp. If r is false, then the player moves away from the whomp.
func (d dkjaWhompEvent) Handle(r Response, g *Game) {
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

//dkjaCoinBranchEvent let's the player decide to go through the coin
//blockade or not.
type dkjaCoinBranchEvent struct {
	Player   int
	Moves    int
	Blockade int
}

//Responses return a slice of bools (true/false).
func (d dkjaCoinBranchEvent) Responses() []Response {
	return []Response{true, false}
}

func (d dkjaCoinBranchEvent) ControllingPlayer() int {
	return d.Player
}

//Handle moves the player based on r. If r is true, the player moves
//through the coin blockade. If r is false, the player continues down the main path.
func (d dkjaCoinBranchEvent) Handle(r Response, g *Game) {
	accept := r.(bool)
	data := g.Board.Data.(dkjaBoardData)
	var pos ChainSpace
	if accept {
		pos = data.CoinAcceptDestination[d.Blockade]
	} else {
		pos = data.CoinRejectDestination[d.Blockade]
	}
	g.Players[d.Player].CurrentSpace = pos
	g.MovePlayer(d.Player, d.Moves-1)
}
