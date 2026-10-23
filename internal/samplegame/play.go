package samplegame

import (
	"fmt"

	"frontage/pkg"
	"frontage/pkg/engine/impl/action/card_action"
	"frontage/pkg/engine/impl/action/piece_action"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/model"
)

func (g *Game) validate(command Command) error {
	if g.Over() {
		return fmt.Errorf("対戦は終了しています")
	}
	if command.Kind == "end" && len(command.Targets) == 0 {
		return nil
	}
	if g.turnActions >= 32 {
		return fmt.Errorf("この手番の操作上限です。endしてください")
	}
	for _, p := range command.Targets {
		if !inBounds(p) {
			return fmt.Errorf("盤面の外です")
		}
	}
	switch command.Kind {
	case "play":
		hand := g.player(g.Side()).Hand
		if command.Card < 0 || command.Card >= len(hand) {
			return fmt.Errorf("その手札番号はありません")
		}
		instance := hand[command.Card]
		def, ok := DefinitionFor(instance.LocalizeTag())
		if !ok {
			return fmt.Errorf("未対応のカードです")
		}
		if len(command.Targets) != def.Targets {
			return fmt.Errorf("%sの対象は%d個です", def.Name, def.Targets)
		}
		if instance.PlayCost()[model.AETHER] > g.Mana(g.Side()) {
			return fmt.Errorf("マナが足りません")
		}
		if def.Effect == Summon {
			p := command.Targets[0]
			if g.at(p) != nil {
				return fmt.Errorf("召喚先は空きマスを指定してください")
			}
			if (g.Side() == 0 && p.Y < Size-2) || (g.Side() == 1 && p.Y > 1) {
				return fmt.Errorf("召喚できるのは自陣の2列です")
			}
			if len(g.pieces(g.Side()))-1 >= UnitLimit {
				return fmt.Errorf("ユニット上限は%d体です", UnitLimit)
			}
			return nil
		}
		if def.Targets == 0 {
			return nil
		}
		target := g.at(command.Targets[0])
		if target == nil {
			return fmt.Errorf("対象のマスに駒がありません")
		}
		switch def.Effect {
		case Damage, Drain, Freeze:
			if g.owner(target) == g.Side() {
				return fmt.Errorf("敵の駒を指定してください")
			}
		case Heal, Buff, Haste, Blink:
			if g.owner(target) != g.Side() {
				return fmt.Errorf("自分の駒を指定してください")
			}
		}
		if (def.Effect == Buff || def.Effect == Haste || def.Effect == Freeze || def.Effect == Blink) && isCore(target) {
			return fmt.Errorf("コアは対象にできません")
		}
		switch def.Effect {
		case Heal:
			if target.HP() >= MaxHP(target) {
				return fmt.Errorf("対象のHPは満タンです")
			}
		case Haste:
			if target.UsedActionCost() == 0 && !g.Frozen(target.Id()) {
				return fmt.Errorf("対象はすでに行動可能です")
			}
		case Freeze:
			if g.Frozen(target.Id()) {
				return fmt.Errorf("対象はすでに凍結しています")
			}
		case Blink:
			if g.at(command.Targets[1]) != nil {
				return fmt.Errorf("転移先は空きマスを指定してください")
			}
		}
		return nil
	case "move", "attack":
		if len(command.Targets) != 2 {
			return fmt.Errorf("移動元と対象を指定してください")
		}
		from, to := command.Targets[0], command.Targets[1]
		actor := g.at(from)
		if actor == nil || g.owner(actor) != g.Side() || isCore(actor) {
			return fmt.Errorf("自分のユニットを指定してください")
		}
		if actor.UsedActionCost() >= actor.MaxActionCost() || g.Frozen(actor.Id()) {
			return fmt.Errorf("そのユニットは行動できません")
		}
		offsets := actor.LegalMoves()
		if command.Kind == "move" {
			if g.at(to) != nil {
				return fmt.Errorf("移動先が埋まっています")
			}
		} else {
			target := g.at(to)
			if target == nil || g.owner(target) == g.Side() {
				return fmt.Errorf("敵の駒を攻撃してください")
			}
			offsets = actor.AttackRanges()
		}
		for _, delta := range offsets {
			if from.X+delta.X == to.X && from.Y+delta.Y == to.Y {
				return nil
			}
		}
		return fmt.Errorf("移動・攻撃の範囲外です")
	}
	return fmt.Errorf("不明な操作です")
}

