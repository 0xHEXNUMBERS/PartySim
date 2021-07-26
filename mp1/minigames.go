package mp1

type MinigameAwards [4]int

var MinigameRewardsFFA = []Response{
	MinigameAwards{10, 10, 10, 10},
	MinigameAwards{10, 10, 10, 0},
	MinigameAwards{10, 10, 10, 0},
	MinigameAwards{10, 10, 0, 10},
	MinigameAwards{10, 10, 0, 0},
	MinigameAwards{10, 0, 10, 10},
	MinigameAwards{10, 0, 10, 0},
	MinigameAwards{10, 0, 0, 10},
	MinigameAwards{10, 0, 0, 0},
	MinigameAwards{0, 10, 10, 10},
	MinigameAwards{0, 10, 10, 0},
	MinigameAwards{0, 10, 0, 10},
	MinigameAwards{0, 10, 0, 0},
	MinigameAwards{0, 0, 10, 10},
	MinigameAwards{0, 0, 10, 0},
	MinigameAwards{0, 0, 0, 10},
	MinigameAwards{0, 0, 0, 0},
}

var MinigameRewards2V2 = []Response{
	MinigameAwards{10, 10, 0, 0},
	MinigameAwards{0, 0, 10, 10},
	MinigameAwards{0, 0, 0, 0},
}

var MinigameRewards1V3 = []Response{
	MinigameAwards{15, -5, -5, -5},
	MinigameAwards{-15, 5, 5, 5},
	MinigameAwards{0, 0, 0, 0},
}

var MinigameRewards1P = []Response{
	MinigameAwards{-5, 0, 0, 0},
	MinigameAwards{0, 0, 0, 0},
	MinigameAwards{1, 0, 0, 0},
	MinigameAwards{2, 0, 0, 0},
	MinigameAwards{3, 0, 0, 0},
	MinigameAwards{4, 0, 0, 0},
	MinigameAwards{5, 0, 0, 0},
	MinigameAwards{6, 0, 0, 0},
	MinigameAwards{7, 0, 0, 0},
	MinigameAwards{8, 0, 0, 0},
	MinigameAwards{9, 0, 0, 0},
	MinigameAwards{10, 0, 0, 0},
	MinigameAwards{11, 0, 0, 0},
	MinigameAwards{12, 0, 0, 0},
	MinigameAwards{13, 0, 0, 0},
	MinigameAwards{14, 0, 0, 0},
	MinigameAwards{15, 0, 0, 0},
	MinigameAwards{16, 0, 0, 0},
	MinigameAwards{17, 0, 0, 0},
	MinigameAwards{18, 0, 0, 0},
	MinigameAwards{19, 0, 0, 0},
	MinigameAwards{20, 0, 0, 0},
	MinigameAwards{21, 0, 0, 0},
	MinigameAwards{22, 0, 0, 0},
	MinigameAwards{23, 0, 0, 0},
	MinigameAwards{24, 0, 0, 0},
	MinigameAwards{25, 0, 0, 0},
	MinigameAwards{26, 0, 0, 0},
	MinigameAwards{27, 0, 0, 0},
	MinigameAwards{28, 0, 0, 0},
	MinigameAwards{29, 0, 0, 0},
	MinigameAwards{30, 0, 0, 0},
	MinigameAwards{31, 0, 0, 0},
	MinigameAwards{32, 0, 0, 0},
	MinigameAwards{33, 0, 0, 0},
	MinigameAwards{34, 0, 0, 0},
	MinigameAwards{35, 0, 0, 0},
	MinigameAwards{36, 0, 0, 0},
}

type MinigameType int

const (
	MinigameFFA MinigameType = iota
	Minigame2V2
	Minigame1V3
	Minigame1P
)

type MinigameEvent struct {
	//Player IDs
	//FFA -> [Player0, Player1, Player2, Player3]
	//1V3 -> [Team1, Team2, Team2, Team2]
	//2V2 -> [Team1, Team1, Team2, Team2]
	//1P  -> [Team1, nil, nil, nil]
	PlayerIDs [4]int
	Type      MinigameType
}

func (m MinigameEvent) Responses() []Response {
	switch m.Type {
	case MinigameFFA:
		return MinigameRewardsFFA
	case Minigame2V2:
		return MinigameRewards2V2
	case Minigame1V3:
		return MinigameRewards1V3
	case Minigame1P:
		return MinigameRewards1P
	}
	//Unreachable
	return nil
}

func (m MinigameEvent) ControllingPlayer() int {
	return CPU_PLAYER
}

func (m MinigameEvent) Handle(r Response, g *Game) {
	awards := r.(MinigameAwards)
	for i, player := range m.PlayerIDs {
		g.AwardCoins(player, awards[i], true)
	}
	g.EndGameTurn()
}

type MinigameTeam int

const (
	BlueTeam MinigameTeam = iota
	RedTeam
	GreenTeam
)

func SpaceToTeam(s SpaceType) MinigameTeam {
	switch s {
	case Blue, Mushroom:
		return BlueTeam
	case Red, Bowser:
		return RedTeam
	default:
		return GreenTeam
	}
}

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

	var minigameType MinigameType
	switch len(blueTeam) {
	case 0, 4:
		minigameType = MinigameFFA
	case 1, 3:
		minigameType = Minigame1V3
	case 2:
		minigameType = Minigame2V2
	}

	minigame := MinigameEvent{Type: minigameType}
	var players []int
	if len(redTeam) == 1 { //Put 1 person team in front
		players = append(redTeam, blueTeam...)
	} else {
		players = append(blueTeam, redTeam...)
	}
	for i := range blueTeam {
		minigame.PlayerIDs[i] = players[i]
	}
	g.ExtraEvent = minigame
}

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
