package board

import (
	"fmt"

	"github.com/0xhexnumbers/partysim/mp1"
)

type YTIThwompBranchResponse int

const (
	YTIThwompBranchPay YTIThwompBranchResponse = iota
	YTIThwompBranchIgnore
)

func (y YTIThwompBranchResponse) String() string {
	switch y {
	case YTIThwompBranchPay:
		return "Pay to pass Thwomp"
	case YTIThwompBranchIgnore:
		return "Do not pay Thwomp"
	}
	return ""
}

//YTIThwompBranchEvent let's the player decide to go and pay the thwomp an
//amount of coins or ignore the thwomp.
type YTIThwompBranchEvent struct {
	Player int
	Moves  int
	Thwomp int
}

func (y YTIThwompBranchEvent) Question(g *mp1.Game) string {
	return fmt.Sprintf("Does %s pay to pass the Thwomp?",
		g.Players[y.Player].Char)
}

func (y YTIThwompBranchEvent) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (y YTIThwompBranchEvent) ControllingPlayer() int {
	return y.Player
}

func (y YTIThwompBranchEvent) Responses() []mp1.Response {
	return []mp1.Response{
		YTIThwompBranchPay,
		YTIThwompBranchIgnore,
	}
}

//Handle calculates the next action based on r. If r is true, then the
//game's next event is set to pay the thwomp. If r is false, then the
//player moves to the Thwomp's rejection space and move their remaining
//spaces.
func (y YTIThwompBranchEvent) Handle(r mp1.Response, g *mp1.Game) {
	pay := r.(YTIThwompBranchResponse)
	bd := g.Board.Data.(ytiBoardData)
	if pay == YTIThwompBranchPay {
		g.NextEvent = YTIPayThwompEvent{
			mp1.PayRangeEvent{
				Range: mp1.Range{
					Min: bd.Thwomps[y.Thwomp],
					Max: min(50, g.Players[y.Player].Coins),
				},
				Player: y.Player,
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

//YTIPayThwompEvent let's the player pay some amount of coins to the thwomp.
type YTIPayThwompEvent struct {
	mp1.PayRangeEvent
	Moves  int
	Thwomp int
}

func (y YTIPayThwompEvent) Question(g *mp1.Game) string {
	return fmt.Sprintf("How many coins does %s pay to pass the Thwomp?",
		g.Players[y.PayRangeEvent.Player].Char)
}

//Handle pays the thwomp r coins, sets the thwomp's new asking price to r+1
//movess the player to the Thwomp's accept space, and the player moves
//their remaining spaces.
func (y YTIPayThwompEvent) Handle(r mp1.Response, g *mp1.Game) {
	y.PayRangeEvent.Handle(r, g)
	cost := r.(int)
	bd := g.Board.Data.(ytiBoardData)
	bd.Thwomps[y.Thwomp] = min(50, cost+1)
	pos := bd.AcceptThwompPos[y.Thwomp]
	g.Board.Data = bd
	g.Players[y.Player].CurrentSpace = pos
	g.MovePlayer(y.Player, y.Moves-1)
}
