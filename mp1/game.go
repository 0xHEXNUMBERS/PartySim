package mp1

type GameConfig struct {
	MaxTurns     uint8
	NoBonusStars bool
	NoKoopa      bool
	NoBoo        bool
	RedDice      bool
	BlueDice     bool
	WarpDice     bool
	EventsDice   bool
}

type Game struct {
	Board
	StarSpaces    StarData
	Players       [4]Player
	Turn          uint8
	CurrentPlayer int
	ExtraEvent    Event
	Config        GameConfig
}

func InitializeGame(b Board, config GameConfig) *Game {
	g := &Game{
		Board: b,
	}
	var chainSpaces []ChainSpace
	var count uint8
	var startSpace ChainSpace
	for ci, chain := range *b.Chains {
		for si, space := range chain {
			if space.Type == Star {
				chainSpaces = append(
					chainSpaces,
					ChainSpace{ci, si},
				)
				count++
			} else if space.Type == Start {
				startSpace = ChainSpace{ci, si}
			}
		}
	}
	if count > 0 {
		g.StarSpaces.IndexToPosition = &chainSpaces
		g.StarSpaces.StarSpaceCount = count
	}
	g.Config = config

	for i := 0; i < len(g.Players); i++ {
		g.Players[i].CurrentSpace = startSpace
		g.Players[i].LastSpaceType = Start
		g.Players[i].Coins = 10
	}

	if g.StarSpaces.StarSpaceCount == 1 { //Go ahead and set star pos
		g.StarSpaces.CurrentStarSpace = (*g.StarSpaces.IndexToPosition)[0]
	}
	if g.StarSpaces.StarSpaceCount <= 1 {
		g.ExtraEvent = PickDiceBlock{0, g.Config}
	} else {
		g.ExtraEvent = StarLocationEvent{g.StarSpaces, 0, 0}
	}
	return g
}

func (g *Game) LastFiveTurns() bool {
	return g.Config.MaxTurns-g.Turn <= 5
}

func (g *Game) AwardCoins(player, coins int, minigame bool) {
	g.Players[player].Coins += coins
	if g.Players[player].Coins < 0 {
		g.Players[player].Coins = 0
	}
	if minigame {
		g.Players[player].MinigameCoins += coins
		if g.Players[player].MinigameCoins < 0 {
			g.Players[player].MinigameCoins = 0
		}
	}
	if g.Players[player].Coins > g.Players[player].MaxCoins {
		g.Players[player].MaxCoins = g.Players[player].Coins
	}
}

func (g *Game) MovePlayer(playerIdx, moves int) {
	chains := *g.Board.Chains
	playerPos := g.Players[playerIdx].CurrentSpace
	for moves > 0 {
		playerPos.Space++
		g.Players[playerIdx].CurrentSpace = playerPos
		if playerPos.Space >= len(chains[playerPos.Chain]) {
			playerPos.Space = 0
		}
		curSpace := chains[playerPos.Chain][playerPos.Space]
		switch curSpace.Type {
		case Invisible:
			g.ExtraEvent = nil
			curSpace.PassingEvent(g, playerIdx, moves)
			if g.ExtraEvent != nil {
				g.Players[playerIdx].CurrentSpace = playerPos
				return
			}
		case Start:
			if !g.Config.NoKoopa {
				g.AwardCoins(playerIdx, 10, false)
			}
		case Star:
			if playerPos == g.StarSpaces.CurrentStarSpace &&
				g.Players[playerIdx].Coins >= 20 {
				g.Players[playerIdx].Stars++
				g.AwardCoins(playerIdx, -20, false)
				if g.StarSpaces.StarSpaceCount > 1 {
					g.ExtraEvent = StarLocationEvent{
						g.StarSpaces,
						playerIdx,
						moves,
					}
					return
				}
			}
			moves--
		case Boo:
			if !g.Config.NoBoo {
				booEvt := BooEvent{
					playerIdx,
					g.Players,
					moves,
					g.Players[playerIdx].Coins,
				}
				if len(booEvt.Responses()) != 0 {
					g.ExtraEvent = booEvt
					return
				}
			}
		default:
			moves--
		}
	}
	//Stop on Space
	g.Players[playerIdx].CurrentSpace = playerPos
	//Activate Space
	curSpace := chains[playerPos.Chain][playerPos.Space]
	g.Players[playerIdx].LastSpaceType = curSpace.Type
	//Star space requires preprocessing
	if curSpace.Type == Star {
		starIndex := g.StarSpaces.GetIndex(playerPos)
		if g.StarSpaces.AbsoluteVisited&(1<<starIndex) > 0 {
			g.Players[playerIdx].LastSpaceType = Chance
		} else {
			g.Players[playerIdx].LastSpaceType = Blue
		}
	}
	//Perform space action
	switch g.Players[playerIdx].LastSpaceType {
	case Blue:
		if g.LastFiveTurns() {
			g.AwardCoins(playerIdx, 6, false)
		} else {
			g.AwardCoins(playerIdx, 3, false)
		}
		g.EndCharacterTurn()
	case Red:
		if g.LastFiveTurns() {
			g.AwardCoins(playerIdx, -6, false)
		} else {
			g.AwardCoins(playerIdx, -3, false)
		}
		g.EndCharacterTurn()
	case Mushroom:
		g.ExtraEvent = MushroomEvent{playerIdx}
	case Happening:
		g.Players[playerIdx].HappeningCount++
		g.ExtraEvent = nil
		curSpace.StoppingEvent(g)
		if g.ExtraEvent == nil {
			g.EndCharacterTurn()
		}
	case Bowser:
		g.PreBowserCheck(playerIdx)
	case MinigameSpace:
		g.ExtraEvent = MinigameEvent{[4]int{playerIdx, 0, 0, 0}, Minigame1P}
	case Chance:
		g.ExtraEvent = ChanceTime{Player: playerIdx}
	}
}

