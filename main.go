package orbengine

import (
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

/*
Controller is the central struct from which the entire engine can be controlled.
*/
type Controller struct {
	window       *sdl.Window              // window is the window instance.
	wrapped      *Renderer                // wrapped is a locked down version of the SDL renderer.
	renderable   *ordered                 // renderable is a sorted list of entities that are drawn.
	entities     map[string]interface{}   // entities is a list of all active entities.
	keyreceivers map[string][]*keyReceive // keyreceivers holds structs that trigger on key presses.
	textCache    map[string]*sdl.Texture  // textCache is a list of cached textures used to draw entities.
}

type keyReceive struct {
	lastState bool
	onPress   func()
	onRelease func()
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
	// prepare basic Controller stuff
	c := &Controller{}
	c.renderable = makeOrdered()
	c.entities = make(map[string]interface{})
	c.keyreceivers = make(map[string][]*keyReceive)
	c.textCache = make(map[string]*sdl.Texture)
	// window
	window, err := sdl.CreateWindow("Orbiting", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}
	c.window = window
	// create renderer (reference kept by window, no need to do it ourselves)
	renderer, err := sdl.CreateRenderer(c.window, -1, 0)
	if err != nil {
		return nil, err
	}
	// however, we DO want to keep just a single instance of the wrapper
	c.wrapped = &Renderer{renderer: renderer}
	return c, nil
}
