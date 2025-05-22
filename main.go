package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/aaneeley/cube/geometry"
	"github.com/aaneeley/cube/model"
	"github.com/aaneeley/cube/timer"
	win "github.com/aaneeley/cube/window"
	"golang.org/x/term"
)

const (
	paddingX = 50
	paddingY = 10

	rotationSpeed = 1.5
	sideLength    = 45
)

var frameTimer timer.FrameTimer

func main() {
	termW, termH, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	model.ClearTerminal()
	model.HideCursor()
	defer model.ShowCursor()

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Exit handlers
	quit := make(chan struct{})
	pause := make(chan struct{})
	go func() {
		buf := make([]byte, 1)
		for {
			_, err := os.Stdin.Read(buf)
			if err != nil {
				continue
			}
			if buf[0] == 'q' {
				quit <- struct{}{}
				return
			} else if buf[0] == 'p' {
				pause <- struct{}{}
			}
		}
	}()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		quit <- struct{}{}
	}()

	window := win.CreateWindow(paddingX, paddingY, termW-(paddingX*2-1), termH-(paddingY*2-1))
	window.Render()

	cube := geometry.NewCube(sideLength)
	cube.SetOrigin(model.NewVec3(float64(termW)/2, float64(termH)/2, 50))

	targetFrameDuration := time.Second / 60
	lastFrameTime := time.Now()
	rot := 0.5
	running := true
	animate := true
	for running {
		select {
		case <-sig:
			running = false
		case <-quit:
			running = false
		case <-pause:
			animate = !animate
		default:
			frameStart := time.Now()
			delta := frameStart.Sub(lastFrameTime).Seconds()
			lastFrameTime = frameStart

			cube.SetRotation(model.NewVec3(rot*0.3, -rot, 0.5))
			cube.DrawToBuf(window, float64(termH)/2)

			window.Render()

			frameTimer.DrawToBuf(window)
			frameTimer.Add(time.Since(frameStart))

			if animate {
				rot += delta * float64(rotationSpeed)
			}
			if sleepTime := targetFrameDuration - time.Since(frameStart); sleepTime > 0 {
				time.Sleep(sleepTime)
			}
		}
	}
	model.ClearTerminal()
}
