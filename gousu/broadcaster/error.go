package broadcaster

type Error = Generic[error]

var _ Base = (*Error)(nil)

func NewError(initialValue error) *Error {
	return NewGeneric(initialValue)
}
