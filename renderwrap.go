package orbengine

import "github.com/veandco/go-sdl2/sdl"

type Renderer struct {
	renderer *sdl.Renderer
}

func (w *Renderer) SetDrawColor(r, g, b, a uint8) {
	w.renderer.SetDrawColor(r, g, b, a)
}

func (w *Renderer) FillRect(rect *sdl.Rect) {
	w.renderer.FillRect(rect)
}

func (w *Renderer) GetRendererOutputSize() (int, int, error) {
	return w.renderer.GetRendererOutputSize()
}
