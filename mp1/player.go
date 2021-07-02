package mp1

type Player struct {
	Char           string
	Stars          int
	Coins          int
	CurrentSpace   ChainSpace
	SkipTurn       bool
	MaxCoins       int
	HappeningCount int
	MinigameCoins  int
}

//CPU_PLAYER acts as a separate player (player 5) to control
//events that normal players have no control over
const CPU_PLAYER int = 4
