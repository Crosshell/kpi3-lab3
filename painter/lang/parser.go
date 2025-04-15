package lang

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"image/color"
)

type Parser struct{}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	var ops []painter.Operation

	for scanner.Scan() {
		commandLine := scanner.Text()
		if len(strings.TrimSpace(commandLine)) == 0 {
			continue
		}

		fields := strings.Fields(commandLine)
		if len(fields) == 0 {
			continue
		}

		command := fields[0]
		args := fields[1:]

		switch command {
		case "white":
			ops = append(ops, painter.ColorFill{Color: color.White})
		case "green":
			ops = append(ops, painter.ColorFill{Color: color.RGBA{G: 0xff, A: 0xff}})
		case "update":
			ops = append(ops, painter.UpdateOp)
		case "bgrect":
			if len(args) != 4 {
				continue
			}
			x1, err1 := strconv.ParseFloat(args[0], 64)
			y1, err2 := strconv.ParseFloat(args[1], 64)
			x2, err3 := strconv.ParseFloat(args[2], 64)
			y2, err4 := strconv.ParseFloat(args[3], 64)
			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				continue
			}
			ops = append(ops, painter.BgRect{X1: x1, Y1: y1, X2: x2, Y2: y2})
		case "figure":
			if len(args) != 2 {
				continue
			}
			x, err1 := strconv.ParseFloat(args[0], 64)
			y, err2 := strconv.ParseFloat(args[1], 64)
			if err1 != nil || err2 != nil {
				continue
			}
			ops = append(ops, painter.Figure{X: x, Y: y})
		case "move":
			if len(args) != 2 {
				continue
			}
			x, err1 := strconv.ParseFloat(args[0], 64)
			y, err2 := strconv.ParseFloat(args[1], 64)
			if err1 != nil || err2 != nil {
				continue
			}
			ops = append(ops, painter.Move{X: x, Y: y})
		case "reset":
			ops = append(ops, painter.Reset{})
		}
	}

	return ops, nil
}