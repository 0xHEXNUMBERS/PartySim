package board

import "github.com/0xhexnumbers/partysim/mp1"

//LERRedFork handles the RBR fork in the middle of the board when the red
//gates are down.
type LERRedFork struct {
	Player int
	Moves  int
}

var lerRedForkDestinations = []mp1.Response{
	mp1.NewChainSpace(3, 0),
	mp1.NewChainSpace(11, 0),
}

//Responses returns the 2 destinations from the fork.
func (l LERRedFork) Responses() []mp1.Response {
	return lerRedForkDestinations
}

func (l LERRedFork) ControllingPlayer() int {
	return l.Player
}

//Handle moves the player to the mp1.ChainSpace r.
func (l LERRedFork) Handle(r mp1.Response, g *mp1.Game) {
	dest := r.(mp1.ChainSpace)
	g.Players[l.Player].CurrentSpace = dest
	g.MovePlayer(l.Player, l.Moves-1)
}

//LERRobot let's the player decide to either pay and raise/lower gates or ignore the robot.
type LERRobot struct {
	mp1.Boolean
	Player int
	Moves  int
}

func (l LERRobot) ControllingPlayer() int {
	return l.Player
}

//Handle pays the robot 20 coins and switches the gates only if r is true.
func (l LERRobot) Handle(r mp1.Response, g *mp1.Game) {
	pay := r.(bool)
	if pay {
		g.AwardCoins(l.Player, -20, false)
		lerSwapGates(g, l.Player)
	}
	g.MovePlayer(l.Player, l.Moves)
}
