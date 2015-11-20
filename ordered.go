package orbengine

import "log"

/*
ordered is a helper struct that sorts Renderable and Drawable on append. Allows
the controller to get the sorted list for each iteration.
*/
type ordered struct {
	list []Placeable
}

func makeOrdered() *ordered {
	return &ordered{
		list: make([]Placeable, 32)} // TODO is 32 a good default?
}

func (o *ordered) append(entity Placeable) error {
	switch entity.(type) {
	case Drawable:
	case Renderable:
	default:
		return ErrMissingComponents
	}
	log.Println("TODO: sort")
	o.list = append(o.list, entity)
	return nil
}

func (o *ordered) get() []Placeable {
	return o.list
}
