package mp1

//MinigameFFAReward handles Free-For-All minigame rewards. One player
//will gain coins from this event. It contians a CoinMinigame if extra
//coins can be gained by any player.
type MinigameFFAReward struct {
	IsCoinMinigame bool
	Coin           CoinMinigameFFAReward
}

//DrawableFFAReward handles drawable Free-For-All minigame rewards. One
//player may gain coins from this event. It contians a CoinMinigame if
//extra coins can be gained by any player.
type DrawableFFAReward struct {
	MinigameFFAReward
}

var DrawableFFAPlayers = []Response{0, 1, 2, 3, 4}

var MinigameFFAPlayers = DrawableFFAPlayers[:4]

//Responses returns a slice of ints from [0, 4]
func (d DrawableFFAReward) Responses() []Response {
	return DrawableFFAPlayers
}

//Responses returns a slice of ints from [0, 3]
func (m MinigameFFAReward) Responses() []Response {
	return MinigameFFAPlayers
}

func (m MinigameFFAReward) ControllingPlayer() int {
	return CPU_PLAYER
}

//MinigameFFAReward gives out 10 coins to player r. If r == 4, then no one
//gains coins. If m.IsCoinMinigame is true, then the game's next event is
//set to the containing coin minigame. Otherwise, the game's turn ends.
func (m MinigameFFAReward) Handle(r Response, g *Game) {
	player := r.(int)
	if player != 4 {
		g.AwardCoins(player, 10, true)
	}
	if m.IsCoinMinigame {
		g.ExtraEvent = m.Coin
	} else {
		g.EndGameTurn()
	}
}

//CoinMinigameFFAReward distributes coins gained from coin minigames.
type CoinMinigameFFAReward struct {
	Player int
	Max    int
}

//Responses returns a slice of ints from [0, c.Max].
func (c CoinMinigameFFAReward) Responses() []Response {
	return CPURangeEvent{0, c.Max}.Responses()
}

