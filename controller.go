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
func (c *Controller) AddEntity(entity Identifier) error {
	// check if entity already exists -> id must be unique
	_, exists := c.entities[entity.Identification()]
	if exists {
		return ErrEntityIDExists
	}
	// include is the flag for entites that can be drawn
	var include bool
	// ensure that interface fullfills at least ONE interface
	switch entity.(type) {
	case Drawable:
		include = true
	case Renderable:
		include = true
	case Actionable:
	case Placeable: // Renderable and Drawable fullfull this
	default:
		return ErrMissingComponents
	}
	// if include is true, include the entity in the ordered draw list
	if include {
		// include entity in sorted renderable
		p, valid := entity.(Placeable)
		if !valid {
			// shouldn't happen because of above switch, so catch
			log.Fatal("include was triggered for non-includeable entity!")
		}
		// append: this sorts the new entity into the complete list
		err := c.renderable.append(p)
		if err != nil {
			return err
		}
	}
	// if we reach this -> add
	c.entities[entity.Identification()] = entity
	return nil
}

/*
RegisterKey allows entities to receive key presses by registering the functions
to execute when a key is either pressed or released.
*/
func (c *Controller) RegisterKey(key string, onPress, onRelease func()) error {
	// check if key is valid
	code := sdl.GetKeyFromName(key)
	// FIXME determine how we can check for illegal values
	log.Println("DEBUG: illegal code? :", code)
	// make keyReceive struct
	receive := &keyReceive{
		lastState: false,
		onPress:   onPress,
		onRelease: onRelease}
	_, exists := c.keyreceivers[key]
	if !exists {
		// if new key --> create array with first chan
		c.keyreceivers[key] = []*keyReceive{receive}
	} else {
		// if other entities already registered -> append chan
		c.keyreceivers[key] = append(c.keyreceivers[key], receive)
	}
	return nil
}

/*
Iterate advances the controller by one step.
*/
func (c *Controller) Iterate() {
	// first handle all the events
	c.handleEvents()
	// then step each entity
	c.iterateEntities()
}

func (c *Controller) handleEvents() {
	// read all events
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch et := event.(type) {
		case *sdl.QuitEvent:
			// TODO handle special stuff, for exampe quit
			log.Println("TODO: how do we apply this... ?")
		case *sdl.KeyDownEvent:
			c.sendKeyEvents(sdl.GetKeyName(et.Keysym.Sym), true)
		case *sdl.KeyUpEvent:
			c.sendKeyEvents(sdl.GetKeyName(et.Keysym.Sym), false)
		// these are events we know but currently don't use FIXME implement
		case *sdl.WindowEvent:
		case *sdl.TextInputEvent:
		case *sdl.TextEditingEvent:
		case *sdl.MouseButtonEvent:
		case *sdl.MouseMotionEvent:
		case *sdl.MouseWheelEvent:
		default:
			log.Printf("Unknown event: %T\n", et)
		}
	}
}

func (c *Controller) iterateEntities() {
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
	// run for all ordered, renderable entities
	for _, entity := range c.renderable.get() {
		// allow execute of actions
		if actionEntity, actionable := entity.(Actionable); actionable {
			actionEntity.Action()
		}
		// determine draw method. Texture is preferred over Render!
		eR, renderable := entity.(Renderable)
		eD, drawable := entity.(Drawable)
		if drawable || renderable {
			// drawable or renderable both are placeable
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
			// check if we need to draw / redraw entity
			text, cacheExists := c.textCache[entity.Identification()]
			if !cacheExists {
				log.Println(e.Identification(), "cache miss: (re)drawing texture")
			}
			if !cacheExists || e.Redraw() {
				// if previous existed, destroy
				if cacheExists {
					text.Destroy()
				}
				if drawable {
					// create texture to draw to
					text, err = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STATIC,
						e.Width(), e.Height())
					if err != nil {
						log.Println("renderer.CreateTexture error:", err)
						continue
					}
					// allow entity to draw to texture
					eD.Texture(text)
				} else if renderable {
					text, _ = renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_TARGET,
						e.Width(), e.Height())
					renderer.SetRenderTarget(text)
					eR.Render(c.wrapped)
					renderer.SetRenderTarget(nil)
				} else {
					// fail loudly if invalid entity
					log.Fatal("iterateEntities:", e.Identification(), "neither drawable nor renderable!")
				}
				// update cache
				c.textCache[entity.Identification()] = text
			}
			// TODO last nil is rotation center, use offset to calculate
			renderer.CopyEx(c.textCache[entity.Identification()], nil, worldBound, e.Rotation(), nil, 0)
		}
	} // for entities
	renderer.Present()
}

/*
sendKeyEvents is a helper function that sends the given key with the given state
to all known keyreceivers that have registered for it.
*/
func (c *Controller) sendKeyEvents(key string, state bool) {
	// check if any interest exists
	receivers, exist := c.keyreceivers[key]
	if !exist {
		return
	}
	// if yes send state to all channels
	for index := range receivers { // use index because we'll set a value directly
		// don't send multiple times
		if receivers[index].lastState == state {
			continue
		}
		receivers[index].lastState = state
		// execute if available
		if state {
			if receivers[index].onPress != nil {
				receivers[index].onPress()
			}
		} else {
			if receivers[index].onRelease != nil {
				receivers[index].onRelease()
			}
		}
	}
}
