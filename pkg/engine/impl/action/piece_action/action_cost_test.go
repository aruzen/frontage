package piece_action

import (
	"strings"
	"testing"

	"frontage/pkg"
	"github.com/google/uuid"
)

func TestActionCostParsing(t *testing.T) {
	var moveState PieceMoveActionState
	var attackState PieceAttackActionState
	var moveContext PieceMoveActionContext
	var attackContext PieceAttackActionContext
	for _, target := range []struct {
		name    string
		restore func(map[string]interface{}) error
		cost    func() int
	}{
		{"move state", func(m map[string]interface{}) error { return moveState.FromMap(nil, m) }, func() int { return moveState.ActionCost() }},
		{"attack state", func(m map[string]interface{}) error { return attackState.FromMap(nil, m) }, func() int { return attackState.ActionCost() }},
		{"move context", moveContext.FromMap, func() int { return moveContext.ActionCost }},
		{"attack context", attackContext.FromMap, func() int { return attackContext.ActionCost }},
	} {
		t.Run(target.name, func(t *testing.T) {
			for _, tc := range []struct {
				name    string
				value   interface{}
				omitted bool
				want    int
				invalid bool
			}{
				{name: "integer", value: 2, want: 2},
				{name: "JSON number", value: float64(3), want: 3},
				{name: "zero", value: 0, want: 0},
				{name: "omitted", omitted: true, want: 0},
				{name: "string", value: "2", invalid: true},
				{name: "boolean", value: true, invalid: true},
				{name: "null", value: nil, invalid: true},
				{name: "object", value: map[string]interface{}{}, invalid: true},
			} {
				t.Run(tc.name, func(t *testing.T) {
					input := map[string]interface{}{
						"piece_id": uuid.New().String(), "point": pkg.PointToMap(pkg.Point{}),
						"from": pkg.PointToMap(pkg.Point{}), "to": pkg.PointToMap(pkg.Point{}),
						"value": 1, "action_cost": 7,
					}
					if err := target.restore(input); err != nil {
						t.Fatal(err)
					}
					if tc.omitted {
						delete(input, "action_cost")
					} else {
						input["action_cost"] = tc.value
					}
					err := target.restore(input)
					if tc.invalid {
						if err == nil || !strings.Contains(err.Error(), "action_cost") {
							t.Fatalf("expected an action_cost error, got %v", err)
						}
						return
					}
					if err != nil || target.cost() != tc.want {
						t.Fatalf("cost = %d, error = %v; want %d", target.cost(), err, tc.want)
					}
				})
			}
		})
	}
}
