package samplegame

import (
	"frontage/pkg"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

type Effect string

const (
	Summon Effect = "summon"
	Damage Effect = "damage"
	Heal   Effect = "heal"
	Buff   Effect = "buff"
	Storm  Effect = "storm"
	Draw   Effect = "draw"
	Drain  Effect = "drain"
	Haste  Effect = "haste"
	Freeze Effect = "freeze"
	Blink  Effect = "blink"
)

type Definition struct {
	Key, Name, Description string
	Effect                 Effect
	Cost, Targets          int
	HP, ATK, Range         int
	Icon                   string
}

var definitions = []Definition{
	{Key: "soldier", Name: "前衛兵", Description: "HP4/攻撃2。自陣の空きマスに召喚", Effect: Summon, Cost: 1, Targets: 1, HP: 4, ATK: 2, Range: 1, Icon: "S"},
	{Key: "archer", Name: "森の射手", Description: "HP3/攻撃2/射程2。自陣の空きマスに召喚", Effect: Summon, Cost: 2, Targets: 1, HP: 3, ATK: 2, Range: 2, Icon: "A"},
	{Key: "guardian", Name: "石の守護者", Description: "HP7/攻撃2。自陣の空きマスに召喚", Effect: Summon, Cost: 3, Targets: 1, HP: 7, ATK: 2, Range: 1, Icon: "G"},
	{Key: "fireball", Name: "火球", Description: "敵1体（コア可）に4ダメージ", Effect: Damage, Cost: 2, Targets: 1},
	{Key: "heal", Name: "癒やしの泉", Description: "味方1体（コア可）を4回復。初期HPが上限", Effect: Heal, Cost: 1, Targets: 1},
	{Key: "rally", Name: "鼓舞", Description: "味方ユニット1体の攻撃力を2上げる", Effect: Buff, Cost: 1, Targets: 1},
	{Key: "storm", Name: "火炎嵐", Description: "敵全体（コアを含む）に2ダメージ。対象指定なし", Effect: Storm, Cost: 4},
	{Key: "insight", Name: "叡智", Description: "カードを2枚引く（手札上限8）。対象指定なし", Effect: Draw, Cost: 1},
	{Key: "drain", Name: "吸魂", Description: "敵1体に3ダメージ、自分のコアを3回復", Effect: Drain, Cost: 3, Targets: 1},
	{Key: "haste", Name: "疾風", Description: "味方ユニット1体を再行動可能にし、凍結も解除", Effect: Haste, Cost: 1, Targets: 1},
	{Key: "freeze", Name: "氷結", Description: "敵ユニット1体を次の手番に行動不能にする", Effect: Freeze, Cost: 1, Targets: 1},
	{Key: "blink", Name: "転移", Description: "味方ユニットを空きマスへ転移。対象: 移動元 移動先", Effect: Blink, Cost: 1, Targets: 2},
}

func Catalogue() []Definition { return append([]Definition(nil), definitions...) }

func DefinitionFor(key pkg.LocalizeTag) (Definition, bool) {
	for _, def := range definitions {
		if def.Key == string(key) {
			return def, true
		}
	}
	return Definition{}, false
}

func newCard(def Definition) model.Card {
	base := model.NewBaseCard(pkg.LocalizeTag(def.Key), model.Materials{model.AETHER: def.Cost})
	if def.Effect == Summon {
		return card.NewBasePiece(base, def.HP, 1, def.ATK, offsets(1), offsets(def.Range))
	}
	return &spellCard{BaseCard: *base}
}

type spellCard struct{ model.BaseCard }

func (*spellCard) Type() model.CardType { return model.CARD_TYPE_MAGIC }
func (c *spellCard) CardCopy() model.MutableCard {
	return &spellCard{BaseCard: *c.BaseCard.CardCopy().(*model.BaseCard)}
}
func (c *spellCard) CardMirror(id uuid.UUID) model.MutableCard {
	return &spellCard{BaseCard: *c.BaseCard.CardMirror(id).(*model.BaseCard)}
}

func offsets(distance int) []pkg.Point {
	var result []pkg.Point
	for y := -distance; y <= distance; y++ {
		for x := -distance; x <= distance; x++ {
			if d := abs(x) + abs(y); d > 0 && d <= distance {
				result = append(result, pkg.Point{X: x, Y: y})
			}
		}
	}
	return result
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
