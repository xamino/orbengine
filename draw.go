package orbengine

import "github.com/veandco/go-sdl2/sdl"

/*
Draw is the interface for drawable entities.
*/
type Draw interface {
	Draw(renderer *sdl.Renderer)
}
