package pvp

import (
	"context"
	"frontage/pkg"
	"frontage/pkg/data"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"frontage/pkg/handler"
	"math/rand"
	"time"
)

type GameInfo struct {
	TurnTimeLimit time.Duration
}

func Game(info GameInfo, playersData [2]data.Player) {
	first := rand.Int() % 2
	second := (first + 1) % 2
	board := model.NewBoard(model.NewBoardInfo(pkg.Size{7, 7}, model.GENERATION_STRATEGY_SWAP), [2]*model.Player{
		model.NewPlayer(handler.InstantiationDeck(playersData[first].MainDeck), handler.InstantiationDeck(playersData[first].SubDeck)),
		model.NewPlayer(handler.InstantiationDeck(playersData[second].MainDeck), handler.InstantiationDeck(playersData[second].SubDeck)),
	})

	es := &logic.EventSystem{
		Board: board,
	}
	es.Emit(logic.NewEffectEvent(logic.GAME_START_ACTION, logic.GAME_START_ACTION.MakeState()))

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	for {

	}

	cancel()
}

func listenPlayerInput() {

}
