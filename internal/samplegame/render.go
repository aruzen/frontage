package samplegame

import (
	"fmt"
	"io"

	"frontage/pkg/engine/model"
)

func (g *Game) Render(out io.Writer) {
	fmt.Fprintf(out, "\n=== 手番 %d: %s ===\n", g.Board.Turn()+1, SideName(g.Side()))
	fmt.Fprintf(out, "あなた: コア %d/%d  マナ %d/%d  手札 %d  山札 %d\n", g.CoreHealth(0), CoreHP, g.Mana(0), g.MaxMana(0), len(g.Hand(0)), g.DeckCount(0))
	fmt.Fprintf(out, "NPC   : コア %d/%d  マナ %d/%d  手札 %d  山札 %d\n\n", g.CoreHealth(1), CoreHP, g.Mana(1), g.MaxMana(1), len(g.Hand(1)), g.DeckCount(1))
	fmt.Fprintln(out, "        A      B      C      D      E")
	for y := 0; y < Size; y++ {
		fmt.Fprintf(out, "  %d  ", y+1)
		for x := 0; x < Size; x++ {
			label := "."
			if piece := g.Board.Entities()[x][y]; piece != nil {
				icon := "C"
				if def, ok := DefinitionFor(piece.LocalizeTag()); ok {
					icon = def.Icon
				}
				owner := "n"
				if g.owner(piece) == 0 {
					owner = "p"
				}
				label = fmt.Sprintf("%s%s%02d", owner, icon, piece.HP())
			}
			fmt.Fprintf(out, "[%5s]", label)
		}
		fmt.Fprintln(out)
	}
	fmt.Fprintln(out, "  p=あなた / n=NPC / C=コア S=前衛兵 A=射手 G=守護者 / 数字=HP")
	for side := 0; side < 2; side++ {
		for _, piece := range g.pieces(side) {
			if isCore(piece) {
				continue
			}
			state := "行動済"
			if piece.UsedActionCost() < piece.MaxActionCost() {
				state = "行動可"
			}
			if g.Frozen(piece.Id()) {
				state = "凍結"
			}
			fmt.Fprintf(out, "  %s %sの%s HP%d/%d 攻撃%d [%s]\n", PointName(piece.Position()), SideName(side), PieceName(piece), piece.HP(), MaxHP(piece), piece.ATK(), state)
		}
	}
	fmt.Fprintln(out, "\nあなたの手札:")
	for i, instance := range g.Hand(0) {
		def, _ := DefinitionFor(instance.LocalizeTag())
		fmt.Fprintf(out, "  %d. [%d] %s — %s\n", i+1, instance.PlayCost()[model.AETHER], def.Name, def.Description)
	}
	if len(g.Hand(0)) == 0 {
		fmt.Fprintln(out, "  （手札なし）")
	}
}

func RenderCatalogue(out io.Writer) {
	fmt.Fprintln(out, "\nカード一覧（各2枚・計24枚）:")
	for _, def := range definitions {
		fmt.Fprintf(out, "  [%d] %s: %s\n", def.Cost, def.Name, def.Description)
	}
}

func (g *Game) RenderResult(out io.Writer) {
	result := "中断"
	if g.Winner == 2 {
		result = "引き分け"
	} else if g.Winner >= 0 {
		result = SideName(g.Winner) + "の勝利"
	}
	fmt.Fprintf(out, "\n結果: %s / 手番%d / コア %d:%d\n", result, g.Board.Turn()+1, g.CoreHealth(0), g.CoreHealth(1))
	fmt.Fprintf(out, "効果解決 %d回 / 攻撃 %d回 / 移動 %d回 / 撃破 %d体 / 使用カード %d種類\n", g.Stats.Effects, g.Stats.Attacks, g.Stats.Moves, g.Stats.Deaths, len(g.Stats.CardsPlayed))
	for _, def := range definitions {
		if n := g.Stats.CardsPlayed[def.Key]; n > 0 {
			fmt.Fprintf(out, "  %s ×%d\n", def.Name, n)
		}
	}
}
