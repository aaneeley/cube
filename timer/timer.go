package timer

import (
	"fmt"
	"time"

	"github.com/aaneeley/cube/model"
	"github.com/aaneeley/cube/window"
)

type FrameTimer struct {
	times [100]float64 // store in milliseconds
	index int
	count int
}

func (ft *FrameTimer) Add(dt time.Duration) {
	ms := float64(dt.Microseconds()) / 1000.0
	ft.times[ft.index] = ms
	ft.index = (ft.index + 1) % len(ft.times)
	if ft.count < len(ft.times) {
		ft.count++
	}
}

func (ft *FrameTimer) Average() float64 {
	if ft.count == 0 {
		return 0
	}
	var sum float64
	for i := range ft.count {
		sum += ft.times[i]
	}
	return sum / float64(ft.count)
}

func (ft *FrameTimer) DrawToBuf(window *window.Window) {
	window.Buffer.WriteString(model.PosCode(window.X, window.Y+window.H+2))
	window.Buffer.WriteString("                               ")
	window.Buffer.WriteString(model.PosCode(window.X, window.Y+window.H+2))
	window.Buffer.WriteString(fmt.Sprintf("Avg frame time: \033[92m%.1fms\033[0m", ft.Average()))
}
