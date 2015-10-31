package orbengine

import "github.com/veandco/go-sdl2/sdl"

/*
Controller is the central struct from which the entire engine can be controlled.
*/
type Controller struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

/*
Build creates a Controller struct.
*/
func Build() (*Controller, error) {
	sdl.Init(sdl.INIT_EVERYTHING)
	c := &Controller{}
	// window
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}
	c.window = window
	// renderer
	// FIXME: check if OR is correct (or is ^ XOR)
	renderer, err := sdl.CreateRenderer(c.window, -1, sdl.RENDERER_ACCELERATED^sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return nil, err
	}
	c.renderer = renderer
	return c, nil
}
