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
	ENTITY_SUMMON_ACTION           EffectActionTag = "entityaction.entity.summon"
	ENTITY_MOVE_ACTION             EffectActionTag = "entityaction.entity.move"
	ENTITY_ATTACK_ACTION           EffectActionTag = "entityaction.entity.attack"
	ENTITY_INVASION_ACTION         EffectActionTag = "entityaction.entity.invasion"
	ENTITY_HP_INCREASE_ACTION      EffectActionTag = "entityaction.entity.hp_increase"
	ENTITY_HP_DECREASE_ACTION      EffectActionTag = "entityaction.entity.hp_decrease"
	ENTITY_HP_FIX_ACTION           EffectActionTag = "entityaction.entity.hp_fix"
	ENTITY_MP_INCREASE_ACTION      EffectActionTag = "entityaction.entity.mp_increase"
	ENTITY_MP_DECREASE_ACTION      EffectActionTag = "entityaction.entity.mp_decrease"
	ENTITY_MP_FIX_ACTION           EffectActionTag = "entityaction.entity.mp_fix"
	ENTITY_ATK_INCREASE_ACTION     EffectActionTag = "entityaction.entity.atk_increase"
	ENTITY_ATK_DECREASE_ACTION     EffectActionTag = "entityaction.entity.atk_decrease"
	ENTITY_ATK_FIX_ACTION          EffectActionTag = "entityaction.entity.atk_fix"
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
