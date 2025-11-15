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

type GameStartAction struct {
	systemNoticeAction
}

type GameFinishAction struct {
	systemNoticeAction
}

type TurnStartAction struct {
	systemNoticeAction
}

type TurnEndAction struct {
	systemNoticeAction
}

type PlayerWinAction struct {
	systemNoticeAction
}

type PlayerLoseAction struct {
	systemNoticeAction
}

func (s SystemNoticeContext) IsCanceled() bool {
	return false
}

func (s SystemNoticeContext) Cancel() {}

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
	GAME_START_ACTION  = GameStartAction{}
	GAME_FINISH_ACTION = GameFinishAction{}
	TURN_START_ACTION  = TurnStartAction{}
	TURN_END_ACTION    = TurnEndAction{}
	PLAYER_WIN_ACTION  = PlayerWinAction{}
	PLAYER_LOSE_ACTION = PlayerLoseAction{}
)
