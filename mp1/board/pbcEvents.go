package board

import (
	"fmt"

	"github.com/0xhexnumbers/partysim/mp1"
)

type PBCSeedCheckResponse int

const (
	PBCSeedCheckBowser PBCSeedCheckResponse = iota
	PBCSeedCheckToad
)

func (p PBCSeedCheckResponse) String() string {
	switch p {
	case PBCSeedCheckBowser:
		return "Collect Bowser seed"
	case PBCSeedCheckToad:
		return "Collect Toad seed"
	}
	return ""
}

//PBCSeedCheck decides if the player got a toad seed or a bowser seed.
type PBCSeedCheck struct {
	Player int
	Moves  int
}

func (p PBCSeedCheck) Question(g *mp1.Game) string {
	return fmt.Sprintf("What seed did %s collect?",
		g.Players[p.Player].Char)
}

func (p PBCSeedCheck) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (p PBCSeedCheck) ControllingPlayer() int {
	return mp1.CPU_PLAYER
}

func (p PBCSeedCheck) Responses() []mp1.Response {
	return []mp1.Response{PBCSeedCheckBowser, PBCSeedCheckToad}
}

//Handle moves the player based on r. If r is true, the player moves to the
//bowser path. If r is false, the player moves to the toad path.
func (p PBCSeedCheck) Handle(r mp1.Response, g *mp1.Game) {
	seed := r.(PBCSeedCheckResponse)
	if seed == PBCSeedCheckBowser {
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
type PBCPiranhaDecisionResponse int

const (
	PBCPiranhaDecisionPay PBCPiranhaDecisionResponse = iota
	PBCPiranhaDecisionIgnore
)

func (p PBCPiranhaDecisionResponse) String() string {
	switch p {
	case PBCPiranhaDecisionPay:
		return "Pay 30 coins to plant Piranha"
	case PBCPiranhaDecisionIgnore:
		return "Do not pay"
	}
	return ""
}

//piranha.
type PBCPiranhaDecision struct {
	Player  int
	Piranha int
}

func (p PBCPiranhaDecision) Question(g *mp1.Game) string {
	return fmt.Sprintf("Does %s plant a Piranha seed?",
		g.Players[p.Player].Char)
}

func (p PBCPiranhaDecision) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (p PBCPiranhaDecision) ControllingPlayer() int {
	return p.Player
}

func (p PBCPiranhaDecision) Responses() []mp1.Response {
	return []mp1.Response{PBCPiranhaDecisionPay, PBCPiranhaDecisionIgnore}
}

//Handle performs the decision r. If r is true, then the player pays 30
//coins and gains a piranha at their current space. If r is false, nothing
//happens.
func (p PBCPiranhaDecision) Handle(r mp1.Response, g *mp1.Game) {
	plantPiranha := r.(PBCPiranhaDecisionResponse)
	data := g.Board.Data.(pbcBoardData)
	if plantPiranha == PBCPiranhaDecisionPay {
		data.PiranhaPlant[p.Piranha] = p.Player
		data.PiranhaOccupied[p.Piranha] = true
		g.AwardCoins(p.Player, -30, false)
	}
	g.Board.Data = data
	g.EndCharacterTurn()
}
