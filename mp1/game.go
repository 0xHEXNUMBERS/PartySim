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
	KoopaPasses   int
	ExtraEvent    Event
	Config        GameConfig
}

func (g *Game) SetDiceBlock() {
	if g.Config.RedDice || g.Config.BlueDice || g.Config.WarpDice || g.Config.EventsDice {
		g.ExtraEvent = PickDiceBlock{g.CurrentPlayer, g.Config}
	} else {
		g.ExtraEvent = NormalDiceBlock{g.CurrentPlayer}
	}
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
		g.SetDiceBlock()
	} else {
		g.ExtraEvent = StarLocationEvent{g.StarSpaces, 0, 0}
	}
	return g
}

func (g *Game) LastFiveTurns() bool {
	return g.Config.MaxTurns-g.Turn <= 5
}

func (g *Game) AwardCoins(player, coins int, minigame bool) int {
	coins0 := g.Players[player].Coins
	g.Players[player].Coins += coins
	g.Players[player].Coins = max(g.Players[player].Coins, 0)
	coinsGiven := g.Players[player].Coins - coins0
	if minigame {
		g.Players[player].MinigameCoins += coins
	}
	g.Players[player].MaxCoins = max(
		g.Players[player].MaxCoins,
		g.Players[player].Coins,
	)
	return coinsGiven
}

func (g *Game) GiveCoins(givingPlayer, takingPlayer, coins int, minigame bool) {
	coinsTaken := -g.AwardCoins(givingPlayer, -coins, minigame)
	g.AwardCoins(takingPlayer, coinsTaken, minigame)
}

func (g *Game) CheckLinks(player, chain, moves int) (branch bool) {
	//Check for links
	if g.Board.Links != nil {
		chainLinks := *(g.Board.Links)
		linksPtr, ok := chainLinks[chain]
		if !ok {
			g.Players[player].CurrentSpace.Space = 0
			return
		}
		links := *linksPtr
		switch len(links) {
		case 1:
			g.Players[player].CurrentSpace = links[0]
		default:
			g.ExtraEvent = BranchEvent{
				player,
				moves - 1,
				linksPtr,
			}
			return true
		}
	}
	return false
}

func (g *Game) ActivateSpace(player int) {
	//Activate Space
	chains := *g.Board.Chains
	playerPos := g.Players[player].CurrentSpace
	curSpace := chains[playerPos.Chain][playerPos.Space]
	//Star space requires preprocessing
	if curSpace.Type == Star {
		starIndex := g.StarSpaces.GetIndex(playerPos)
		if g.StarSpaces.AbsoluteVisited&(1<<starIndex) > 0 {
			g.Players[player].LastSpaceType = Chance
		} else {
			g.Players[player].LastSpaceType = Blue
		}
	}
	//Perform space action
	switch g.Players[player].LastSpaceType {
	case Invisible:
		//Stopping Event should set LastSpaceType
		curSpace.StoppingEvent(g, player)
		g.ActivateSpace(player)
	case Blue:
		if g.LastFiveTurns() {
			g.AwardCoins(player, 6, false)
		} else {
			g.AwardCoins(player, 3, false)
		}
		g.EndCharacterTurn()
	case Red:
		if g.LastFiveTurns() {
			g.AwardCoins(player, -6, false)
		} else {
			g.AwardCoins(player, -3, false)
		}
		g.EndCharacterTurn()
	case Mushroom:
		g.ExtraEvent = MushroomEvent{player}
	case Happening:
		g.Players[player].HappeningCount++
		g.ExtraEvent = nil
		curSpace.StoppingEvent(g, player)
		if g.ExtraEvent == nil {
			g.EndCharacterTurn()
		}
	case Bowser:
		g.PreBowserCheck(player)
	case MinigameSpace:
		g.ExtraEvent = Minigame1PSelector{player}
	case Chance:
		g.ExtraEvent = ChanceTime{Player: player}
	}
}

func (g *Game) MovePlayer(playerIdx, moves int) {
	chains := *g.Board.Chains
	playerPos := g.Players[playerIdx].CurrentSpace
	for moves > 0 {
		playerPos.Space++
		if playerPos.Space >= len(chains[playerPos.Chain]) {
			if g.CheckLinks(playerIdx, playerPos.Chain, moves) {
				return
			}
			playerPos = g.Players[playerIdx].CurrentSpace
		}
		g.Players[playerIdx].CurrentSpace = playerPos
		curSpace := chains[playerPos.Chain][playerPos.Space]
		switch curSpace.Type {
		case Invisible:
			if curSpace.PassingEvent != nil {
				g.ExtraEvent = nil
				moves = curSpace.PassingEvent(g, playerIdx, moves)
				if g.ExtraEvent != nil {
					return
				}
				//The PassingEvent() sets the new player position.
				//As such, we must update our understanding of where
				//the player is at
				playerPos = g.Players[playerIdx].CurrentSpace
				curSpace = chains[playerPos.Chain][playerPos.Space]
			} else {
				moves--
			}
		case Start:
			if !g.Config.NoKoopa {
				g.KoopaPasses++
				if g.KoopaPasses%10 == 0 {
					g.AwardCoins(playerIdx, 20, false)
				} else {
					g.AwardCoins(playerIdx, 10, false)
				}
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
		case BogusItem:
			g.AwardCoins(playerIdx, -g.Board.BowserCoins, false)
		default:
			moves--
		}
	}
	curSpace := chains[playerPos.Chain][playerPos.Space]
	g.Players[playerIdx].LastSpaceType = curSpace.Type
	if g.Config.EventsDice &&
		(curSpace.Type == Blue || curSpace.HiddenBlock) {
		g.ExtraEvent = HiddenBlockEvent{playerIdx}
		return
	}
	g.ActivateSpace(playerIdx)
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
		g.SetDiceBlock()
	}
}

func (g *Game) StartMinigamePrep() {
	g.FindGreenPlayer()
}

func (g *Game) EndCharacterTurn() {
	if g.Board.EndCharacterTurnEvent != nil {
		g.Board.EndCharacterTurnEvent(g, g.CurrentPlayer)
	}
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
	g.SetDiceBlock()
}
