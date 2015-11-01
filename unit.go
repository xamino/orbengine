package orbengine

/*
Unit is an abstraction for designating number values.
*/
type Unit struct {
	aspect float32
}

/*
CreateUnit creates a new unit for converting arbitrary measurement units to
engine units.
*/
func CreateUnit(aspect float32) *Unit {
	return &Unit{aspect: aspect}
}

/*
Convert applies the aspect to a unit value.
*/
func (u *Unit) convert(value int) int {
	return int(float32(value) * u.aspect)
}
