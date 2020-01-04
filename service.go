package gousu

// IService is the base type for all services
type IService interface {
	Start() error
	Stop() error
	Health() error
}
