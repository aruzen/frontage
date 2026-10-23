package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"frontage/internal/samplegame"
)

func main() {
	seed := flag.Int64("seed", 1, "山札の乱数シード")
	auto := flag.Bool("auto", false, "両側を自動操作して対戦を最後まで進める")
	quiet := flag.Bool("quiet", false, "自動対戦の途中表示を省略する")
	flag.Parse()
	if err := run(os.Stdin, os.Stdout, *seed, *auto, *quiet); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const help = `
play 番号 [対象] [移動先]  カードを使用（例: play 1 B4）
move 元 先                ユニットを上下左右に1マス移動
attack 元 対象            ユニットで攻撃（例: attack B2 C2）
end                       手番を終了してNPCへ
hint                      おすすめの操作を表示（状態は変えない）
board / hand              盤面と手札を再表示
cards                     全カードの説明
auto                      残りの対戦を自動操作
help / quit               ヘルプ / 終了

敵のコアHPを0にすると勝利。召喚は自陣2列の空きマスへ。
各ユニットは自分の手番に移動か攻撃を1回。召喚直後は行動不可。
毎手番マナ全回復・最大値+1（初回4、最大10）。初回以外は1枚ドロー。
手札上限8・ユニット上限6。山札切れはコアに1,2,3…ダメージ。
`

func run(input io.Reader, output io.Writer, seed int64, auto, quiet bool) error {
	g := samplegame.New(seed)
	fmt.Fprintf(output, "Frontage / sample-game / seed=%d\n", seed)
	if auto {
		if err := autoplay(g, output, quiet, false); err != nil {
			return err
		}
		g.RenderResult(output)
		return nil
	}
	fmt.Fprint(output, help)
	g.Render(output)
	scanner := bufio.NewScanner(input)
	for !g.Over() {
		fmt.Fprint(output, "\nsample-game> ")
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch strings.ToLower(line) {
		case "quit", "q":
			return nil
		case "help", "?":
			fmt.Fprint(output, help)
			continue
		case "board", "hand":
			g.Render(output)
			continue
		case "cards":
			samplegame.RenderCatalogue(output)
			continue
		case "hint":
			fmt.Fprintln(output, "おすすめ:", g.BestCommand())
			continue
		case "auto":
			if err := autoplay(g, output, false, false); err != nil {
				return err
			}
			continue
		}
		command, err := samplegame.ParseCommand(line)
		if err == nil {
			err = g.Apply(command)
		}
		if err != nil {
			fmt.Fprintln(output, "入力エラー:", err)
			continue
		}
		printEvents(g, output)
		if !g.Over() && g.Side() == 1 {
			if err := autoplay(g, output, false, true); err != nil {
				return err
			}
		}
		g.Render(output)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	g.RenderResult(output)
	return nil
}

func autoplay(g *samplegame.Game, output io.Writer, quiet, npcOnly bool) error {
	for !g.Over() && (!npcOnly || g.Side() == 1) {
		command := g.BestCommand()
		if !quiet {
			fmt.Fprintf(output, "\n%s> %s\n", samplegame.SideName(g.Side()), command)
		}
		if err := g.Apply(command); err != nil {
			return fmt.Errorf("自動操作 %s: %w", command, err)
		}
		if !quiet {
			printEvents(g, output)
			if !npcOnly && (command.Kind == "end" || g.Over()) {
				g.Render(output)
			}
		}
	}
	return nil
}

func printEvents(g *samplegame.Game, output io.Writer) {
	for _, message := range g.Events {
		fmt.Fprintln(output, " ", message)
	}
}
