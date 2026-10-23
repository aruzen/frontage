package pve

import (
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	action_register "frontage/pkg/engine/impl/action/register"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/model"
	"frontage/pkg/network/game_dispatcher"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
	"time"
)

type GameInfo struct {
	DrawPhaseLimit     time.Duration
	TurnTimeLimit      time.Duration
	ExtraTimeLimit     time.Duration
	SyncResendDuration time.Duration
	SyncResendMax      int
}

func DefaultGameInfo() GameInfo {
	return GameInfo{
		DrawPhaseLimit:     15 * time.Second,
		TurnTimeLimit:      time.Minute,
		SyncResendDuration: time.Second,
		SyncResendMax:      3,
	}
}

type RequireContents struct {
	ActionRepo   *repository.ActionRepository
	CardRepo     *repository.CardRepository
	ActEventDisp *game_dispatcher.ActEventDispatcher
	GameInitDisp *game_dispatcher.GameInitializeDispatcher
}

// ============================

func InitRequireRepositories(requireRepos *RequireContents) {
	action_register.Init()
	if requireRepos.ActionRepo == nil {
		requireRepos.ActionRepo = repository.NewActionRepository(action.FindActionModify, action.FindActionEffect)
	}
	if requireRepos.CardRepo == nil {
		requireRepos.CardRepo = repository.NewCardRepository()
	}
	if requireRepos.ActEventDisp == nil {
		requireRepos.ActEventDisp = game_dispatcher.NewActEventDispatcher(nil, nil)
	}
	if requireRepos.GameInitDisp == nil {
		requireRepos.GameInitDisp = game_dispatcher.NewGameInitializeDispatcher()
	}

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
}

func DemoDecks(requireRepos *RequireContents) (model.Cards, model.Cards, model.Cards, model.Cards) {
	mainNaturoDeck := make(model.Cards, 0, 8)
	subNaturoDeck := make(model.Cards, 0, 8)
	for k, v := range map[pkg.LocalizeTag]int{
		"少年グルーシャ":    4,
		"爛漫に咲く花・姫百子": 4,
	} {
		find, err := requireRepos.CardRepo.Find(k)
		if err != nil {
			continue
		}
		for i := 0; i < v; i++ {
			mainNaturoDeck = append(mainNaturoDeck, find.CardMirror(uuid.New()))
			subNaturoDeck = append(subNaturoDeck, find.CardMirror(uuid.New()))
		}
	}

	mainPyroDeck := make(model.Cards, 0, 8)
	subPyroDeck := make(model.Cards, 0, 8)
	for k, v := range map[pkg.LocalizeTag]int{
		"生まれたばかりの灯・ベビードレイク": 4,
		"旗将炎猿・ドモルドス":        4,
	} {
		find, err := requireRepos.CardRepo.Find(k)
		if err != nil {
			continue
		}
		for i := 0; i < v; i++ {
			mainPyroDeck = append(mainPyroDeck, find.CardMirror(uuid.New()))
			subPyroDeck = append(subPyroDeck, find.CardMirror(uuid.New()))
		}
	}
	return mainNaturoDeck, subNaturoDeck, mainPyroDeck, subPyroDeck
}
