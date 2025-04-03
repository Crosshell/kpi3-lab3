package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

type Operation interface {
	Do(t screen.Texture) (ready bool)
}

type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool {
	return true
}

type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

type BgRect struct {
	X1, Y1, X2, Y2 float64
}

func (op BgRect) Do(t screen.Texture) bool {
	bounds := t.Bounds()
	rect := image.Rect(
		int(float64(bounds.Dx())*op.X1),
		int(float64(bounds.Dy())*op.Y1),
		int(float64(bounds.Dx())*op.X2),
		int(float64(bounds.Dy())*op.Y2),
	)
	t.Fill(rect, color.Black, screen.Src)
	return false
}

type Figure struct {
	X, Y float64
}

func (op Figure) Do(t screen.Texture) bool {
	drawFigure(t, op.X, op.Y)
	return false
}

type Move struct {
	X, Y float64
}

func (op Move) Do(t screen.Texture) bool {
	// Move operation needs to be handled by the Loop's state
	// This is just a placeholder implementation
	return false
}

type Reset struct{}

func (op Reset) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), color.Black, screen.Src)
	return true
}

func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

func drawFigure(t screen.Texture, x, y float64) {
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
			centerX-verticalWidth/2 - horizontalWidth, 
			centerY-horizontalHeight/2,
			centerX-verticalWidth/2, 
			centerY+horizontalHeight/2,
	)
	
	figureColor := color.RGBA{B: 255, A: 255}
	
	t.Fill(mainRect, figureColor, screen.Src)
	t.Fill(extensionRect, figureColor, screen.Src)
}