func (c CoinMinigameFFAReward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives r coins to c.Player. The games next event is either set to
//the same event with c.Player incremented, or ends the game's turn.
func (c CoinMinigameFFAReward) Handle(r Response, g *Game) {
	coins := r.(int)
	g.AwardCoins(c.Player, coins, true)
	if c.Player == 3 {
		g.EndGameTurn()
	} else {
		c.Player++
		g.ExtraEvent = c
	}
}

//MinigameFFAMultiWinReward handles Free-For-All minigame rewards. 0 or
//more players will gain coins from this event, and 0 or more players will
//lose coins from this event.
type MinigameFFAMultiWinReward struct {
	CoinsToWin      int
	CoinsToLose     int
	CoinsIfNoWinner int
}

//Responses returns a slice of ints from [0, 15].
func (m MinigameFFAMultiWinReward) Responses() []Response {
	return CPURangeEvent{0, 15}.Responses()
}

func (m MinigameFFAMultiWinReward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives out gives out coins to players based on the value of r. The
//bits 0-3 in r determine if player x won the minigame. For example: if r
//is 9, player 0 and player 2 both win coins, while player 1 and player 3
//lose coins.
func (m MinigameFFAMultiWinReward) Handle(r Response, g *Game) {
	wins := r.(int)
	defer g.EndGameTurn()
	if wins == 0 {
		for p := 0; p < 4; p++ {
			g.AwardCoins(p, m.CoinsIfNoWinner, true)
		}
		return
	}
	for p := 0; p < 4; p++ {
		if wins&(1<<p) > 0 {
			g.AwardCoins(p, m.CoinsToWin, true)
		} else {
			g.AwardCoins(p, m.CoinsToLose, true)
		}
	}
}

//MinigameFFA1Loser handles Free-For-All minigame rewards. Either 1 player
//gives the other 3 5 coins each, or all players win 10 coins.
type MinigameFFA1Loser struct{}

//Responses returns a slice of ints from [0, 4].
func (m MinigameFFA1Loser) Responses() []Response {
	return CPURangeEvent{0, 4}.Responses()
}

func (m MinigameFFA1Loser) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives coins to player based on r. If r == 4, then all players win
//10 coins. Otherwise, player r gives 5 coins to every other player.
func (m MinigameFFA1Loser) Handle(r Response, g *Game) {
	player := r.(int)
	defer g.EndGameTurn()
	if player == 4 {
		for i := 0; i < 4; i++ {
			g.AwardCoins(i, 10, true)
		}
		return
	}
	for i := 0; i < 4; i++ {
		if i == player {
			g.AwardCoins(i, -15, true)
		} else {
			g.AwardCoins(i, 5, true)
		}
	}
}

//MinigameFFACoop handles Free-For-All cooperative minigame rewards.
//Players either win 10 coins each or lose 5 coins each.
type MinigameFFACoop struct{}

//Responses returns a slice of bools (true/false).
func (m MinigameFFACoop) Responses() []Response {
	return []Response{true, false}
}

func (m MinigameFFACoop) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives coins to players based on r. If r == true, then each player
//gains 10 coins. Otherwise, each player loses 5 coins.
func (m MinigameFFACoop) Handle(r Response, g *Game) {
	won := r.(bool)
	coins := -5
	if won {
		coins = 10
	}
	for i := 0; i < 4; i++ {
		g.AwardCoins(i, coins, true)
	}
	g.EndGameTurn()
}

//MinigameGrabBag handles the grab bag FFA minigame. Players steal coins
//from each other in the duration of this minigame.
type MinigameGrabBag struct {
	Player   int
	Acc      int
	MaxCoins int
}

//Responses returns a slice of ints from [-m.MaxCoins,m.MaxCoins]
func (m MinigameGrabBag) Responses() []Response {
	return CPURangeEvent{-m.MaxCoins, m.MaxCoins}.Responses()
}

func (m MinigameGrabBag) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives player m.Player r coins. An accumulator is kept to save one
//round of event execution. If m.Player is 2, then Player 3 receives m.Acc
//coins. Otherwise, the next event is set to this event with m.Player
//incremented.
func (m MinigameGrabBag) Handle(r Response, g *Game) {
	coins := r.(int)
	m.Acc += coins
	g.AwardCoins(m.Player, coins, true)
	m.Player++
	if m.Player >= 3 {
		g.AwardCoins(3, m.Acc, true)
		g.EndGameTurn()
	} else {
		g.ExtraEvent = m
	}
}

//MinigameFFAGame is a enumeration of the available FFA minigames.
type MinigameFFAGame int

const (
	MinigameFFABurriedTreasure MinigameFFAGame = iota
	MinigameFFATreasureDivers
	MinigameFFAHotBobomb
	MinigameFFAMusicalMushroom
	MinigameFFACrazyCutter
	MinigameFFAFaceLift
	MinigameFFABalloonBurst
	MinigameFFACoinBlockBlitz
	MinigameFFASkateboardScamper
	MinigameFFABoxMountainMayhem
	MinigameFFAPlatformPeril
	MinigameFFAMushroomMixup
	MinigameFFAGrabBag
	MinigameFFABumperBalls
	MinigameFFATipsyTourney
	MinigameFFABombsAway
	MinigameFFAMarioBandstand
	MinigameFFAShyGuySays
	MinigameFFACastAways
	MinigameFFAKeypaWay
	MinigameFFARunningoftheBulb
	MinigameFFAHotRopeJump
	MinigameFFAHammerDrop
	MinigameFFASlotCarDerby
)

//MinigameFFASelector selects which FFA minigame to play.
type MinigameFFASelector struct{}

//Responses returns a slice of all MinigameFFAGame enumerations.
func (m MinigameFFASelector) Responses() []Response {
	return []Response{
		MinigameFFABurriedTreasure,
		MinigameFFATreasureDivers,
		MinigameFFAHotBobomb,
		MinigameFFAMusicalMushroom,
		MinigameFFACrazyCutter,
		MinigameFFAFaceLift,
		MinigameFFABalloonBurst,
		MinigameFFACoinBlockBlitz,
		MinigameFFASkateboardScamper,
		MinigameFFABoxMountainMayhem,
		MinigameFFAPlatformPeril,
		MinigameFFAMushroomMixup,
		MinigameFFAGrabBag,
		MinigameFFABumperBalls,
		MinigameFFATipsyTourney,
		MinigameFFABombsAway,
		MinigameFFAMarioBandstand,
		MinigameFFAShyGuySays,
		MinigameFFACastAways,
		MinigameFFAKeypaWay,
		MinigameFFARunningoftheBulb,
		MinigameFFAHotRopeJump,
		MinigameFFAHammerDrop,
		MinigameFFASlotCarDerby,
	}
}

func (m MinigameFFASelector) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the next event to the selected Minigame.
func (m MinigameFFASelector) Handle(r Response, g *Game) {
	game := r.(MinigameFFAGame)
	switch game {
	case MinigameFFABurriedTreasure:
		g.ExtraEvent = MinigameFFAReward{}
	case MinigameFFATreasureDivers:
		g.ExtraEvent = CoinMinigameFFAReward{0, 50}
	case MinigameFFAHotBobomb:
		g.ExtraEvent = MinigameFFA1Loser{}
	case MinigameFFAMusicalMushroom:
		g.ExtraEvent = MinigameFFAReward{}
	case MinigameFFACrazyCutter:
		g.ExtraEvent = MinigameFFAMultiWinReward{10, -5, 0}
	case MinigameFFAFaceLift:
		g.ExtraEvent = MinigameFFAMultiWinReward{10, -5, 0}
	case MinigameFFABalloonBurst:
		g.ExtraEvent = MinigameFFAReward{}
	case MinigameFFACoinBlockBlitz:
		g.ExtraEvent = CoinMinigameFFAReward{0, 40}
	case MinigameFFASkateboardScamper:
		//TODO: Separate coin from coinbag
		g.ExtraEvent = MinigameFFAReward{true, CoinMinigameFFAReward{0, 10}}
	case MinigameFFABoxMountainMayhem:
		g.ExtraEvent = CoinMinigameFFAReward{0, 25}
	case MinigameFFAPlatformPeril:
		//TODO: Separate coin from coinbag
		g.ExtraEvent = MinigameFFAReward{true, CoinMinigameFFAReward{0, 10}}
	case MinigameFFAMushroomMixup:
		g.ExtraEvent = DrawableFFAReward{}
	case MinigameFFAGrabBag:
		//TODO: I'm not sure how I feel about max of 50.
		//Theoritically, players can steal > 50 coins, but probably not
		//feasible
		g.ExtraEvent = MinigameGrabBag{0, 0, 50}
	case MinigameFFABumperBalls:
		g.ExtraEvent = DrawableFFAReward{}
	case MinigameFFATipsyTourney:
		g.ExtraEvent = MinigameFFAReward{}
	case MinigameFFABombsAway:
		g.ExtraEvent = MinigameFFAReward{}
	case MinigameFFAMarioBandstand:
		//TODO: Find out how many coins are distributed
		g.ExtraEvent = MinigameFFAMultiWinReward{}
	case MinigameFFAShyGuySays:
		g.ExtraEvent = MinigameFFAReward{}
	case MinigameFFACastAways:
		//TODO: Should optimize; ask for chests/bags/coins
		g.ExtraEvent = CoinMinigameFFAReward{0, 80}
	case MinigameFFAKeypaWay:
		g.ExtraEvent = MinigameFFACoop{}
	case MinigameFFARunningoftheBulb:
		g.ExtraEvent = MinigameFFAMultiWinReward{10, 0, -5}
	case MinigameFFAHotRopeJump:
		g.ExtraEvent = MinigameFFA1Loser{}
	case MinigameFFAHammerDrop:
		g.ExtraEvent = CoinMinigameFFAReward{0, 20}
	case MinigameFFASlotCarDerby:
		g.ExtraEvent = MinigameFFAReward{}
	}
}

//Minigame2V2Reward handles 2v2 minigame rewards. One team will gain coins
//from this event, while the other team will lose coins.
type Minigame2V2Reward struct {
	BlueTeam [2]int
	RedTeam  [2]int
}

//Drawable2V2Reward handles 2v2 minigame rewards. Zero or One team will
//gain coins from this event, while the other team may lose coins.
type DrawableMinigame2V2Reward Minigame2V2Reward

var Drawable2V2Players = []Response{0, 1, 2}
var Minigame2V2Players = Drawable2V2Players[:2]

//Responses returns a slice of ints from [0, 2]
func (d DrawableMinigame2V2Reward) Responses() []Response {
	return Drawable2V2Players
}

//Responses returns a slice of ints from [0, 1]
func (m Minigame2V2Reward) Responses() []Response {
	return Minigame2V2Players
}

func (m Minigame2V2Reward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives coins out to team members based on r. If r == 0, then the
//blue team takes coins from the red team. If r == 1, then the red team
//takes coins from the blue team. Otherwise, it is considered a draw.
func (m Minigame2V2Reward) Handle(r Response, g *Game) {
	team := r.(int)
	if team == 0 {
		g.AwardCoins(m.BlueTeam[0], 10, true)
		g.AwardCoins(m.BlueTeam[1], 10, true)
		g.AwardCoins(m.RedTeam[0], -10, true)
		g.AwardCoins(m.RedTeam[1], -10, true)
	} else if team == 1 {
		g.AwardCoins(m.RedTeam[0], 10, true)
		g.AwardCoins(m.RedTeam[1], 10, true)
		g.AwardCoins(m.BlueTeam[0], -10, true)
		g.AwardCoins(m.BlueTeam[1], -10, true)
	}
	g.EndGameTurn()
}

//CoinMinigame2V2Reward distributes coins gained from 2v2 coin minigames.
type CoinMinigame2V2Reward struct {
	BlueTeam [2]int
	RedTeam  [2]int
	Team     int
	Max      int
}

//Responses returns a slice of ints from [0,c.Max].
func (c CoinMinigame2V2Reward) Responses() []Response {
	return CPURangeEvent{0, c.Max}.Responses()
}

func (c CoinMinigame2V2Reward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives r coins to team c.Team (0 == Blue, 1 == Red).
func (c CoinMinigame2V2Reward) Handle(r Response, g *Game) {
	coins := r.(int)
	if c.Team == 0 {
		g.AwardCoins(c.BlueTeam[0], coins, true)
		g.AwardCoins(c.BlueTeam[1], coins, true)
	} else if c.Team == 1 {
		g.AwardCoins(c.RedTeam[0], coins, true)
		g.AwardCoins(c.RedTeam[1], coins, true)
	} else if c.Team == 2 {
		g.EndGameTurn()
		return
	}
	c.Team++
	g.ExtraEvent = c
}

//Minigame2V2Game is a enumeration of the available 2V2 minigames.
type Minigame2V2Game int

const (
	Minigame2V2BobsledRun Minigame2V2Game = iota
	Minigame2V2DesertDash
	Minigame2V2Bombsketball
	Minigame2V2HandcarHavoc
	Minigame2V2DeepSeaDivers
)

//Minigame2V2Selector selects which 2V2 minigame to play.
type Minigame2V2Selector struct {
	Team1 [2]int
	Team2 [2]int
}

//Responses returns a slice of all Minigame2V2Game enumerations.
func (m Minigame2V2Selector) Responses() []Response {
	return []Response{
		Minigame2V2BobsledRun,
		Minigame2V2DesertDash,
		Minigame2V2Bombsketball,
		Minigame2V2HandcarHavoc,
		Minigame2V2DeepSeaDivers,
	}
}

func (m Minigame2V2Selector) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the next event to the selected Minigame.
func (m Minigame2V2Selector) Handle(r Response, g *Game) {
	game := r.(Minigame2V2Game)
	switch game {
	case Minigame2V2BobsledRun:
		g.ExtraEvent = Minigame2V2Reward{
			m.Team1, m.Team2,
		}
	case Minigame2V2DesertDash:
		g.ExtraEvent = Minigame2V2Reward{
			m.Team1, m.Team2,
		}
	case Minigame2V2Bombsketball:
		g.ExtraEvent = Minigame2V2Reward{
			m.Team1, m.Team2,
		}
	case Minigame2V2HandcarHavoc:
		g.ExtraEvent = Minigame2V2Reward{
			m.Team1, m.Team2,
		}
	case Minigame2V2DeepSeaDivers:
		g.ExtraEvent = CoinMinigame2V2Reward{
			m.Team1, m.Team2, 0, 50,
		}
	}
}

//Minigame1V3Reward handles 1v3 minigame rewards. One team will gain coins
//from this event, while the other team will lose coins.
type Minigame1V3Reward struct {
	SingleTeam int
}

//Drawable1V3Reward handles 1v3 minigame rewards. Zero or One team will
//gain coins from this event, while the other team may lose coins.
type Drawable1V3Reward Minigame1V3Reward

var Drawable1V3Players = []Response{0, 1, 2}
var Minigame1V3Players = Drawable1V3Players[:2]

//Responses returns a slice of ints from [0, 2]
func (d Drawable1V3Reward) Responses() []Response {
	return Minigame1V3Players
}

//Responses returns a slice of ints from [0, 1]
func (m Minigame1V3Reward) Responses() []Response {
	return Minigame1V3Players
}

func (m Minigame1V3Reward) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives coins out to team members based on r. If r == 0, then the
//solo player takes coins from the 3player team. If r == 1, then the
//3player takes coins from the solo player. Otherwise, it is considered a
//draw.
func (m Minigame1V3Reward) Handle(r Response, g *Game) {
	team := r.(int)
	if team == 0 {
		g.AwardCoins(m.SingleTeam, 15, true)
		for i := 0; i < 4; i++ {
			if i != m.SingleTeam {
				g.AwardCoins(i, -5, true)
			}
		}
	} else if team == 1 {
		g.AwardCoins(m.SingleTeam, -15, true)
		for i := 0; i < 4; i++ {
			if i != m.SingleTeam {
				g.AwardCoins(i, 5, true)
			}
		}
	}
	g.EndGameTurn()
}

//Throwable1V3Minigame is a minigame that the Solo player may choose to
//lose, causing no one to gain coins.
type Throwable1V3Minigame struct {
	Player   int
	Minigame Event
}

//Responses returns a slice of bools (true/false).
func (t Throwable1V3Minigame) Responses() []Response {
	return []Response{false, true}
}

func (t Throwable1V3Minigame) ControllingPlayer() int {
	return t.Player
}

//Handle chooses whether to throw the minigame based on r. If r == true,
//then the game turn ends. If r == false, the next event becomes the
//underlying minigame.
func (t Throwable1V3Minigame) Handle(r Response, g *Game) {
	throw := r.(bool)
	if throw {
		g.EndGameTurn()
	} else {
		g.ExtraEvent = t.Minigame
	}
}

//Specific 1V3 Minigames

//MinigamePipeMaze holds the implementation for Pipe Maze.
type MinigamePipeMaze struct {
	Player int
}

//Responses returns a slice of ints from [0,3].
func (m MinigamePipeMaze) Responses() []Response {
	return []Response{0, 1, 2, 3}
}

func (m MinigamePipeMaze) ControllingPlayer() int {
	return m.Player
}

//Handle gives player r 10 coins.
func (m MinigamePipeMaze) Handle(r Response, g *Game) {
	player := r.(int)
	g.AwardCoins(player, 10, true)
	g.EndGameTurn()
}

//MinigameBashnCash holds the implementation for Bash n Cash.
type MinigameBashnCash struct {
	BowsersBashnCash
}

//Handle calculates the number of coins taken from the solo player. The
//folowing events are set to distribute those coins among the other 3
//players.
func (m MinigameBashnCash) Handle(r Response, g *Game) {
	//TODO: code is copied from BowsersBashnCash.Handle()
	//Event interface needs to include *Game to remove this
	//effectively
	timesHit := r.(int)
	if timesHit == 0 {
		g.EndGameTurn()
		return
	}
	coinsLost := 0
	if m.Coins/5 < timesHit {
		coinsLost += m.Coins - (m.Coins % 5)
		timesHit -= m.Coins / 5
		coinsLost += timesHit
	} else {
		coinsLost += timesHit * 5
	}
	g.AwardCoins(m.Player, -coinsLost, true)
	nextEvent := MinigameBashnCashCoinAwards{0, m.Player, coinsLost}
	if m.Player == 0 {
		nextEvent.CurrentPlayer = 1
	}
	g.ExtraEvent = nextEvent
}

//MinigameBashnCashCoinAwards distributes a set of coins from a player to
//the other 3 players.
type MinigameBashnCashCoinAwards struct {
	CurrentPlayer int
	LosingPlayer  int
	Coins         int
}

//Responses returns a slice of ints from [0, m.Coins].
func (m MinigameBashnCashCoinAwards) Responses() []Response {
	return CPURangeEvent{0, m.Coins}.Responses()
}

func (m MinigameBashnCashCoinAwards) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives r coins to m.CurrentPlayer. The next event is set to the
//next player to receive coins, or ends the Game's turn if all players
//have received coins.
func (m MinigameBashnCashCoinAwards) Handle(r Response, g *Game) {
	coins := r.(int)
	g.AwardCoins(m.CurrentPlayer, coins, true)
	m.CurrentPlayer++
	if m.CurrentPlayer == m.LosingPlayer {
		m.CurrentPlayer++
	}
	if m.CurrentPlayer >= 4 {
		g.EndGameTurn()
	} else {
		g.ExtraEvent = m
	}
}

//MinigameBowlOver holds the implementation for Bowl Over.
type MinigameBowlOver struct {
	Player int
}

//MinigameBowlOverResponse is an ending result of the Bowl Over minigame.
type MinigameBowlOverResponse struct {
	Pins     int
	CharPins [3]bool
}

var MinigameBowlOverResponses = []Response{
	MinigameBowlOverResponse{0, [3]bool{false, false, false}},
	MinigameBowlOverResponse{1, [3]bool{false, false, false}},
	MinigameBowlOverResponse{2, [3]bool{false, false, false}},
	MinigameBowlOverResponse{0, [3]bool{false, false, true}},
	MinigameBowlOverResponse{1, [3]bool{false, false, true}},
	MinigameBowlOverResponse{2, [3]bool{false, false, true}},
	MinigameBowlOverResponse{0, [3]bool{false, true, false}},
	MinigameBowlOverResponse{1, [3]bool{false, true, false}},
	MinigameBowlOverResponse{2, [3]bool{false, true, false}},
	MinigameBowlOverResponse{0, [3]bool{false, true, true}},
	MinigameBowlOverResponse{1, [3]bool{false, true, true}},
	MinigameBowlOverResponse{2, [3]bool{false, true, true}},
	MinigameBowlOverResponse{0, [3]bool{true, false, false}},
	MinigameBowlOverResponse{1, [3]bool{true, false, false}},
	MinigameBowlOverResponse{2, [3]bool{true, false, false}},
	MinigameBowlOverResponse{0, [3]bool{true, false, true}},
	MinigameBowlOverResponse{1, [3]bool{true, false, true}},
	MinigameBowlOverResponse{2, [3]bool{true, false, true}},
	MinigameBowlOverResponse{0, [3]bool{true, true, false}},
	MinigameBowlOverResponse{1, [3]bool{true, true, false}},
	MinigameBowlOverResponse{2, [3]bool{true, true, false}},
	MinigameBowlOverResponse{0, [3]bool{true, true, true}},
	MinigameBowlOverResponse{1, [3]bool{true, true, true}},
	MinigameBowlOverResponse{2, [3]bool{true, true, true}},
}

//Responses returns all valid Bowl Over minigame endings.
func (m MinigameBowlOver) Responses() []Response {
	return MinigameBowlOverResponses
}

func (m MinigameBowlOver) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle calculates coin awards given a valid Response r.
func (m MinigameBowlOver) Handle(r Response, g *Game) {
	res := r.(MinigameBowlOverResponse)
	g.AwardCoins(m.Player, res.Pins, true)
	pinIndex := 0
	for i := 0; i < 4; i++ {
		if i == m.Player {
			continue
		}
		if res.CharPins[pinIndex] {
			g.GiveCoins(i, m.Player, 3, true)
		}
		pinIndex++
	}
	g.EndGameTurn()
}

//MinigameCraneGameCoins holds the implementation of the coin portion for
//Crane Game.
type MinigameCraneGameCoins struct {
	Player int
}

//Responses returns a slice of ints: {0, 1, 5, 10}.
func (m MinigameCraneGameCoins) Responses() []Response {
	return []Response{0, 1, 5, 10}
}

func (m MinigameCraneGameCoins) ControllingPlayer() int {
	return m.Player
}

//Handle gives the solo player r coins. If r == 0, the solo player decided
//to go after one of the 3 other players, and the next event is set to
//MinigameCraneGamePlayers.
func (m MinigameCraneGameCoins) Handle(r Response, g *Game) {
	coins := r.(int)
	if coins > 0 {
		g.AwardCoins(m.Player, coins, true)
		g.EndGameTurn()
	} else {
		groupPlayers := [3]int{}
		for i := 0; i < 3; i++ {
			if i >= m.Player {
				groupPlayers[i] = i + 1
			} else {
				groupPlayers[i] = i
			}
		}
		g.ExtraEvent = MinigameCraneGamePlayers{
			m.Player,
			groupPlayers,
		}
	}
}

//MinigameCraneGamePlayers holds the implementation of the player portion
//for Crane Game.
type MinigameCraneGamePlayers struct {
	SoloPlayer int
	Team       [3]int
}

//Responses returns a slice of ints containing the indexes of the 3 other
//players, and 4.
func (m MinigameCraneGamePlayers) Responses() []Response {
	return []Response{m.Team[0], m.Team[1], m.Team[2], 4}
}

func (m MinigameCraneGamePlayers) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives a third of Player r's coins to m.SoloPlayer. If r == 4, the
//solo player gains no coins.
func (m MinigameCraneGamePlayers) Handle(r Response, g *Game) {
	losingPlayer := r.(int)
	if losingPlayer != 4 {
		coins := g.Players[losingPlayer].Coins / 3
		g.GiveCoins(losingPlayer, m.SoloPlayer, coins, true)
	}
	g.EndGameTurn()
}

//MinigamePaddleBattle holds the implementation for Paddle Battle
type MinigamePaddleBattle struct {
	Player int
}

//Responses return a slice of int from [-10, 10].
func (m MinigamePaddleBattle) Responses() []Response {
	return CPURangeEvent{-10, 10}.Responses() //TODO: Find out max number hits possible
}

func (m MinigamePaddleBattle) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle calculates the amount of coins going from the solo player to/from
//the other 3 players. The solo player effectively loses r*3 coins.
func (m MinigamePaddleBattle) Handle(r Response, g *Game) {
	hits := r.(int)
	//TODO: Find out how coins are distributed when m.Player has 1 or 2 coins
	for i := 0; i < 4; i++ {
		if i == m.Player {
			i++
			continue
		}
		g.GiveCoins(m.Player, i, hits, true)
	}
	g.EndGameTurn()
}

//Minigame1V3Game is a enumeration of the available 1V3 minigames.
type Minigame1V3Game int

const (
	Minigame1V3PipeMaze Minigame1V3Game = iota
	Minigame1V3BashnCash
	Minigame1V3BowlOver
	Minigame1V3CoinBlockBash
	Minigame1V3TightropeTreachery
	Minigame1V3CraneGame
	Minigame1V3PiranhaPursuit
	Minigame1V3TugoWar
	Minigame1V3PaddleBattle
	Minigame1V3CoinShowerFlower
)

//Minigame1V3Selector selects which 1V3 minigame to play.
type Minigame1V3Selector struct {
	Player    int
	SoloCoins int
}

//Responses returns a slice of all Minigame1V3Game enumerations. If the
//solo player has 0 coins, then BashnCash is not selected.
func (m Minigame1V3Selector) Responses() []Response {
	if m.SoloCoins == 0 {
		return []Response{
			Minigame1V3PipeMaze,
			Minigame1V3BowlOver,
			Minigame1V3CoinBlockBash,
			Minigame1V3TightropeTreachery,
			Minigame1V3CraneGame,
			Minigame1V3PiranhaPursuit,
			Minigame1V3TugoWar,
			Minigame1V3PaddleBattle,
			Minigame1V3CoinShowerFlower,
		}
	}
	return []Response{
		Minigame1V3PipeMaze,
		Minigame1V3BashnCash,
		Minigame1V3BowlOver,
		Minigame1V3CoinBlockBash,
		Minigame1V3TightropeTreachery,
		Minigame1V3CraneGame,
		Minigame1V3PiranhaPursuit,
		Minigame1V3TugoWar,
		Minigame1V3PaddleBattle,
		Minigame1V3CoinShowerFlower,
	}
}

func (m Minigame1V3Selector) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the next event to the selected Minigame.
func (m Minigame1V3Selector) Handle(r Response, g *Game) {
	minigame := r.(Minigame1V3Game)
	switch minigame {
	case Minigame1V3PipeMaze:
		g.ExtraEvent = MinigamePipeMaze{m.Player}
	case Minigame1V3BashnCash:
		coins := g.Players[m.Player].Coins
		g.ExtraEvent = MinigameBashnCash{BowsersBashnCash{m.Player, coins}}
	case Minigame1V3BowlOver:
		g.ExtraEvent = MinigameBowlOver{m.Player}
	case Minigame1V3CoinBlockBash:
		g.ExtraEvent = CoinMinigameFFAReward{0, 30}
	case Minigame1V3TightropeTreachery:
		g.ExtraEvent = Minigame1V3Reward{m.Player}
	case Minigame1V3CraneGame:
		g.ExtraEvent = MinigameCraneGameCoins{m.Player}
	case Minigame1V3PiranhaPursuit:
		g.ExtraEvent = Minigame1V3Reward{m.Player}
	case Minigame1V3TugoWar:
		g.ExtraEvent = Minigame1V3Reward{m.Player}
	case Minigame1V3PaddleBattle:
		g.ExtraEvent = MinigamePaddleBattle{m.Player}
	case Minigame1V3CoinShowerFlower:
		g.ExtraEvent = Throwable1V3Minigame{m.Player, CoinMinigameFFAReward{0, 30}}
	}
}

//Minigame1PRewards handles 1P minigame rewards. The player will either
//gain 10 coins, or lose 5 coins.
type Minigame1PRewards struct {
	Player int
}

//Responses returns a slice of ints, {-5, 10}.
func (m Minigame1PRewards) Responses() []Response {
	return []Response{-5, 10}
}

func (m Minigame1PRewards) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle gives m.Player r coins.
func (m Minigame1PRewards) Handle(r Response, g *Game) {
	coins := r.(int)
	g.AwardCoins(m.Player, coins, true)
	g.EndCharacterTurn()
}

//MinigameMemoryMatch holds the implementation for Memory Match.
type MinigameMemoryMatch struct {
	Minigame1PRewards
}

//Responses returns a slice of ints, {0, 2, 4, 6, 10}.
func (m MinigameMemoryMatch) Responses() []Response {
	return []Response{0, 2, 4, 6, 10}
}

//MinigameSlotMachine hodls the implementation for Slot Machine.
type MinigameSlotMachine struct {
	Minigame1PRewards
}

//Responses returns a slice of ints, {0, 1, 3, 5, 6, 8, 10, 20}.
func (m MinigameSlotMachine) Responses() []Response {
	return []Response{0, 1, 3, 5, 6, 8, 10, 20}
}

//MinigameWhackaPlant holds the implementation for Whack a Plant.
type MinigameWhackaPlant struct {
	Minigame1PRewards
}

//Responses returns a slice of ints from [0, 36].
func (m MinigameWhackaPlant) Responses() []Response {
	return CPURangeEvent{0, 36}.Responses()
}

//MinigameTeeteringTowers holds the implementation for Teetering Towers.
type MinigameTeeteringTowers struct {
	Minigame1PRewards
}

//Responses returns a slice of ints, {-5, 10, 11, 15, 16}.
func (m MinigameTeeteringTowers) Responses() []Response {
	return []Response{-5, 10, 11, 15, 16} //Mix of coin and coinbag
}

//Minigame1PGame is a enumeration of the available 1P minigames.
type Minigame1PGame int

const (
	Minigame1PMemoryMatch Minigame1PGame = iota
	Minigame1PSlotMachine
	Minigame1PShellGame
	Minigame1PGhostGuess
	Minigame1PPedalPower
	Minigame1PWhackaPlant
	Minigame1PGroundPound
	Minigame1PTeeteringTowers
	Minigame1PKnockBlockTower
	Minigame1PLimboDance
)

//Minigame1PSelector selects which 1P minigame to play.
type Minigame1PSelector struct {
	Player int
}

//Responses returns a slice of all Minigame1PGame enumerations.
func (m Minigame1PSelector) Responses() []Response {
	return []Response{
		Minigame1PMemoryMatch,
		Minigame1PSlotMachine,
		Minigame1PShellGame,
		Minigame1PGhostGuess,
		Minigame1PPedalPower,
		Minigame1PWhackaPlant,
		Minigame1PGroundPound,
		Minigame1PTeeteringTowers,
		Minigame1PKnockBlockTower,
		Minigame1PLimboDance,
	}
}

func (m Minigame1PSelector) ControllingPlayer() int {
	return CPU_PLAYER
}

//Handle sets the next event to the selected Minigame.
func (m Minigame1PSelector) Handle(r Response, g *Game) {
	game := r.(Minigame1PGame)
	baseGame := Minigame1PRewards{m.Player}
	switch game {
	case Minigame1PMemoryMatch:
		g.ExtraEvent = MinigameMemoryMatch{baseGame}
	case Minigame1PSlotMachine:
		g.ExtraEvent = MinigameSlotMachine{baseGame}
	case Minigame1PShellGame:
		g.ExtraEvent = baseGame
	case Minigame1PGhostGuess:
		g.ExtraEvent = baseGame
	case Minigame1PPedalPower:
		g.ExtraEvent = baseGame
	case Minigame1PWhackaPlant:
		g.ExtraEvent = MinigameWhackaPlant{baseGame}
	case Minigame1PGroundPound:
		g.ExtraEvent = baseGame
	case Minigame1PTeeteringTowers:
		g.ExtraEvent = MinigameTeeteringTowers{baseGame}
	case Minigame1PKnockBlockTower:
		g.ExtraEvent = baseGame
	case Minigame1PLimboDance:
		g.ExtraEvent = baseGame
	}
}

//MinigameTeam is an enumeration of the available teams players can be on.
type MinigameTeam int

const (
	BlueTeam MinigameTeam = iota
	RedTeam
	GreenTeam
)

//SpaceToTeam is a mapping from SpaceType to MinigameTeam.
func SpaceToTeam(s SpaceType) MinigameTeam {
	switch s {
	case Blue, Mushroom, MinigameSpace, Chance:
		return BlueTeam
	case Red, Bowser:
		return RedTeam
	default:
		return GreenTeam
	}
}

//GetMinigame sets up the next minigame type (FFA/2v2/1v3). If any players
//are on the *green* team, the next event is set to determine that
//player's team selection.
func (g *Game) GetMinigame() {
	var blueTeam []int
	var redTeam []int
	for i, p := range g.Players {
		if SpaceToTeam(p.LastSpaceType) == BlueTeam {
			blueTeam = append(blueTeam, i)
		} else if SpaceToTeam(p.LastSpaceType) == RedTeam {
			redTeam = append(redTeam, i)
		}
	}

	var minigame Event
	switch len(blueTeam) {
	case 0, 4:
		minigame = MinigameFFASelector{}
	case 1:
		minigame = Minigame1V3Selector{blueTeam[0], g.Players[blueTeam[0]].Coins}
	case 2:
		minigame = Minigame2V2Selector{
			[2]int{blueTeam[0], blueTeam[1]},
			[2]int{redTeam[0], redTeam[1]},
		}
	case 3:
		minigame = Minigame1V3Selector{redTeam[0], g.Players[redTeam[0]].Coins}
	}
	g.ExtraEvent = minigame
}

//FindGreenPlayer looks through the 4 players to find one that is on the
//*green* team.
func (g *Game) FindGreenPlayer() {
	for i, p := range g.Players {
		if SpaceToTeam(p.LastSpaceType) == GreenTeam {
			g.ExtraEvent = DeterminePlayerTeamEvent{
				Player: i,
			}
			return
		}
	}
	g.GetMinigame()
}
