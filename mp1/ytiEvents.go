package mp1

//ytiThwompBranchEvent let's the player decide to go and pay the thwomp an
//amount of coins or ignore the thwomp.
type ytiThwompBranchEvent struct {
	Player int
	Moves  int
	Thwomp int
}

//Responses returns a slice of bools (true/false).
func (y ytiThwompBranchEvent) Responses() []Response {
	return []Response{true, false}
}

func (y ytiThwompBranchEvent) ControllingPlayer() int {
	return y.Player
}

//Handle calculates the next action based on r. If r is true, then the
//game's next event is set to pay the thwomp. If r is false, then the
//player moves to the Thwomp's rejection space and move their remaining
//spaces.
func (y ytiThwompBranchEvent) Handle(r Response, g *Game) {
	pay := r.(bool)
	bd := g.Board.Data.(ytiBoardData)
	if pay {
		g.ExtraEvent = ytiPayThwompEvent{
			PayRangeEvent{
				y.Player,
				bd.Thwomps[y.Thwomp],
				min(50, g.Players[y.Player].Coins),
			},
			y.Moves,
			y.Thwomp,
		}
	} else {
		pos := bd.RejectThwompPos[y.Thwomp]
		g.Players[y.Player].CurrentSpace = pos
		g.MovePlayer(y.Player, y.Moves-1)
	}
}

//ytiPayThwompEvent let's the player pay some amount of coins to the thwomp.
type ytiPayThwompEvent struct {
	PayRangeEvent
	Moves  int
	Thwomp int
}

//Handle pays the thwomp r coins, sets the thwomp's new asking price to r+1
//movess the player to the Thwomp's accept space, and the player moves
//their remaining spaces.
func (y ytiPayThwompEvent) Handle(r Response, g *Game) {
	y.PayRangeEvent.Handle(r, g)
	cost := r.(int)
	bd := g.Board.Data.(ytiBoardData)
	bd.Thwomps[y.Thwomp] = min(50, cost+1)
	pos := bd.AcceptThwompPos[y.Thwomp]
	g.Board.Data = bd
	g.Players[y.Player].CurrentSpace = pos
	g.MovePlayer(y.Player, y.Moves-1)
}
