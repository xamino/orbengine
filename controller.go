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
	case Actionable:
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
	wrapped := &Renderer{renderer: renderer}
	// viewport is used to check if we even have to draw something
	viewport := &sdl.Rect{}
	renderer.GetViewport(viewport)
	// run for all entities
	for _, entity := range c.entities {
		// TODO update everything else BEFORE draw except state thingy
		// execute actions
		if actionEntity, actionable := entity.(Actionable); actionable {
			actionEntity.Action()
		}
		// determine draw method. Texture is preferred over Render!
		eR, renderable := entity.(Renderable)
		eD, drawable := entity.(Drawable)
		if drawable || renderable {
			// drawbale or renderable both are placeable
			e, _ := entity.(Placeable)
			// TODO this can be precalculated and reused if not changed, plus multithreaded for all entities
			// TODO apply scale / unit
			worldBound := &sdl.Rect{
				X: e.Position().X - int32(e.Width()/2),
				Y: e.Position().Y - int32(e.Height()/2),
				W: int32(e.Width()),
				H: int32(e.Height())}
			// check if we even need to draw it
			if _, intersects := viewport.Intersect(worldBound); !intersects {
				continue
			}
			var text *sdl.Texture
			if drawable {
				// create texture to draw to
				text, err = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STATIC,
					e.Width(), e.Height())
				if err != nil {
					log.Println("renderer.CreateTexture error:", err)
					continue
				}
				// allow entity to draw to texture FIXME: reuse?
				eD.Texture(text)
			} else {
				text, _ = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET,
					e.Width(), e.Height())
				renderer.SetRenderTarget(text)
				eR.Render(wrapped)
				renderer.SetRenderTarget(nil)
			}
			// TODO last nil is rotation center, use offset to calculate
			renderer.CopyEx(text, nil, worldBound, e.Rotation(), nil, 0)
		}
	} // for entities
	renderer.Present()
}
