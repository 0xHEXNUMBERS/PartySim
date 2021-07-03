package mp1

//CPU_PLAYER acts as a separate player (player 5) to control
//events that normal players have no control over
const CPU_PLAYER int = 4

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
