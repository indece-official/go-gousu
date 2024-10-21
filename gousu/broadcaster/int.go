package broadcaster

type Int = Generic[int]

var _ Base = (*Int)(nil)

func NewInt(initialValue int) *Int {
	return NewGeneric(initialValue)
}