func (g *Game) play(command Command) error {
	side := g.Side()
	instance := g.player(side).Hand[command.Card]
	def, _ := DefinitionFor(instance.LocalizeTag())
	g.Board = g.Board.Next()
	player := g.player(side)
	player.SetMaterials(model.Materials{model.AETHER: g.Mana(side) - instance.PlayCost()[model.AETHER]})
	if def.Effect == Summon {
		current, _ := player.Hand.FindById(instance.Id())
		state := card_action.NewPieceSummonActionState(player, model.DECK_TYPE_HAND, current.(card.Piece), command.Targets[0])
		if err := g.resolve(card_action.PieceSummonAction{}, &state); err != nil {
			return err
		}
	} else {
		if !player.Hand.RemoveById(instance.Id()) {
			return fmt.Errorf("手札が見つかりません")
		}
	}
	g.player(side).Graveyard = append(g.player(side).Graveyard, instance.CardCopy())
	g.Stats.CardsPlayed[def.Key]++
	g.Events = append(g.Events, fmt.Sprintf("%sが「%s」を使用（マナ%d）", SideName(side), def.Name, def.Cost))
	var target model.Piece
	if len(command.Targets) > 0 {
		target = g.at(command.Targets[0])
	}
	switch def.Effect {
	case Summon:
		g.Events = append(g.Events, def.Name+"を"+PointName(command.Targets[0])+"に召喚（次の手番から行動）")
	case Damage:
		g.Events = append(g.Events, fmt.Sprintf("%sの%sに4ダメージ", PointName(command.Targets[0]), PieceName(target)))
		return g.changeHP(target.Id(), -4)
	case Heal:
		g.Events = append(g.Events, fmt.Sprintf("%sの%sを%d回復", PointName(command.Targets[0]), PieceName(target), min(4, MaxHP(target)-target.HP())))
		return g.heal(target.Id(), 4)
	case Buff:
		g.Events = append(g.Events, PointName(command.Targets[0])+"の"+PieceName(target)+"の攻撃力+2")
		return g.perform(piece_action.PieceATKIncreaseAction{}, map[string]interface{}{"piece_id": target.Id().String(), "value": 2})
	case Storm:
		g.Events = append(g.Events, "敵全体に2ダメージ")
		for _, enemy := range g.pieces(1 - side) {
			if err := g.changeHP(enemy.Id(), -2); err != nil {
				return err
			}
		}
	case Draw:
		return g.draw(side, 2)
	case Drain:
		g.Events = append(g.Events, fmt.Sprintf("%sの%sに3ダメージ、コアを%d回復", PointName(command.Targets[0]), PieceName(target), min(3, CoreHP-g.CoreHealth(side))))
		if err := g.changeHP(target.Id(), -3); err != nil {
			return err
		}
		return g.heal(g.CoreIDs[side], 3)
	case Haste:
		g.Events = append(g.Events, PointName(command.Targets[0])+"の"+PieceName(target)+"が再行動可能になった")
		delete(g.frozen, target.Id())
		return g.setActionCost(target.Id(), 0)
	case Freeze:
		g.Events = append(g.Events, PointName(command.Targets[0])+"の"+PieceName(target)+"は次の手番に行動不能")
		g.frozen[target.Id()] = 1
		return g.setActionCost(target.Id(), target.MaxActionCost())
	case Blink:
		return g.move(command.Targets[0], command.Targets[1], 0)
	}
	return nil
}

func (g *Game) move(from, to pkg.Point, cost int) error {
	piece := g.at(from)
	err := g.perform(piece_action.PieceMoveAction{}, map[string]interface{}{
		"piece_id": piece.Id().String(), "from": pkg.PointToMap(from), "to": pkg.PointToMap(to), "action_cost": cost,
	})
	if err != nil {
		return err
	}
	g.Stats.Moves++
	g.Events = append(g.Events, fmt.Sprintf("%sの%sが %s → %s", SideName(g.owner(piece)), PieceName(piece), PointName(from), PointName(to)))
	return nil
}

func (g *Game) attack(from, to pkg.Point) error {
	actor, target := g.at(from), g.at(to)
	err := g.perform(piece_action.PieceAttackAction{}, map[string]interface{}{
		"piece_id": actor.Id().String(), "point": pkg.PointToMap(to), "value": actor.ATK(), "action_cost": 1,
	})
	if err != nil {
		return err
	}
	g.Stats.Attacks++
	g.Events = append(g.Events, fmt.Sprintf("%sの%sが%sの%sへ%dダメージ（%s → %s）", SideName(g.owner(actor)), PieceName(actor), SideName(g.owner(target)), PieceName(target), actor.ATK(), PointName(from), PointName(to)))
	return nil
}
