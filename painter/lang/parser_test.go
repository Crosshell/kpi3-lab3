package lang_test

import (
	"strings"
	"testing"

	"image/color"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
)

func TestParser_Parse(t *testing.T) {
	input := `
white
green
update
bgrect 10 20 30 40
figure 50 60
move 70 80
reset
invalid
bgrect wrong args
`

	parser := lang.Parser{}
	ops, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ops) != 7 {
		t.Fatalf("expected 7 operations, got %d", len(ops))
	}

	if cf, ok := ops[0].(painter.ColorFill); !ok || cf.Color != color.White {
		t.Errorf("expected ColorFill with white, got %+v", ops[0])
	}

	if cf, ok := ops[1].(painter.ColorFill); !ok || cf.Color != (color.RGBA{G: 0xff, A: 0xff}) {
		t.Errorf("expected ColorFill with green, got %+v", ops[1])
	}

	if ops[2] != painter.UpdateOp {
		t.Errorf("expected UpdateOp, got %+v", ops[2])
	}

	if bg, ok := ops[3].(painter.BgRect); !ok || bg.X1 != 10 || bg.Y1 != 20 || bg.X2 != 30 || bg.Y2 != 40 {
		t.Errorf("expected BgRect with coords, got %+v", ops[3])
	}

	if fig, ok := ops[4].(painter.Figure); !ok || fig.X != 50 || fig.Y != 60 {
		t.Errorf("expected Figure, got %+v", ops[4])
	}

	if mv, ok := ops[5].(painter.Move); !ok || mv.X != 70 || mv.Y != 80 {
		t.Errorf("expected Move, got %+v", ops[5])
	}

	if _, ok := ops[6].(painter.Reset); !ok {
		t.Errorf("expected Reset, got %+v", ops[6])
	}
}
