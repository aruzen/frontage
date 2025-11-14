package action

import (
	cardaction "frontage/pkg/engine/impl/action/card_action"
	entityaction "frontage/pkg/engine/impl/action/entity_action"
	"frontage/pkg/engine/logic"
)

var (
	CARD_PIECE_HP_INCREASE_ACTION  logic.EffectAction = cardaction.PieceHPIncreaseAction{}
	CARD_PIECE_HP_DECREASE_ACTION  logic.EffectAction = cardaction.PieceHPDecreaseAction{}
	CARD_PIECE_HP_FIX_ACTION       logic.EffectAction = cardaction.PieceHPFixAction{}
	CARD_PIECE_MP_INCREASE_ACTION  logic.EffectAction = cardaction.PieceMPIncreaseAction{}
	CARD_PIECE_MP_DECREASE_ACTION  logic.EffectAction = cardaction.PieceMPDecreaseAction{}
	CARD_PIECE_MP_FIX_ACTION       logic.EffectAction = cardaction.PieceMPFixAction{}
	CARD_PIECE_ATK_INCREASE_ACTION logic.EffectAction = cardaction.PieceATKIncreaseAction{}
	CARD_PIECE_ATK_DECREASE_ACTION logic.EffectAction = cardaction.PieceATKDecreaseAction{}
	CARD_PIECE_ATK_FIX_ACTION      logic.EffectAction = cardaction.PieceATKFixAction{}
	ENTITY_SUMMON_ACTION           logic.EffectAction = entityaction.EntitySummonAction{}
	ENTITY_MOVE_ACTION             logic.EffectAction = entityaction.EntityMoveAction{}
	ENTITY_ATTACK_ACTION           logic.EffectAction = entityaction.EntityAttackAction{}
	ENTITY_INVASION_ACTION         logic.EffectAction = entityaction.EntityInvasionAction{}
	ENTITY_HP_INCREASE_ACTION      logic.EffectAction = entityaction.EntityHPIncreaseAction{}
	ENTITY_HP_DECREASE_ACTION      logic.EffectAction = entityaction.EntityHPDecreaseAction{}
	ENTITY_HP_FIX_ACTION           logic.EffectAction = entityaction.EntityHPFixAction{}
	ENTITY_MP_INCREASE_ACTION      logic.EffectAction = entityaction.EntityMPIncreaseAction{}
	ENTITY_MP_DECREASE_ACTION      logic.EffectAction = entityaction.EntityMPDecreaseAction{}
	ENTITY_MP_FIX_ACTION           logic.EffectAction = entityaction.EntityMPFixAction{}
	ENTITY_ATK_INCREASE_ACTION     logic.EffectAction = entityaction.EntityATKIncreaseAction{}
	ENTITY_ATK_DECREASE_ACTION     logic.EffectAction = entityaction.EntityATKDecreaseAction{}
	ENTITY_ATK_FIX_ACTION          logic.EffectAction = entityaction.EntityATKFixAction{}
)
