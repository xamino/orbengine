package orbengine

import "github.com/veandco/go-sdl2/sdl"

/*
Drawable is the interface an entity must fulfill to be drawable.
*/
type Drawable interface {
	Placeable
	Texture(texture *sdl.Texture)
}

/*
Renderable is an alternative interface for Drawable for manually drawing entities.
*/
type Renderable interface {
	Placeable
	Render(renderer *Renderer)
}

/*
Placeable is the interface used to determine where to place an object.
*/
type Placeable interface {
	Position() *sdl.Point
	Offset() *sdl.Point // Offset is center of orientation offset from position. FIXME use and apply
	Scale() int
	Width() int
	Height() int
	Rotation() float64
	Redraw() bool // TODO this should be somewhere else, but for now...
}

/*
Actionable is the interface for entity actions.
TODO: this is where we can animate entities. Check how this will need to be
implemented.
*/
type Actionable interface {
	Action()
}
