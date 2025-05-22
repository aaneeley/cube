package window

import (
	"os"
	"strings"

	"github.com/aaneeley/cube/model"
)

type Window struct {
	Buffer     *strings.Builder
	X, Y, W, H int
}

const (
	horizontal = "─"
	vertical   = "│"
	tl         = "╭"
	tr         = "╮"
	bl         = "╰"
	br         = "╯"
)

func (w *Window) resetBuffer() {
	w.Buffer.Reset()

}

func (w *Window) Render() {
	os.Stdout.Write([]byte(w.Buffer.String()))
	w.resetBuffer()
}

func CreateWindow(x, y, width, height int) *Window {
	bb := &strings.Builder{}
	bb.WriteString("\033[94m")
	// Top and bottom
	for curX := x - 1; curX < x+width+1; curX++ {
		bb.WriteString(model.PosCode(curX, y-1))
		bb.WriteString(horizontal)
		bb.WriteString(model.PosCode(curX, y+height+1))
		bb.WriteString(horizontal)
	}
	// Left and right
	for curY := y - 1; curY <= y+height+1; curY++ {
		bb.WriteString(model.PosCode(x-1, curY))
		bb.WriteString(vertical)
		bb.WriteString(model.PosCode(x+width+1, curY))
		bb.WriteString(vertical)
	}
	// Corners
	bb.WriteString(model.PosCode(x-1, y-1))
	bb.WriteString(tl)
	bb.WriteString(model.PosCode(x+width+1, y-1))
	bb.WriteString(tr)
	bb.WriteString(model.PosCode(x+width+1, y+height+1))
	bb.WriteString(br)
	bb.WriteString(model.PosCode(x-1, y+height+1))
	bb.WriteString(bl)

	bb.WriteString("\033[0m")

	return &Window{
		Buffer: bb,
		X:      x,
		Y:      y,
		W:      width,
		H:      height,
	}
}