func (g *Game) AwardBonusStars() {
	if g.Config.NoBonusStars {
		return
	}

	maxCoins := g.Players[0].MaxCoins
	for i := 1; i < 4; i++ {
		maxCoins = max(maxCoins, g.Players[i].MaxCoins)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].MaxCoins == maxCoins {
			g.Players[i].Stars++
		}
	}

	maxMinigameCoins := g.Players[0].MinigameCoins
	for i := 1; i < 4; i++ {
		maxMinigameCoins = max(maxMinigameCoins, g.Players[i].MinigameCoins)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].MinigameCoins == maxMinigameCoins {
			g.Players[i].Stars++
		}
	}

	maxHappening := g.Players[0].HappeningCount
	for i := 1; i < 4; i++ {
		maxHappening = max(maxHappening, g.Players[i].HappeningCount)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].HappeningCount == maxHappening {
			g.Players[i].Stars++
		}
	}
}

func (g *Game) Winners() []int {
	maxStarHolders := []int{}
	maxStars := g.Players[0].Stars
	for i := 1; i < 4; i++ {
		maxStars = max(maxStars, g.Players[i].Stars)
	}
	for i := 0; i < 4; i++ {
		if g.Players[i].Stars == maxStars {
			maxStarHolders = append(maxStarHolders, i)
		}
	}

	winners := []int{}
	maxCoins := g.Players[maxStarHolders[0]].Coins
	for i := 1; i < len(maxStarHolders); i++ {
		maxCoins = max(maxCoins, g.Players[maxStarHolders[i]].Coins)
	}
	for i := 0; i < len(maxStarHolders); i++ {
		if g.Players[maxStarHolders[i]].Coins == maxCoins {
			winners = append(winners, maxStarHolders[i])
		}
	}
	return winners
}

func (g *Game) EndGameTurn() {
	g.Turn++
	if g.Turn == g.Config.MaxTurns {
		g.AwardBonusStars()
		//Game is over, no more events
		g.ExtraEvent = nil
	} else {
		if g.Players[g.CurrentPlayer].SkipTurn {
			g.Players[g.CurrentPlayer].SkipTurn = false
			g.EndCharacterTurn()
			return
		}
		g.ExtraEvent = PickDiceBlock{g.CurrentPlayer, g.Config}
	}
}

func (g *Game) StartMinigamePrep() {
	g.FindGreenPlayer()
}

func (g *Game) EndCharacterTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	if g.CurrentPlayer == 0 {
		g.StartMinigamePrep()
		return
	}
	if g.Players[g.CurrentPlayer].SkipTurn {
		g.Players[g.CurrentPlayer].SkipTurn = false
		g.EndCharacterTurn()
		return
	}
	g.ExtraEvent = PickDiceBlock{g.CurrentPlayer, g.Config}
}
