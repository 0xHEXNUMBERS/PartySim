package board

import "github.com/0xhexnumbers/partysim/mp1"

//PBCSeedCheck decides if the player got a toad seed or a bowser seed.
type PBCSeedCheck struct {
	Player int
	Moves  int
}

//Responses returns a slice of bools (true/false).
func (p PBCSeedCheck) Responses() []mp1.Response {
	return []mp1.Response{true, false}
}

func (p PBCSeedCheck) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

//Handle moves the player based on r. If r is true, the player moves to the
//bowser path. If r is false, the player moves to the toad path.
func (p PBCSeedCheck) Handle(r mp1.Response, g *mp1.Game) {
	bowser := r.(bool)
	if bowser {
		g.Players[p.Player].CurrentSpace = mp1.NewChainSpace(1, 0)
		data := g.Board.Data.(pbcBoardData)
		data.BowserSeedPlanted = true
		g.Board.Data = data
	} else {
		g.Players[p.Player].CurrentSpace = mp1.NewChainSpace(0, 0)
	}
	g.MovePlayer(p.Player, p.Moves-1)
}

//PBCPiranhaDecision decides if the player wants to pay 30 coins for a
//piranha.
type PBCPiranhaDecision struct {
	Player  int
	Piranha int
}

//Responses returns a slice of bools (true/false).
func (p PBCPiranhaDecision) Responses() []mp1.Response {
	return []mp1.Response{true, false}
}

func (p PBCPiranhaDecision) ControllingPlayer() int {
	return p.Player
}

//Handle performs the decision r. If r is true, then the player pays 30
//coins and gains a piranha at their current space. If r is false, nothing
//happens.
func (p PBCPiranhaDecision) Handle(r mp1.Response, g *mp1.Game) {
	plantPiranha := r.(bool)
	data := g.Board.Data.(pbcBoardData)
	if plantPiranha {
		data.PiranhaPlant[p.Piranha] = p.Player
		data.PiranhaOccupied[p.Piranha] = true
		g.AwardCoins(p.Player, -30, false)
	}
	g.Board.Data = data
	g.EndCharacterTurn()
}
