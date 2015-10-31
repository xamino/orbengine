package orbengine

import (
	"log"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

/*
Destroy the controller object.
*/
func (c *Controller) Destroy() {
	c.window.Destroy()
	// TODO can I call this anywhere? is this... ok?
	sdl.Quit()
	// MUST unlock OS thread
	runtime.UnlockOSThread()
}

/*
AddEntity will add the entity to the controller. If the entity does not match
any orbengine interface, an error is returned.
*/
func (c *Controller) AddEntity(id string, entity interface{}) error {
	// TODO check if entity already exists? do we allow adding entities multiple times?
	switch entity.(type) {
	case Drawable:
		c.drawables[id], _ = entity.(Drawable)
	default:
		return ErrMissingComponents
	}
	return nil
}

/*
Iterate TODO
*/
func (c *Controller) Iterate() {
	renderer, err := c.window.GetRenderer()
	if err != nil {
		log.Println("window.GetRenderer error:", err, renderer)
		return
	}
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	for _, d := range c.drawables {
		d.Draw(renderer)
	}
	renderer.Present()
}
