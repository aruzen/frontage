package action

import (
	"fmt"
	"frontage/pkg/engine/logic"
)

type ActionTag string
type EffectActionTag ActionTag
type ModifierActionTag ActionTag

var (
	CARD_PIECE_HP_INCREASE_ACTION  EffectActionTag = "cardaction.piece.hp_increase"
	CARD_PIECE_HP_DECREASE_ACTION  EffectActionTag = "cardaction.piece.hp_decrease"
	CARD_PIECE_HP_FIX_ACTION       EffectActionTag = "cardaction.piece.hp_fix"
	CARD_PIECE_MP_INCREASE_ACTION  EffectActionTag = "cardaction.piece.mp_increase"
	CARD_PIECE_MP_DECREASE_ACTION  EffectActionTag = "cardaction.piece.mp_decrease"
	CARD_PIECE_MP_FIX_ACTION       EffectActionTag = "cardaction.piece.mp_fix"
	CARD_PIECE_ATK_INCREASE_ACTION EffectActionTag = "cardaction.piece.atk_increase"
	CARD_PIECE_ATK_DECREASE_ACTION EffectActionTag = "cardaction.piece.atk_decrease"
	CARD_PIECE_ATK_FIX_ACTION      EffectActionTag = "cardaction.piece.atk_fix"
	ENTITY_SUMMON_ACTION           EffectActionTag = "pieceaction.piece.summon"
	ENTITY_MOVE_ACTION             EffectActionTag = "pieceaction.piece.move"
	ENTITY_ATTACK_ACTION           EffectActionTag = "pieceaction.piece.attack"
	ENTITY_INVASION_ACTION         EffectActionTag = "pieceaction.piece.invasion"
	ENTITY_HP_INCREASE_ACTION      EffectActionTag = "pieceaction.piece.hp_increase"
	ENTITY_HP_DECREASE_ACTION      EffectActionTag = "pieceaction.piece.hp_decrease"
	ENTITY_HP_FIX_ACTION           EffectActionTag = "pieceaction.piece.hp_fix"
	ENTITY_MP_INCREASE_ACTION      EffectActionTag = "pieceaction.piece.mp_increase"
	ENTITY_MP_DECREASE_ACTION      EffectActionTag = "pieceaction.piece.mp_decrease"
	ENTITY_MP_FIX_ACTION           EffectActionTag = "pieceaction.piece.mp_fix"
	ENTITY_ATK_INCREASE_ACTION     EffectActionTag = "pieceaction.piece.atk_increase"
	ENTITY_ATK_DECREASE_ACTION     EffectActionTag = "pieceaction.piece.atk_decrease"
	ENTITY_ATK_FIX_ACTION          EffectActionTag = "pieceaction.piece.atk_fix"
)

var effectActionTable map[EffectActionTag]logic.EffectAction

func Register(tag EffectActionTag, action logic.EffectAction) error {
	if _, ok := effectActionTable[tag]; ok {
		return fmt.Errorf("effect action tag %s is already registered", tag)
	}
	effectActionTable[tag] = action
	return nil
}

func FindActionEffect(tag EffectActionTag) logic.EffectAction {
	return effectActionTable[tag]
}
