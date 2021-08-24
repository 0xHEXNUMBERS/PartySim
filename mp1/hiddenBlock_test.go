package mp1

import "testing"

func TestHiddenBlock(t *testing.T) {
	g := *InitializeGame(YTI, GameConfig{MaxTurns: 20, EventsDice: true})
	g.Players[0].CurrentSpace = ChainSpace{0, 21}
	g.MovePlayer(0, 1) //Blue Space
	EventIs(HiddenBlockEvent{0}, g.NextEvent, "", t)

	gHidden := g
	gHidden.NextEvent.Handle(true, &gHidden)
	EventIs(EventDiceBlock{0}, gHidden.NextEvent, "Dice", t)

	gBlue := g
	gBlue.NextEvent.Handle(false, &gBlue)
	EventIs(PickDiceBlock{1, gBlue.Config}, gBlue.NextEvent, "Player", t)
	CoinsIs(13, 0, gBlue, "", t)
}
