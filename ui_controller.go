package gousu

import "os"

// IUIController is a special controller that can be used for an ui and is executed on the main thread
type IUIController interface {
	IController

	Run(chan os.Signal) error
}

// UIControllerFactory defines the factory function for creating new instances of
// the ui controller
type UIControllerFactory func(ctx IContext) IUIController
