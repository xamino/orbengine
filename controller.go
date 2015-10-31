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
	_, exists := c.entities[id]
	if exists {
		return ErrEntityIDExists
	}
	switch entity.(type) {
	case Drawable:
	case Boundable:
	case Actionable:
	case Positionable:
	default:
		return ErrMissingComponents
	}
	c.entities[id] = entity
	return nil
}

/*
Iterate TODO
*/
func (c *Controller) Iterate() {
	// prepare to draw
	renderer, err := c.window.GetRenderer()
	if err != nil {
		log.Println("window.GetRenderer error:", err, renderer)
		return
	}
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()
	// viewport is used to check if we even have to draw something
	viewport := &sdl.Rect{}
	renderer.GetViewport(viewport)
	// run for all entities
	for _, entity := range c.entities {
		// TODO update everything else BEFORE: Action, ...
		// execute actions
		if actionEntity, actionable := entity.(Actionable); actionable {
			actionEntity.Action()
		}
		// draw entities with bounds check if available
		boundEntity, boundable := entity.(Boundable)
		if drawEntity, drawable := entity.(Drawable); drawable {
			if boundable {
				_, intersects := viewport.Intersect(boundEntity.Bounds())
				// if it doesn't intersect we don't need to even draw it, so continue to next entity
				if !intersects {
					continue
				}
				// if it DOES intersect, keep going to draw
			}
			// draw if applicable
			drawEntity.Draw(renderer)
		}
	} // for entities
	renderer.Present()
}
