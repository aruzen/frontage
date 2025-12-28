package pvp

import (
	"context"
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
	"time"
)

type GameInfo struct {
	TurnTimeLimit time.Duration
}

func DefaultGameInfo() GameInfo {
	return GameInfo{
		TurnTimeLimit: time.Minute,
	}
}

type RequireRepositories struct {
	CardRepo *repository.CardRepository
}

func Game(requireRepos RequireRepositories, ids [2]uuid.UUID, info GameInfo) {
	_ = repository.NewActionRepository(func(tag logic.ModifyActionTag) logic.ModifyAction {
		return nil
	}, func(tag logic.EffectActionTag) logic.EffectAction {
		return action.FindActionEffect(tag)
	})

	requireRepos.CardRepo.Insert(card.NewBasePiece(model.NewBaseCard("少年グルーシャ", model.Materials{
		model.NATURO: 1,
	}),
		2, 2, 1,
		[]pkg.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}},
		[]pkg.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}))

	requireRepos.CardRepo.Insert(card.NewBasePiece(model.NewBaseCard("爛漫に咲く花・姫百子", model.Materials{
		model.NATURO: 4,
	}),
		6, 4, 2,
		[]pkg.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {-1, -1}, {1, -1}},
		[]pkg.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {-1, -1}, {1, -1}}))

	requireRepos.CardRepo.Insert(card.NewBasePiece(model.NewBaseCard("生まれたばかりの灯・ベビードレイク", model.Materials{
		model.PYRO: 1,
	}),
		1, 0, 2,
		[]pkg.Point{{1, 0}, {-1, 0}, {-2, 0}},
		[]pkg.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}))

	requireRepos.CardRepo.Insert(card.NewBasePiece(model.NewBaseCard("旗将炎猿・ドモルドス", model.Materials{
		model.PYRO: 4,
	}),
		6, 4, 2,
		[]pkg.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {-1, -1}, {1, -1}},
		[]pkg.Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {-1, -1}, {1, -1}}))

	mainNaturoDeck := make(model.Cards, 8)
	subNaturoDeck := make(model.Cards, 8)
	for k, v := range map[pkg.LocalizeTag]int{
		"少年グルーシャ":    4,
		"爛漫に咲く花・姫百子": 4,
	} {
		find, err := requireRepos.CardRepo.Find(k)
		if err != nil {
			continue
		}
		for i := 0; i < v; i++ {
			mainNaturoDeck = append(mainNaturoDeck, find.CardCopy())
			subNaturoDeck = append(subNaturoDeck, find.CardCopy())
		}
	}

	mainPyroDeck := make(model.Cards, 8)
	subPyroDeck := make(model.Cards, 8)
	for k, v := range map[pkg.LocalizeTag]int{
		"生まれたばかりの灯・ベビードレイク": 4,
		"旗将炎猿・ドモルドス":        4,
	} {
		find, err := requireRepos.CardRepo.Find(k)
		if err != nil {
			continue
		}
		for i := 0; i < v; i++ {
			mainPyroDeck = append(mainPyroDeck, find.CardCopy())
			subPyroDeck = append(subPyroDeck, find.CardCopy())
		}
	}

	players := [2]model.Player{
		model.NewLocalPlayer(ids[0], mainNaturoDeck, subNaturoDeck),
		model.NewLocalPlayer(ids[1], mainPyroDeck, subPyroDeck),
	}

	board := model.NewBoard(model.NewBoardInfo(pkg.Size{7, 7}, model.GENERATION_STRATEGY_SWAP), players)

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
