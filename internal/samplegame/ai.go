package samplegame

import (
	"frontage/pkg"
	"frontage/pkg/engine/model"
)

// BestCommand uses one-ply evaluation on isolated boards. It never changes the
// live game, and equal scores keep the stable hand/coordinate ordering.
func (g *Game) BestCommand() Command {
	best := Command{Kind: "end"}
	if g.Over() || g.turnActions >= 32 {
		return best
	}
	side := g.Side()
	bestScore := g.score(side) + 0.01
	for _, command := range g.candidates() {
		if g.validate(command) != nil {
			continue
		}
		trial := g.clone()
		trial.Events = nil
		if trial.execute(command) != nil {
			continue
		}
		if score := trial.score(side); score > bestScore {
			best, bestScore = command, score
		}
	}
	return best
}

func (g *Game) candidates() []Command {
	var result []Command
	var points []pkg.Point
	for y := 0; y < Size; y++ {
		for x := 0; x < Size; x++ {
			points = append(points, pkg.Point{X: x, Y: y})
		}
	}
	seen := make(map[pkg.LocalizeTag]bool)
	for i, instance := range g.player(g.Side()).Hand {
		if seen[instance.LocalizeTag()] {
			continue
		}
		seen[instance.LocalizeTag()] = true
		def, _ := DefinitionFor(instance.LocalizeTag())
		if instance.PlayCost()[model.AETHER] > g.Mana(g.Side()) {
			continue
		}
		switch def.Targets {
		case 0:
			result = append(result, Command{Kind: "play", Card: i})
		case 1:
			for _, point := range points {
				result = append(result, Command{Kind: "play", Card: i, Targets: []pkg.Point{point}})
			}
		case 2:
			for _, piece := range g.pieces(g.Side()) {
				if isCore(piece) {
					continue
				}
				for _, point := range points {
					result = append(result, Command{Kind: "play", Card: i, Targets: []pkg.Point{piece.Position(), point}})
				}
			}
		}
	}
	for _, piece := range g.pieces(g.Side()) {
		if isCore(piece) || piece.UsedActionCost() >= piece.MaxActionCost() || g.Frozen(piece.Id()) {
			continue
		}
		for _, delta := range piece.LegalMoves() {
			to := pkg.Point{X: piece.Position().X + delta.X, Y: piece.Position().Y + delta.Y}
			result = append(result, Command{Kind: "move", Targets: []pkg.Point{piece.Position(), to}})
		}
		for _, target := range g.pieces(1 - g.Side()) {
			result = append(result, Command{Kind: "attack", Targets: []pkg.Point{piece.Position(), target.Position()}})
		}
	}
	return result
}

func (g *Game) score(side int) float64 {
	if g.Winner == side {
		return 100000
	}
	if g.Winner == 1-side {
		return -100000
	}
	if g.Winner == 2 {
		return 0
	}
	score := float64(g.CoreHealth(side)-g.CoreHealth(1-side)) * 2
	score += float64(len(g.player(side).Hand)-len(g.player(1-side).Hand)) * 0.35
	score += float64(g.Mana(side)) * 0.05
	for owner := 0; owner < 2; owner++ {
		goal := pkg.Point{X: Size / 2, Y: 0}
		if owner == 1 {
			goal.Y = Size - 1
		}
		for _, piece := range g.pieces(owner) {
			if isCore(piece) {
				continue
			}
			distance := abs(piece.Position().X-goal.X) + abs(piece.Position().Y-goal.Y)
			def, _ := DefinitionFor(piece.LocalizeTag())
			value := 1 + float64(piece.HP())*0.35 + float64(piece.ATK())*0.8 + float64(8-distance)*0.75 + float64(def.Range-1)*1.2
			if owner == side {
				if piece.UsedActionCost() < piece.MaxActionCost() && !g.Frozen(piece.Id()) {
					value += 0.6
				}
				score += value
			} else {
				if !g.Frozen(piece.Id()) {
					value += 0.8
				}
				score -= value
			}
		}
	}
	return score
}
