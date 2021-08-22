package mp1

//CPU_PLAYER acts as a separate player (player 5) to control
//events that normal players have no control over
const CPU_PLAYER int = 4

//Player holds all player data, including bonus star stats.
type Player struct {
	Char          string
	Stars         int
	Coins         int
	CurrentSpace  ChainSpace
	SkipTurn      bool
	LastSpaceType SpaceType

	//Bonus Star Data
	MaxCoins       int
	HappeningCount int
	MinigameCoins  int
}

//NewPlayer generates a new player with a given name.
func NewPlayer(name string, stars, coins int, space ChainSpace) Player {
	return Player{
		name,
		stars,
		coins,
		space,
		false,
		Start,
		0,
		0,
		0,
	}
}
