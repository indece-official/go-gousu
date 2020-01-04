package gousu

// IController is the base type for all controllers
type IController interface {
	Start() error
	Stop() error
	Health() error
}
