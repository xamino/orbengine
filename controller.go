package orbengine

import "github.com/veandco/go-sdl2/sdl"

/*
Destroy the controller object.
*/
func (c *Controller) Destroy() {
	c.renderer.Destroy()
	c.window.Destroy()
	// TODO can I call this anywhere? is this... ok?
	sdl.Quit()
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
	c.renderer.Clear()
	for _, d := range c.drawables {
		d.Draw(c.renderer)
	}
	c.renderer.Present()
}
