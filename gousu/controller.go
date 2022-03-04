package gousu

// IController is the base type for all controllers
type IController interface {
	Name() string

	Start() error
	Stop() error
	Health() error
}

// ControllerFactory defines the factory function for creating new instances of
// the controller
type ControllerFactory func(ctx IContext) IController
