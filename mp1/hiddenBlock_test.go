package mp1

import "testing"

var BlueBoard = MakeSimpleBoard(Blue)

func TestHiddenBlock(t *testing.T) {
	g := *InitializeGame(BlueBoard, GameConfig{MaxTurns: 20, EventsDice: true})

	//Engine always sets first diceblock for all players to normal dice
	//block. We must tell the engine we're not on the first turn to test
	//hidden block behavior.
	g.Turn = 1

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
