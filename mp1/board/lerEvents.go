package board

import (
	"fmt"

	"github.com/0xhexnumbers/partysim/mp1"
)

type LERRobotResponse int

const (
	LERRobotPay LERRobotResponse = iota
	LERRobotIgnore
)

func (l LERRobotResponse) String() string {
	switch l {
	case LERRobotPay:
		return "Pay 20 coins to flip switch"
	case LERRobotIgnore:
		return "Ignore robot"
	}
	return ""
}

//LERRobot let's the player decide to either pay and raise/lower gates or ignore the robot.
type LERRobot struct {
	Player int
	Moves  int
}

func (l LERRobot) Question(g *mp1.Game) string {
	return fmt.Sprintf(
		"Does %s pay 20 coins to flip the Red/Blue Switch?",
		g.Players[l.Player].Char,
	)
}

func (l LERRobot) Type() mp1.EventType {
	return mp1.ENUM_EVT_TYPE
}

func (l LERRobot) ControllingPlayer() int {
	return l.Player
}

func (l LERRobot) Responses() []mp1.Response {
	return []mp1.Response{LERRobotPay, LERRobotIgnore}
}

//Handle pays the robot 20 coins and switches the gates only if r is true.
func (l LERRobot) Handle(r mp1.Response, g *mp1.Game) {
	pay := r.(LERRobotResponse)
	if pay == LERRobotPay {
		g.AwardCoins(l.Player, -20, false)
		lerSwapGates(g, l.Player)
	}
	g.MovePlayer(l.Player, l.Moves)
}
