package samplegame

import (
	"fmt"
	"strconv"
	"strings"

	"frontage/pkg"
)

type Command struct {
	Kind    string
	Card    int // Zero-based hand index.
	Targets []pkg.Point
}

func ParseCommand(line string) (Command, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return Command{}, fmt.Errorf("コマンドを入力してください")
	}
	command := Command{Kind: strings.ToLower(fields[0])}
	args := fields[1:]
	switch command.Kind {
	case "play":
		if len(args) < 1 {
			return Command{}, fmt.Errorf("play 手札番号 [対象] [移動先]")
		}
		number, err := strconv.Atoi(args[0])
		if err != nil || number < 1 {
			return Command{}, fmt.Errorf("手札番号は1以上の整数です")
		}
		command.Card = number - 1
		args = args[1:]
		if len(args) > 2 {
			return Command{}, fmt.Errorf("対象は最大2つです")
		}
	case "move", "attack":
		if len(args) != 2 {
			return Command{}, fmt.Errorf("%s 移動元 対象（例: %s B4 B3）", command.Kind, command.Kind)
		}
	case "end", "pass":
		command.Kind = "end"
		if len(args) != 0 {
			return Command{}, fmt.Errorf("end に引数は不要です")
		}
	default:
		return Command{}, fmt.Errorf("不明なコマンドです。help で操作を確認できます")
	}
	for _, arg := range args {
		point, err := ParsePoint(arg)
		if err != nil {
			return Command{}, err
		}
		command.Targets = append(command.Targets, point)
	}
	return command, nil
}

func ParsePoint(value string) (pkg.Point, error) {
	value = strings.ToUpper(value)
	if len(value) != 2 || value[0] < 'A' || value[0] >= 'A'+Size || value[1] < '1' || value[1] >= '1'+Size {
		return pkg.Point{}, fmt.Errorf("座標はA1〜E5です: %q", value)
	}
	return pkg.Point{X: int(value[0] - 'A'), Y: int(value[1] - '1')}, nil
}

func PointName(point pkg.Point) string { return fmt.Sprintf("%c%d", 'A'+point.X, point.Y+1) }

func (c Command) String() string {
	parts := []string{c.Kind}
	if c.Kind == "play" {
		parts = append(parts, strconv.Itoa(c.Card+1))
	}
	for _, point := range c.Targets {
		parts = append(parts, PointName(point))
	}
	return strings.Join(parts, " ")
}
