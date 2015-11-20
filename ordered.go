package orbengine

import "sort"

/*
ordered is a helper struct that sorts Renderable and Drawable on append. Allows
the controller to get the sorted list for each iteration.
*/
type ordered struct {
	list []Placeable
}

func makeOrdered() *ordered {
	return &ordered{
		list: make([]Placeable, 0)}
}

func (o *ordered) append(entity Placeable) error {
	switch entity.(type) {
	case Drawable:
	case Renderable:
	default:
		return ErrMissingComponents
	}
	// append
	o.list = append(o.list, entity)
	// sort
	sortable := byLevel(o.list)
	sort.Sort(sortable)
	o.list = []Placeable(sortable)
	// all ok
	return nil
}

func (o *ordered) get() []Placeable {
	return o.list
}

type byLevel []Placeable

func (bl byLevel) Len() int {
	return len(bl)
}

func (bl byLevel) Less(i, j int) bool {
	// if the object is nil, sort that non nil objects are worth less.
	if bl[i] == nil {
		return true
	}
	if bl[j] == nil {
		return false
	}
	// otherwise sort by layer
	return bl[i].Layer() < bl[j].Layer()
}

func (bl byLevel) Swap(i, j int) {
	bl[i], bl[j] = bl[j], bl[i]
}
