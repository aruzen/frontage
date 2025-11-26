package logic

import (
	"frontage/pkg/engine/model"
)

type SystemNoticeState struct {
}

type SystemNoticeContext struct {
}

type systemNoticeAction struct {
	BaseAction[SystemNoticeState, SystemNoticeContext]
}

func (systemNoticeAction) Tag() EffectActionTag { panic("implement me") }

type GameStartAction struct {
	systemNoticeAction
}

func (GameStartAction) Tag() EffectActionTag { return GAME_START_ACTION }

type GameFinishAction struct {
	systemNoticeAction
}

func (GameFinishAction) Tag() EffectActionTag { return GAME_FINISH_ACTION }

type TurnStartAction struct {
	systemNoticeAction
}

func (TurnStartAction) Tag() EffectActionTag { return TURN_START_ACTION }

type TurnEndAction struct {
	systemNoticeAction
}

func (TurnEndAction) Tag() EffectActionTag { return TURN_END_ACTION }

type PlayerWinAction struct {
	systemNoticeAction
}

func (PlayerWinAction) Tag() EffectActionTag { return PLAYER_WIN_ACTION }

type PlayerLoseAction struct {
	systemNoticeAction
}

func (PlayerLoseAction) Tag() EffectActionTag { return PLAYER_LOSE_ACTION }

func (s SystemNoticeContext) IsCanceled() bool {
	return false
}

func (s SystemNoticeContext) Cancel() {}

func (s SystemNoticeContext) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (s *SystemNoticeContext) FromMap(_ map[string]interface{}) error {
	return nil
}

func (s systemNoticeAction) MakeState() interface{} {
	return SystemNoticeState{}
}

func (s systemNoticeAction) Act(state interface{}, beforeAction EffectAction, beforeContext EffectContext) EffectContext {
	return &SystemNoticeContext{}
}

func (s systemNoticeAction) Solve(board *model.Board, state interface{}, context EffectContext) *model.Board {
	return board
}

var (
	GAME_START_ACTION  EffectActionTag = "system.game_start"
	GAME_FINISH_ACTION EffectActionTag = "system.game_finish"
	TURN_START_ACTION  EffectActionTag = "system.turn_start"
	TURN_END_ACTION    EffectActionTag = "system.turn_end"
	PLAYER_WIN_ACTION  EffectActionTag = "system.player_win"
	PLAYER_LOSE_ACTION EffectActionTag = "system.player_lose"
)
