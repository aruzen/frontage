package action

import (
	"fmt"
	"frontage/pkg/engine/logic"
)

var (
	CARD_PIECE_HP_INCREASE_ACTION  logic.EffectActionTag = "cardaction.piece.hp_increase"
	CARD_PIECE_HP_DECREASE_ACTION  logic.EffectActionTag = "cardaction.piece.hp_decrease"
	CARD_PIECE_HP_FIX_ACTION       logic.EffectActionTag = "cardaction.piece.hp_fix"
	CARD_PIECE_MP_INCREASE_ACTION  logic.EffectActionTag = "cardaction.piece.mp_increase"
	CARD_PIECE_MP_DECREASE_ACTION  logic.EffectActionTag = "cardaction.piece.mp_decrease"
	CARD_PIECE_MP_FIX_ACTION       logic.EffectActionTag = "cardaction.piece.mp_fix"
	CARD_PIECE_ATK_INCREASE_ACTION logic.EffectActionTag = "cardaction.piece.atk_increase"
	CARD_PIECE_ATK_DECREASE_ACTION logic.EffectActionTag = "cardaction.piece.atk_decrease"
	CARD_PIECE_ATK_FIX_ACTION      logic.EffectActionTag = "cardaction.piece.atk_fix"
	CARD_PIECE_SUMMON_ACTION       logic.EffectActionTag = "cardaction.piece.summon"
	ENTITY_SUMMON_ACTION           logic.EffectActionTag = "pieceaction.piece.summon"
	ENTITY_MOVE_ACTION             logic.EffectActionTag = "pieceaction.piece.move"
	ENTITY_ATTACK_ACTION           logic.EffectActionTag = "pieceaction.piece.attack"
	ENTITY_INVASION_ACTION         logic.EffectActionTag = "pieceaction.piece.invasion"
	ENTITY_HP_INCREASE_ACTION      logic.EffectActionTag = "pieceaction.piece.hp_increase"
	ENTITY_HP_DECREASE_ACTION      logic.EffectActionTag = "pieceaction.piece.hp_decrease"
	ENTITY_HP_FIX_ACTION           logic.EffectActionTag = "pieceaction.piece.hp_fix"
	ENTITY_MP_INCREASE_ACTION      logic.EffectActionTag = "pieceaction.piece.mp_increase"
	ENTITY_MP_DECREASE_ACTION      logic.EffectActionTag = "pieceaction.piece.mp_decrease"
	ENTITY_MP_FIX_ACTION           logic.EffectActionTag = "pieceaction.piece.mp_fix"
	ENTITY_ATK_INCREASE_ACTION     logic.EffectActionTag = "pieceaction.piece.atk_increase"
	ENTITY_ATK_DECREASE_ACTION     logic.EffectActionTag = "pieceaction.piece.atk_decrease"
	ENTITY_ATK_FIX_ACTION          logic.EffectActionTag = "pieceaction.piece.atk_fix"
)

var (
	MOVE_CANCEL_MODIFY_ACTION logic.ModifyActionTag = "pieceaction.move.cancel"
)

var (
	GAME_START_ACTION  logic.EffectActionTag = logic.GAME_START_ACTION_TAG
	GAME_FINISH_ACTION logic.EffectActionTag = logic.GAME_FINISH_ACTION_TAG
	TURN_START_ACTION  logic.EffectActionTag = logic.TURN_START_ACTION_TAG
	TURN_END_ACTION    logic.EffectActionTag = logic.TURN_END_ACTION_TAG
	PLAYER_WIN_ACTION  logic.EffectActionTag = logic.PLAYER_WIN_ACTION_TAG
	PLAYER_LOSE_ACTION logic.EffectActionTag = logic.PLAYER_LOSE_ACTION_TAG
)

var effectActionTable map[logic.EffectActionTag]logic.EffectAction

func Register(tag logic.EffectActionTag, action logic.EffectAction) error {
	if effectActionTable == nil {
		effectActionTable = make(map[logic.EffectActionTag]logic.EffectAction)
	}
	if _, ok := effectActionTable[tag]; ok {
		return fmt.Errorf("effect action tag %s is already registered", tag)
	}
	effectActionTable[tag] = action
	return nil
}

func FindActionEffect(tag logic.EffectActionTag) logic.EffectAction {
	return effectActionTable[tag]
}
