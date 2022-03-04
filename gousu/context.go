package gousu

// IContext defines the interface of Context used for dependency injection (DI)
type IContext interface {
	RegisterService(service IService)
	RegisterController(controller IController)
	RegisterUIController(uiController IUIController)
	GetService(name string) IService
	GetServices() []IService
	GetController(name string) IController
	GetControllers() []IController
	GetUIController() IUIController
}

// Context is used for dependency injection (DI) from Runner to services and controllers
type Context struct {
	services     map[string]IService
	controllers  map[string]IController
	uiController IUIController
}

var _ (IContext) = (*Context)(nil)

// RegisterService registers a service by its name (returned by IService.Name())
//
// Causes a fatal failure if the name is empty or already in use
func (c *Context) RegisterService(service IService) {
	name := service.Name()

	if name == "" {
		logFatalf("Error registering service %v: empty name", service)

		return
	}

	if _, ok := c.services[name]; ok {
		logFatalf("Error registering service %v: name %s already in use", service, name)

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
		logFatalf("Error registering controller %v: empty name", controller)

		return
	}

	if _, ok := c.controllers[name]; ok {
		logFatalf("Error registering controller %v: name %s already in use", controller, name)

		return
	}

	c.controllers[name] = controller
}

// RegisterUIController registers a UI-Controller
//
// Causes a fatal failure if an UI-Controller was already registered
func (c *Context) RegisterUIController(uiController IUIController) {
	name := uiController.Name()

	if name == "" {
		logFatalf("Error registering UI-Controller %v: empty name", uiController)

		return
	}

	if c.uiController != nil {
		logFatalf("Error registering an UI-Controller %v", c.uiController.Name())

		return
	}

	c.uiController = uiController
}

// GetService returns a service by its name
//
// Causes a fatal failure if no service is registered for this name
func (c *Context) GetService(name string) IService {
	service, ok := c.services[name]

	if !ok {
		logFatalf("Error getting service %s: unknown service", name)

		return nil
	}

	return service
}

// GetServices returns a list of all registered services
func (c *Context) GetServices() []IService {
	services := make([]IService, len(c.services))

	i := 0
	for _, service := range c.services {
		services[i] = service
		i++
	}

	return services
}

// GetController returns a controller by its name
//
// Causes a fatal failure if no controller is registered for this name
func (c *Context) GetController(name string) IController {
	controller, ok := c.controllers[name]

	if !ok {
		logFatalf("Error getting controller %s: unknown controller", name)

		return nil
	}

	return controller
}

// GetControllers returns a list of all registered controllers
func (c *Context) GetControllers() []IController {
	controllers := make([]IController, len(c.controllers))

	i := 0
	for _, controller := range c.controllers {
		controllers[i] = controller
		i++
	}

	return controllers
}

// GetUIController returns the registered UI-Controller or nil
func (c *Context) GetUIController() IUIController {
	return c.uiController
}

// NewContext creates a new initialized instance of Context
func NewContext() *Context {
	return &Context{
		services:     map[string]IService{},
		controllers:  map[string]IController{},
		uiController: nil,
	}
}
