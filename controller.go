package orbengine

import "github.com/veandco/go-sdl2/sdl"

/*
Destroy the controller object.
*/
func (c *Controller) Destroy() {
	c.window.Destroy()
	// TODO can I call this anywhere? is this... ok?
	sdl.Quit()
}
