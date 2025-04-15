package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

type Operation interface {
	Do(state TextureState) TextureState
}

type OperationList []Operation

func (ol OperationList) Do(state TextureState) TextureState {
	for _, o := range ol {
		state = o.Do(state)
	}
	return state
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(state TextureState) TextureState {
	return state 
}

type OperationFunc func(state TextureState) TextureState

func (f OperationFunc) Do(state TextureState) TextureState {
	return f(state)
}


type ColorFill struct {
	Color color.Color
}

func (op ColorFill) Do(state TextureState) TextureState {
	state.Background = op.Color
	return state
}


type BgRect struct {
	X1, Y1, X2, Y2 float64
}

func (op BgRect) Do(state TextureState) TextureState {
	state.BgRect = &op 
	return state
}

func (op BgRect) Draw(t screen.Texture) {
	bounds := t.Bounds()
	rect := image.Rect(
		int(float64(bounds.Dx())*op.X1),
		int(float64(bounds.Dy())*op.Y1),
		int(float64(bounds.Dx())*op.X2),
		int(float64(bounds.Dy())*op.Y2),
	)
	t.Fill(rect, color.Black, screen.Src)
}


type Figure struct {
	X, Y float64
}

func (op Figure) Do(state TextureState) TextureState {
	state.Figures = append(state.Figures, op) 
	return state
}

func (op Figure) Draw(t screen.Texture) {
	drawFigureShape(t, op.X, op.Y)
}


type Move struct {
	X, Y float64
}

func (op Move) Do(state TextureState) TextureState {
	return state
}

type Reset struct{}

func (op Reset) Do(state TextureState) TextureState {
	state.Background = color.Black 
	state.BgRect = nil            
	state.Figures = nil           
	return state
}


func WhiteFill(state TextureState) TextureState {
	state.Background = color.White
	return state
}

func GreenFill(state TextureState) TextureState {
	state.Background = color.RGBA{G: 0xff, A: 0xff}
	return state
}


func drawFigureShape(t screen.Texture, x, y float64) {
	bounds := t.Bounds()
	centerX := int(float64(bounds.Max.X) * x)
	centerY := int(float64(bounds.Max.Y) * y)

	verticalWidth := 200
	verticalHeight := 50
	horizontalWidth := 50
	horizontalHeight := 200

	mainRect := image.Rect(
		centerX-verticalWidth/2,
		centerY-verticalHeight/2,
		centerX+verticalWidth/2,
		centerY+verticalHeight/2,
	)

	extensionRect := image.Rect(
		centerX-verticalWidth/2-horizontalWidth,
		centerY-horizontalHeight/2,
		centerX-verticalWidth/2,
		centerY+horizontalHeight/2,
	)

	figureColor := color.RGBA{B: 255, A: 255}

	t.Fill(mainRect, figureColor, screen.Src)
	t.Fill(extensionRect, figureColor, screen.Src)
}