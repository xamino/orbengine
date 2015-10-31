package orbengine

import "github.com/veandco/go-sdl2/sdl"

/*
Drawable is the interface for drawing entities.

TODO: wrap sdl.Renderer in own renderer from controller to control access to methods?
*/
type Drawable interface {
	Draw(renderer *sdl.Renderer)
}

/*
Boundable is the interface for entity bounds.
*/
type Boundable interface {
	Bounds() *sdl.Rect
}

/*
Positionable is the interface for entity positions.
*/
type Positionable interface {
	Position() (int32, int32)
}

/*
Actionable is the interface for entity actions.
TODO: this is where we can animate entities. Check how this will need to be
implemented.
*/
type Actionable interface {
	Action()
}
