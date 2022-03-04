package gousu

// IService is the base type for all services
type IService interface {
	Name() string

	Start() error
	Stop() error
	Health() error
}

// ServiceFactory defines the factory function for creating new instances of
// the service
type ServiceFactory func(ctx IContext) IService
