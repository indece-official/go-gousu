package broadcaster

type Bool = Generic[bool]

var _ Base = (*Bool)(nil)

func NewBool(initialValue bool) *Bool {
	return NewGeneric(initialValue)
}
