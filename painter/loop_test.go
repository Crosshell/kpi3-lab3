package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr
	l.state = TextureState{Background: color.Black}

	var testOps []string

	l.Start(mockScreen{})
	l.Post(logOp(t, "do white fill", ColorFill{Color: color.White}))
	l.Post(logOp(t, "do green fill", ColorFill{Color: color.RGBA{G: 0xff, A: 0xff}}))
	l.Post(UpdateOp)

	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", ColorFill{Color: color.RGBA{G: 0xff, A: 0xff}}))
	}

	l.Post(operationFunc(func(state TextureState) TextureState {
		testOps = append(testOps, "op 1")
		return state
	}))
	l.Post(operationFunc(func(state TextureState) TextureState {
		testOps = append(testOps, "op 2")
		return state
	}))
	l.Post(operationFunc(func(state TextureState) TextureState {
		testOps = append(testOps, "op 3")
		return state
	}))

	time.Sleep(100 * time.Millisecond)
	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", tr.lastTexture)
	}
	if len(mt.Colors) > 0 {
		lastColor := mt.Colors[len(mt.Colors)-1]
		if _, ok := lastColor.(color.RGBA); !ok || lastColor.(color.RGBA) != (color.RGBA{G: 0xff, A: 0xff}) {
			t.Errorf("Last color is not green, or not RGBA: %+v, Colors: %+v", lastColor, mt.Colors)
		}
	} else {
		t.Error("No colors were filled in mockTexture, but updates should have occurred.")
	}

	if !reflect.DeepEqual(testOps, []string{"op 1", "op 2", "op 3"}) {
		t.Error("Bad order:", testOps)
	}
}

type operationFunc func(state TextureState) TextureState

func (f operationFunc) Do(state TextureState) TextureState {
	return f(state)
}

func logOp(t *testing.T, msg string, op Operation) Operation {
	return operationFunc(func(state TextureState) TextureState {
		t.Log(msg, reflect.TypeOf(op))
		return op.Do(state)
	})
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{Colors: []color.Color{}}, nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}
