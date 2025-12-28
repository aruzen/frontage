package logic

import (
	"frontage/pkg"
	"frontage/pkg/engine/model"
)

type SystemNoticeState struct {
}

func (s SystemNoticeState) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (s *SystemNoticeState) FromMap(_ *model.Board, _ map[string]interface{}) error {
	return nil
}

type SystemNoticeContext struct {
}

type systemNoticeAction struct {
	BaseAction[SystemNoticeState, SystemNoticeContext]
}

func (systemNoticeAction) Tag() EffectActionTag { return SYSTEM_NOTICE_ACTION_TAG }
func (systemNoticeAction) LocalizeTag() pkg.LocalizeTag {
	return pkg.LocalizeTag(SYSTEM_NOTICE_ACTION_TAG)
}

type GameStartAction struct {
	systemNoticeAction
}

func (GameStartAction) Tag() EffectActionTag           { return GAME_START_ACTION_TAG }
func (a GameStartAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }

type GameFinishAction struct {
	systemNoticeAction
}

func (GameFinishAction) Tag() EffectActionTag           { return GAME_FINISH_ACTION_TAG }
func (a GameFinishAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }

type TurnStartAction struct {
	systemNoticeAction
}

func (TurnStartAction) Tag() EffectActionTag           { return TURN_START_ACTION_TAG }
func (a TurnStartAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }

type TurnEndAction struct {
	systemNoticeAction
}

func (TurnEndAction) Tag() EffectActionTag           { return TURN_END_ACTION_TAG }
func (a TurnEndAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }

type PlayerWinAction struct {
	systemNoticeAction
}

func (PlayerWinAction) Tag() EffectActionTag           { return PLAYER_WIN_ACTION_TAG }
func (a PlayerWinAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }

type PlayerLoseAction struct {
	systemNoticeAction
}

func (PlayerLoseAction) Tag() EffectActionTag           { return PLAYER_LOSE_ACTION_TAG }
func (a PlayerLoseAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }

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

func (s systemNoticeAction) Act(state ActionState, beforeAction EffectAction, beforeContext EffectContext) (EffectContext, Summary) {
	return &SystemNoticeContext{}, Summary{}
}

func (s systemNoticeAction) Solve(board *model.Board, state ActionState, context EffectContext) (*model.Board, Summary) {
	return board, Summary{}
}

var (
	SYSTEM_NOTICE_ACTION_TAG EffectActionTag = "system.notice"
	GAME_START_ACTION_TAG    EffectActionTag = "system.game_start"
	GAME_FINISH_ACTION_TAG   EffectActionTag = "system.game_finish"
	TURN_START_ACTION_TAG    EffectActionTag = "system.turn_start"
	TURN_END_ACTION_TAG      EffectActionTag = "system.turn_end"
	PLAYER_WIN_ACTION_TAG    EffectActionTag = "system.player_win"
	PLAYER_LOSE_ACTION_TAG   EffectActionTag = "system.player_lose"

	GAME_START_ACTION  = GameStartAction{}
	GAME_FINISH_ACTION = GameFinishAction{}
	TURN_START_ACTION  = TurnStartAction{}
	TURN_END_ACTION    = TurnEndAction{}
	PLAYER_WIN_ACTION  = PlayerWinAction{}
	PLAYER_LOSE_ACTION = PlayerLoseAction{}
)
