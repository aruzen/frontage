package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestInteractiveGameCommands(t *testing.T) {
	input := "play 1 C5\nplay 1 B4\nplay 3 B4\nmove B4 B3\nend\nhint\nquit\n"
	var output bytes.Buffer
	if err := run(strings.NewReader(input), &output, 1, false, false); err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"入力エラー:", "に召喚", "B4 → B3", "NPC>", "おすすめ:"} {
		if !strings.Contains(output.String(), text) {
			t.Errorf("missing %q in output", text)
		}
	}
}

func TestAutomaticGameFinishes(t *testing.T) {
	var output bytes.Buffer
	if err := run(strings.NewReader(""), &output, 42, true, true); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "の勝利") || !strings.Contains(output.String(), "効果解決") || strings.Contains(output.String(), "中断") {
		t.Fatalf("incomplete automatic game: %s", output.String())
	}
}
