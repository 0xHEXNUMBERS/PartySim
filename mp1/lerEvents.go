package mp1

//lerRedFork handles the RBR fork in the middle of the board when the red
//gates are down.
type lerRedFork struct {
	Player int
	Moves  int
}

var lerRedForkDestinations = []Response{
	ChainSpace{3, 0},
	ChainSpace{11, 0},
}

//Responses returns the 2 destinations from the fork.
func (l lerRedFork) Responses() []Response {
	return lerRedForkDestinations
}

func (l lerRedFork) ControllingPlayer() int {
	return l.Player
}

//Handle moves the player to the ChainSpace r.
func (l lerRedFork) Handle(r Response, g *Game) {
	dest := r.(ChainSpace)
	g.Players[l.Player].CurrentSpace = dest
	g.MovePlayer(l.Player, l.Moves-1)
}

//lerRobot let's the player decide to either pay and raise/lower gates or ignore the robot.
type lerRobot struct {
	Player int
	Moves  int
}

//Responses returns a slice of bools (true/false).
func (l lerRobot) Responses() []Response {
	return []Response{true, false}
}

func (l lerRobot) ControllingPlayer() int {
	return l.Player
}

//Handle pays the robot 20 coins and switches the gates only if r is true.
func (l lerRobot) Handle(r Response, g *Game) {
	pay := r.(bool)
	if pay {
		g.AwardCoins(l.Player, -20, false)
		lerSwapGates(g, l.Player)
	}
	g.MovePlayer(l.Player, l.Moves)
}
