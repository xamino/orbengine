package orbengine

import (
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

/*
Controller is the central struct from which the entire engine can be controlled.
*/
type Controller struct {
	window    *sdl.Window
	drawables map[string]Drawable
}

/*
Build creates a Controller struct. NOTE: locks the OS thread. Make sure to call
all Controller.func from the exact same goroutine, otherwise the SDL context
may crash!
*/
func Build() (*Controller, error) {
	// MUST guarantee controller always runs in same OS thread
	runtime.LockOSThread()
	sdl.Init(sdl.INIT_EVERYTHING)
	c := &Controller{
		drawables: make(map[string]Drawable)}
	// window
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}
	c.window = window
	// create renderer (reference kept by window, no need to do it ourselves)
	_, err = sdl.CreateRenderer(c.window, -1, 0)
	if err != nil {
		return nil, err
	}
	return c, nil
}
