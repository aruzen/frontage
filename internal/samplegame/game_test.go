package samplegame

import (
	"fmt"
	"reflect"
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func testPoint(value string) pkg.Point {
	p, err := ParsePoint(value)
	if err != nil {
		panic(err)
	}
	return p
}

func prepareCard(t *testing.T, key string) *Game {
	t.Helper()
	g := New(1)
	def, ok := DefinitionFor(pkg.LocalizeTag(key))
	if !ok {
		t.Fatal(key)
	}
	g.player(0).Hand = model.Cards{newCard(def)}
	g.player(0).SetMaterials(model.Materials{model.AETHER: 10})
	return g
}

func addUnit(t *testing.T, g *Game, side int, key, point string) model.MutablePiece {
	t.Helper()
	def, ok := DefinitionFor(pkg.LocalizeTag(key))
	if !ok {
		t.Fatal(key)
	}
	piece := model.NewBasePiece(uuid.New(), g.player(side), testPoint(point), def.HP, 1, def.ATK, offsets(1), offsets(def.Range))
	piece.SetLocalizeTag(pkg.LocalizeTag(key))
	if !g.Board.SetPiece(piece.Position(), piece) {
		t.Fatal("failed to place fixture")
	}
	return piece
}

func apply(t *testing.T, g *Game, line string) {
	t.Helper()
	command, err := ParseCommand(line)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.Apply(command); err != nil {
		t.Fatalf("%s: %v", line, err)
	}
}

func pieceHP(g *Game, id uuid.UUID) int {
	if piece, ok := g.Board.GetPiece(id); ok {
		return piece.HP()
	}
	return 0
}

func TestEverySampleCardEffect(t *testing.T) {
	for _, key := range []string{"soldier", "archer", "guardian"} {
		t.Run(key, func(t *testing.T) {
			g := prepareCard(t, key)
			before := g.Board
			apply(t, g, "play 1 B4")
			piece := g.at(testPoint("B4"))
			def, _ := DefinitionFor(pkg.LocalizeTag(key))
			if piece == nil || piece.HP() != def.HP || piece.ATK() != def.ATK || piece.LocalizeTag() != pkg.LocalizeTag(key) {
				t.Fatal("summon lost its definition")
			}
			if piece.UsedActionCost() != piece.MaxActionCost() || g.Mana(0) != 10-def.Cost || len(g.Hand(0)) != 0 || len(g.player(0).Graveyard) != 1 {
				t.Fatal("summon cost, card consumption or summoning sickness failed")
			}
			if before.Entities()[1][3] != nil || len(before.Players()[0].(*model.LocalPlayer).Hand) != 1 {
				t.Fatal("command changed the previous board")
			}
		})
	}
	t.Run("fireball", func(t *testing.T) {
		g := prepareCard(t, "fireball")
		apply(t, g, "play 1 C1")
		if g.CoreHealth(1) != 16 || g.Mana(0) != 8 {
			t.Fatal("damage or payment failed")
		}
	})
	t.Run("heal", func(t *testing.T) {
		g := prepareCard(t, "heal")
		core, _ := g.Board.GetPiece(g.CoreIDs[0])
		core.(model.MutablePiece).SetHP(19)
		apply(t, g, "play 1 C5")
		if g.CoreHealth(0) != CoreHP {
			t.Fatal("healing exceeded maximum HP")
		}
	})
	t.Run("rally", func(t *testing.T) {
		g := prepareCard(t, "rally")
		unit := addUnit(t, g, 0, "soldier", "B4")
		apply(t, g, "play 1 B4")
		changed, _ := g.Board.GetPiece(unit.Id())
		if changed.ATK() != 4 || changed.HP() != 4 || unit.ATK() != 2 {
			t.Fatal("buff failed or leaked into previous board")
		}
	})
	t.Run("storm", func(t *testing.T) {
		g := prepareCard(t, "storm")
		enemy := addUnit(t, g, 1, "soldier", "B2")
		ally := addUnit(t, g, 0, "soldier", "B4")
		apply(t, g, "play 1")
		if pieceHP(g, enemy.Id()) != 2 || pieceHP(g, ally.Id()) != 4 || g.CoreHealth(1) != 18 {
			t.Fatal("area effect targeted the wrong pieces")
		}
	})
	t.Run("insight", func(t *testing.T) {
		g := prepareCard(t, "insight")
		deck := g.player(0).MainDeck
		last, second := deck[len(deck)-1].Id(), deck[len(deck)-2].Id()
		apply(t, g, "play 1")
		hand := g.Hand(0)
		if len(hand) != 2 || hand[0].Id() != last || hand[1].Id() != second || g.DeckCount(0) != len(deck)-2 {
			t.Fatal("draw order or count failed")
		}
	})
	t.Run("drain", func(t *testing.T) {
		g := prepareCard(t, "drain")
		enemy := addUnit(t, g, 1, "soldier", "B2")
		core, _ := g.Board.GetPiece(g.CoreIDs[0])
		core.(model.MutablePiece).SetHP(10)
		apply(t, g, "play 1 B2")
		if pieceHP(g, enemy.Id()) != 1 || g.CoreHealth(0) != 13 {
			t.Fatal("compound damage/healing failed")
		}
	})
	t.Run("haste", func(t *testing.T) {
		g := prepareCard(t, "haste")
		unit := addUnit(t, g, 0, "soldier", "B4")
		unit.SetUsedActionCost(1)
		g.frozen[unit.Id()] = 0
		apply(t, g, "play 1 B4")
		if g.Frozen(unit.Id()) {
			t.Fatal("haste did not thaw the unit")
		}
		apply(t, g, "move B4 B3")
		if g.at(testPoint("B3")) == nil {
			t.Fatal("unit could not act again")
		}
	})
	t.Run("freeze", func(t *testing.T) {
		g := prepareCard(t, "freeze")
		unit := addUnit(t, g, 1, "soldier", "B2")
		apply(t, g, "play 1 B2")
		apply(t, g, "end")
		command, _ := ParseCommand("move B2 B3")
		if err := g.Apply(command); err == nil {
			t.Fatal("frozen unit acted on its next turn")
		}
		apply(t, g, "end")
		apply(t, g, "end")
		if g.Frozen(unit.Id()) {
			t.Fatal("freeze survived two owner turns")
		}
		apply(t, g, "move B2 B3")
	})
	t.Run("blink", func(t *testing.T) {
		g := prepareCard(t, "blink")
		unit := addUnit(t, g, 0, "soldier", "B4")
		unit.SetUsedActionCost(1)
		apply(t, g, "play 1 B4 E2")
		moved, ok := g.Board.GetPiece(unit.Id())
		if !ok || moved.Position() != testPoint("E2") || moved.UsedActionCost() != 1 || g.at(testPoint("B4")) != nil {
			t.Fatal("teleport lost identity, position or action budget")
		}
	})
}

func TestInvalidActionsDoNotSpendOrChangeState(t *testing.T) {
	g := New(1)
	addUnit(t, g, 0, "soldier", "B4")
	addUnit(t, g, 1, "soldier", "B2")
	for _, line := range []string{"play 99", "play 1 C5", "play 1 B2", "play 1 B4", "move B2 B3", "move B4 B2", "attack B4 C5", "attack B4 A1"} {
		t.Run(line, func(t *testing.T) {
			before, mana, hand := g.Board, g.Mana(0), len(g.Hand(0))
			command, err := ParseCommand(line)
			if err != nil {
				t.Fatal(err)
			}
			if err := g.Apply(command); err == nil {
				t.Fatal("illegal action succeeded")
			}
			if g.Board != before || g.Mana(0) != mana || len(g.Hand(0)) != hand || g.Stats.Effects != 0 {
				t.Fatal("invalid action changed live state")
			}
		})
	}
	t.Run("insufficient mana", func(t *testing.T) {
		g := prepareCard(t, "fireball")
		g.player(0).SetMaterials(nil)
		before := g.Board
		cmd, _ := ParseCommand("play 1 C1")
		if g.Apply(cmd) == nil || g.Board != before || len(g.Hand(0)) != 1 {
			t.Fatal("unpaid card was consumed")
		}
	})
}

func TestCombatDeathAndCoreVictory(t *testing.T) {
	g := New(1)
	attacker := addUnit(t, g, 0, "soldier", "B3")
	defender := addUnit(t, g, 1, "archer", "B2")
	defender.SetHP(2)
	apply(t, g, "attack B3 B2")
	if _, exists := g.Board.GetPiece(defender.Id()); exists {
		t.Fatal("dead defender stayed on board")
	}
	current, _ := g.Board.GetPiece(attacker.Id())
	if current.UsedActionCost() != 1 || g.Stats.Attacks != 1 || g.Stats.Deaths != 1 {
		t.Fatal("combat bookkeeping failed")
	}
	g = prepareCard(t, "fireball")
	core, _ := g.Board.GetPiece(g.CoreIDs[1])
	core.(model.MutablePiece).SetHP(4)
	apply(t, g, "play 1 C1")
	if g.Winner != 0 {
		t.Fatal("destroying core did not win")
	}
	before := g.Board
	if g.Apply(Command{Kind: "end"}) == nil || g.Board != before {
		t.Fatal("finished game accepted another action")
	}
}

func TestHandLimitAndFatigue(t *testing.T) {
	g := prepareCard(t, "insight")
	def, _ := DefinitionFor("soldier")
	for len(g.Hand(0)) < HandLimit {
		g.player(0).Hand = append(g.player(0).Hand, newCard(def))
	}
	deck := g.DeckCount(0)
	apply(t, g, "play 1")
	if len(g.Hand(0)) != HandLimit || g.DeckCount(0) != deck-1 {
		t.Fatal("draw exceeded hand limit")
	}
	g = New(1)
	g.player(1).MainDeck = nil
	g.turns[1] = 1
	for len(g.Hand(1)) < HandLimit {
		g.player(1).Hand = append(g.player(1).Hand, newCard(def))
	}
	apply(t, g, "end")
	if g.CoreHealth(1) != 19 {
		t.Fatal("empty deck did not cause fatigue with a full hand")
	}
	apply(t, g, "end")
	apply(t, g, "end")
	if g.CoreHealth(1) != 17 {
		t.Fatal("fatigue did not increase")
	}
}

func TestAIPreviewDoesNotMutateLiveGame(t *testing.T) {
	g := New(42)
	unit := addUnit(t, g, 0, "soldier", "B4")
	unit.SetUsedActionCost(1)
	enemy := addUnit(t, g, 1, "soldier", "B2")
	g.frozen[enemy.Id()] = 1
	before, mana, events := g.Board, g.Mana(0), append([]string(nil), g.Events...)
	stats := g.Stats
	stats.CardsPlayed = map[string]int{}
	first := g.BestCommand()
	second := g.BestCommand()
	if first.String() != second.String() || g.validate(first) != nil {
		t.Fatal("AI was unstable or generated an illegal command")
	}
	if g.Board != before || g.Mana(0) != mana || g.frozen[enemy.Id()] != 1 || unit.UsedActionCost() != 1 || !reflect.DeepEqual(g.Stats, stats) || !reflect.DeepEqual(g.Events, events) {
		t.Fatal("preview changed the real game")
	}
}

func TestAutoplayCompletesWithRealCombat(t *testing.T) {
	seen := make(map[string]bool)
	for _, seed := range []int64{1, 2, 3, 7, 42, 99} {
		t.Run(fmt.Sprint(seed), func(t *testing.T) {
			g := New(seed)
			for actions := 0; !g.Over() && actions < MaxTurns*33; actions++ {
				command := g.BestCommand()
				if err := g.Apply(command); err != nil {
					t.Fatalf("%s: %v", command, err)
				}
			}
			if !g.Over() || g.Board.Turn() > MaxTurns || g.Stats.Attacks == 0 || g.Stats.Deaths == 0 {
				t.Fatalf("match did not complete with combat: turn=%d winner=%d stats=%+v", g.Board.Turn(), g.Winner, g.Stats)
			}
			for key := range g.Stats.CardsPlayed {
				seen[key] = true
			}
			t.Logf("seed=%d turns=%d winner=%d effects=%d card types=%d", seed, g.Board.Turn()+1, g.Winner, g.Stats.Effects, len(g.Stats.CardsPlayed))
		})
	}
	if len(seen) != len(definitions) {
		t.Errorf("autoplay covered %d/%d card types: %v", len(seen), len(definitions), seen)
	}
}
