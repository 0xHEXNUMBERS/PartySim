package board

import (
	"fmt"

	"github.com/0xhexnumbers/partysim/mp1"
)

type DKJAWhompResponse int

const (
	DKJAWhompPay DKJAWhompResponse = iota
	DKJAWhompIgnore
)

func (d DKJAWhompResponse) String() string {
	switch d {
	case DKJAWhompPay:
		return "Pay 10 coins to the whomp"
	case DKJAWhompIgnore:
		return "Do not pay"
	}
	return ""
}

//DKJAWhompEvent let's the player decide to go and pay the whomp 10 coins
//or ignore the whomp.
type DKJAWhompEvent struct {
	Player int
	Moves  int
	Whomp  int
}

func (d DKJAWhompEvent) Question(g *mp1.Game) string {
	return fmt.Sprintf("Does %s pay 10 coins to get past the whomp?",
		g.Players[d.Player].Char)
}

func (d DKJAWhompEvent) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (d DKJAWhompEvent) ControllingPlayer() int {
	return d.Player
}

func (d DKJAWhompEvent) Responses() []mp1.Response {
	return []mp1.Response{DKJAWhompPay, DKJAWhompIgnore}
}

//Handle moves the player to the appropriate space and takes coins away
//based on r. If r is true, then the player loses 10 coins and moves past
//the whomp. If r is false, then the player moves away from the whomp.
func (d DKJAWhompEvent) Handle(r mp1.Response, g *mp1.Game) {
	pay := r.(DKJAWhompResponse)
	data := g.Board.Data.(dkjaBoardData)
	if pay == DKJAWhompPay {
		g.AwardCoins(d.Player, -10, false)
	} else {
		data.WhompPos[d.Whomp] = !data.WhompPos[d.Whomp]
		g.Board.Data = data
	}
	pos := dkjaGetWhompDestination(g, d.Whomp)
	g.Players[d.Player].CurrentSpace = pos
	g.MovePlayer(d.Player, d.Moves-1)
}
