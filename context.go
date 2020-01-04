package gousu

import (
	"log"
)

// IContext defines the interface of Context used for dependency injection (DI)
type IContext interface {
	RegisterService(service IService)
	RegisterController(controller IController)
	GetService(name string) IService
	GetController(name string) IController
}

// Context is used for dependency injection (DI) from Runner to services and controllers
type Context struct {
	services    map[string]IService
	controllers map[string]IController
}

var _ (IContext) = (*Context)(nil)

// RegisterService registers a service by its name (returned by IService.Name())
//
// Causes a fatal failure if the name is empty or already in use
func (c *Context) RegisterService(service IService) {
	name := service.Name()

	if name == "" {
		log.Fatalf("Error registering service %v: empty name", service)

		return
	}

	if _, ok := c.services[name]; ok {
		log.Fatalf("Error registering service %v: name %s already in use", service, name)

		return
	}

	c.services[name] = service
}

// RegisterController registers a controller by its name (returned by IController.Name())
//
// Causes a fatal failure if the name is empty or already in use
func (c *Context) RegisterController(controller IController) {
	name := controller.Name()

	if name == "" {
		log.Fatalf("Error registering controller %v: empty name", controller)

		return
	}

	if _, ok := c.controllers[name]; ok {
		log.Fatalf("Error registering controller %v: name %s already in use", controller, name)

		return
	}

	c.controllers[name] = controller
}

// GetService returns a service by its name
//
// Causes a fatal failure if no service is registered for this name
func (c *Context) GetService(name string) IService {
	service, ok := c.services[name]

	if !ok {
		log.Fatalf("Error getting service %s: unknown service", name)

		return nil
	}

	return service
}

// GetController returns a controller by its name
//
// Causes a fatal failure if no controller is registered for this name
func (c *Context) GetController(name string) IController {
	controller, ok := c.controllers[name]

	if !ok {
		log.Fatalf("Error getting controller %s: unknown controller", name)

		return nil
	}

	return controller
}

// NewContext creates a new initialized instance of Context
func NewContext() *Context {
	return &Context{
		services:    map[string]IService{},
		controllers: map[string]IController{},
	}
}
