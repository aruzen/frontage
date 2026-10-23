package samplegame

import (
	"fmt"
	"math/rand"

	"frontage/pkg"
	"frontage/pkg/engine/impl/action/piece_action"
	"frontage/pkg/engine/impl/action/player_action"
	actionregister "frontage/pkg/engine/impl/action/register"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

const (
	Size      = 5
	CoreHP    = 20
	HandLimit = 8
	UnitLimit = 6
	MaxTurns  = 100
)

type Statistics struct {
	CardsPlayed                     map[string]int
	Effects, Moves, Attacks, Deaths int
}

type Game struct {
	Board       *model.Board
	CoreIDs     [2]uuid.UUID
	Winner      int // -1: ongoing, 0/1: winner, 2: draw.
	Events      []string
	Stats       Statistics
	maxMana     [2]int
	turns       [2]int
	fatigue     [2]int
	frozen      map[uuid.UUID]int // 1: next owner turn, 0: current owner turn.
	turnActions int
}

func New(seed int64) *Game {
	actionregister.Init()
	random := rand.New(rand.NewSource(seed))
	var players [2]model.Player
	for side := range players {
		var deck model.Cards
		for _, def := range definitions {
			deck = append(deck, newCard(def), newCard(def))
		}
		player := model.NewLocalPlayer(uuid.New(), deck, nil)
		for _, key := range []string{"soldier", "archer", "fireball", "haste", "insight"} {
			for i, instance := range player.MainDeck {
				if string(instance.LocalizeTag()) == key {
					player.Hand = append(player.Hand, instance)
					player.MainDeck = append(player.MainDeck[:i], player.MainDeck[i+1:]...)
					break
				}
			}
		}
		random.Shuffle(len(player.MainDeck), func(i, j int) { player.MainDeck[i], player.MainDeck[j] = player.MainDeck[j], player.MainDeck[i] })
		player.SetMaterials(model.Materials{model.AETHER: 4})
		players[side] = player
	}
	g := &Game{
		Board:  model.NewBoard(model.NewBoardInfo(pkg.Size{Width: Size, Height: Size}, model.GENERATION_STRATEGY_CHAIN), players),
		Winner: -1, maxMana: [2]int{4, 4}, turns: [2]int{1, 0}, frozen: make(map[uuid.UUID]int),
		Stats: Statistics{CardsPlayed: make(map[string]int)}, Events: []string{"対戦開始。あなたの手番です。"},
	}
	for side, y := range []int{Size - 1, 0} {
		point := pkg.Point{X: Size / 2, Y: y}
		core := model.NewBasePiece(uuid.New(), players[side], point, CoreHP, 0, 0, nil, nil)
		core.SetLocalizeTag("core")
		core.SetMaxActionCost(0)
		g.CoreIDs[side] = core.Id()
		g.Board.SetPiece(point, core)
	}
	return g
}

func SideName(side int) string {
	if side == 0 {
		return "あなた"
	}
	return "NPC"
}
func (g *Game) Side() int                 { return g.Board.Phase() }
func (g *Game) Over() bool                { return g.Winner >= 0 }
func (g *Game) Mana(side int) int         { return g.player(side).Materials()[model.AETHER] }
func (g *Game) MaxMana(side int) int      { return g.maxMana[side] }
func (g *Game) Hand(side int) model.Cards { return append(model.Cards(nil), g.player(side).Hand...) }
func (g *Game) DeckCount(side int) int    { return len(g.player(side).MainDeck) }
func (g *Game) Frozen(id uuid.UUID) bool  { _, exists := g.frozen[id]; return exists }
func (g *Game) player(side int) *model.LocalPlayer {
	return g.Board.Players()[side].(*model.LocalPlayer)
}
func (g *Game) CoreHealth(side int) int {
	if core, ok := g.Board.GetPiece(g.CoreIDs[side]); ok {
		return max(0, core.HP())
	}
	return 0
}

func (g *Game) clone() *Game {
	copy := *g
	copy.Board = g.Board.Sandbox()
	copy.Events = append([]string(nil), g.Events...)
	copy.frozen = make(map[uuid.UUID]int, len(g.frozen))
	for id, turns := range g.frozen {
		copy.frozen[id] = turns
	}
	copy.Stats.CardsPlayed = make(map[string]int, len(g.Stats.CardsPlayed))
	for key, count := range g.Stats.CardsPlayed {
		copy.Stats.CardsPlayed[key] = count
	}
	return &copy
}

// Apply stages the complete command, including costs and secondary effects.
// Invalid input or a failed effect leaves the live game unchanged.
func (g *Game) Apply(command Command) error {
	if err := g.validate(command); err != nil {
		return err
	}
	trial := g.clone()
	trial.Events = nil
	if err := trial.execute(command); err != nil {
		return err
	}
	*g = *trial
	return nil
}

func (g *Game) execute(command Command) error {
	var err error
	switch command.Kind {
	case "play":
		err = g.play(command)
	case "move":
		err = g.move(command.Targets[0], command.Targets[1], 1)
	case "attack":
		err = g.attack(command.Targets[0], command.Targets[1])
	case "end":
		err = g.endTurn()
	}
	if err != nil {
		return err
	}
	if command.Kind != "end" {
		g.turnActions++
	}
	if err := g.removeDead(); err != nil {
		return err
	}
	first, second := g.CoreHealth(0), g.CoreHealth(1)
	switch {
	case first <= 0 && second <= 0:
		g.Winner = 2
	case first <= 0:
		g.Winner = 1
	case second <= 0:
		g.Winner = 0
	case g.Board.Turn() >= MaxTurns:
		g.Winner = 2
	}
	if g.Over() {
		result := "引き分け"
		if g.Winner < 2 {
			result = SideName(g.Winner) + "の勝利"
		}
		g.Events = append(g.Events, result+"！")
	}
	return nil
}

func (g *Game) perform(action logic.EffectAction, input map[string]interface{}) error {
	state := action.MakeState(g.Board, input)
	if state == nil {
		return fmt.Errorf("効果 %s のStateを作成できません", action.Tag())
	}
	return g.resolve(action, state)
}

func (g *Game) resolve(action logic.EffectAction, state logic.ActionState) error {
	system := logic.NewEventSystem(g.Board)
	system.Emit(logic.NewEffectEvent(action, state))
	system.Resolve()
	if system.Board == nil || len(system.AppliedEffects) == 0 {
		return fmt.Errorf("効果 %s を解決できません", action.Tag())
	}
	g.Board = system.Board
	g.Stats.Effects += len(system.AppliedEffects)
	return nil
}

func (g *Game) changeHP(id uuid.UUID, amount int) error {
	if amount == 0 {
		return nil
	}
	var effect logic.EffectAction = piece_action.PieceHPIncreaseAction{}
	if amount < 0 {
		effect = piece_action.PieceHPDecreaseAction{}
		amount = -amount
	}
	return g.perform(effect, map[string]interface{}{"piece_id": id.String(), "value": amount})
}

func (g *Game) heal(id uuid.UUID, amount int) error {
	piece, exists := g.Board.GetPiece(id)
	if !exists {
		return fmt.Errorf("回復対象が見つかりません")
	}
	return g.changeHP(id, min(amount, max(0, MaxHP(piece)-piece.HP())))
}

func (g *Game) setActionCost(id uuid.UUID, value int) error {
	return g.perform(actionCostEffect{}, map[string]interface{}{"piece_id": id.String(), "value": value})
}

func (g *Game) draw(side, count int) error {
	for i := 0; i < count; i++ {
		player := g.player(side)
		if len(player.MainDeck) == 0 {
			g.fatigue[side]++
			if err := g.changeHP(g.CoreIDs[side], -g.fatigue[side]); err != nil {
				return err
			}
			g.Events = append(g.Events, fmt.Sprintf("%s: 山札切れでコアに%dダメージ", SideName(side), g.fatigue[side]))
			continue
		}
		if len(player.Hand) >= HandLimit {
			g.Events = append(g.Events, SideName(side)+"は手札上限のため追加ドローなし")
			break
		}
		state := player_action.NewDrawActionState(player)
		if err := g.resolve(player_action.DrawMainDeckAction{}, state); err != nil {
			return err
		}
		g.Events = append(g.Events, SideName(side)+"がカードを1枚引いた")
	}
	return nil
}

func (g *Game) endTurn() error {
	g.Board = g.Board.AdvanceTurn()
	g.turnActions = 0
	side := g.Side()
	g.turns[side]++
	g.maxMana[side] = min(10, 3+g.turns[side])
	g.player(side).SetMaterials(model.Materials{model.AETHER: g.maxMana[side]})
	for _, piece := range g.pieces(side) {
		if isCore(piece) {
			continue
		}
		cost := 0
		if turns, exists := g.frozen[piece.Id()]; exists {
			if turns > 0 {
				cost = piece.MaxActionCost()
				g.frozen[piece.Id()] = 0
			} else {
				delete(g.frozen, piece.Id())
			}
		}
		if err := g.setActionCost(piece.Id(), cost); err != nil {
			return err
		}
	}
	g.Events = append(g.Events, fmt.Sprintf("%sの手番（マナ%d）", SideName(side), g.maxMana[side]))
	if g.turns[side] > 1 {
		return g.draw(side, 1)
	}
	return nil
}

func (g *Game) removeDead() error {
	dead := append(g.pieces(0), g.pieces(1)...)
	for _, piece := range dead {
		if piece.HP() > 0 {
			continue
		}
		state := piece_action.NewPieceDeathActionState(piece)
		if err := g.resolve(piece_action.PieceDeathAction{}, &state); err != nil {
			return err
		}
		delete(g.frozen, piece.Id())
		g.Stats.Deaths++
		g.Events = append(g.Events, fmt.Sprintf("%sの%s（%s）が撃破された", SideName(g.owner(piece)), PieceName(piece), PointName(piece.Position())))
	}
	return nil
}

func (g *Game) pieces(side int) []model.Piece {
	var result []model.Piece
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			piece := g.Board.Entities()[x][y]
			if piece != nil && g.owner(piece) == side {
				result = append(result, piece)
			}
		}
	}
	return result
}

func (g *Game) owner(piece model.Piece) int {
	if piece.Owner() != nil {
		for side, player := range g.Board.Players() {
			if player.ID() == piece.Owner().ID() {
				return side
			}
		}
	}
	return -1
}

func (g *Game) at(point pkg.Point) model.Piece {
	if !inBounds(point) {
		return nil
	}
	return g.Board.Entities()[point.X][point.Y]
}
func inBounds(p pkg.Point) bool     { return p.X >= 0 && p.X < Size && p.Y >= 0 && p.Y < Size }
func isCore(piece model.Piece) bool { return piece.LocalizeTag() == "core" }
func MaxHP(piece model.Piece) int {
	if isCore(piece) {
		return CoreHP
	}
	def, _ := DefinitionFor(piece.LocalizeTag())
	return def.HP
}
func PieceName(piece model.Piece) string {
	if isCore(piece) {
		return "コア"
	}
	if def, ok := DefinitionFor(piece.LocalizeTag()); ok {
		return def.Name
	}
	return string(piece.LocalizeTag())
}
