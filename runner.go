package gousu

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/namsral/flag"
)

// IRunner defines the interface of the core Runner
type IRunner interface {
	CreateService(serviceFactory ServiceFactory)
	CreateController(controllerFactory ControllerFactory)
	CreateUIController(uiControllerFactory UIControllerFactory)
	AwaitReady()
	Run()
	Kill()
}

// Runner is the core struct responsible for dependency injection and starting /
// stopping services & controllers
type Runner struct {
	ctx              IContext
	sigReady         chan bool
	sigTerm          chan os.Signal
	log              *Log
	servicesOrder    []string
	controllersOrder []string
	projectName      string
	version          string
}

// CreateService creates a new instance of a service using its factory function and
// registers it in the context
func (r *Runner) CreateService(serviceFactory ServiceFactory) {
	service := serviceFactory(r.ctx)

	r.ctx.RegisterService(service)
	r.servicesOrder = append(r.servicesOrder, service.Name())
}

// CreateController creates a new instance of a controller using its factory function
// and registers it in the context
func (r *Runner) CreateController(controllerFactory ControllerFactory) {
	controller := controllerFactory(r.ctx)

	r.ctx.RegisterController(controller)
	r.controllersOrder = append(r.controllersOrder, controller.Name())
}

// CreateUIController creates a new instance of an UI-Controller using its factory function
// and registers it in the context
func (r *Runner) CreateUIController(uiControllerFactory UIControllerFactory) {
	uiController := uiControllerFactory(r.ctx)

	r.ctx.RegisterUIController(uiController)
}

// Run is the blocking core function starting all services & controllers, waiting
// for a SIGINT or SIGTERM signal an the stopping all
func (r *Runner) Run() {
	r.log.Infof("Starting ...")

	for _, name := range r.servicesOrder {
		r.log.Infof("Starting service '%s' ...", name)

		CheckError(r.ctx.GetService(name).Start())

		r.log.Infof("Service '%s' started", name)
	}

	for _, name := range r.controllersOrder {
		r.log.Infof("Starting controller '%s' ...", name)

		CheckError(r.ctx.GetController(name).Start())

		r.log.Infof("Controller '%s' started", name)
	}

	uiController := r.ctx.GetUIController()

	if uiController != nil {
		r.log.Infof("Starting UI-Controller '%s' ...", uiController.Name())

		CheckError(uiController.Start())

		r.log.Infof("UI-Controller '%s' started", uiController.Name())
	}

	r.sigReady <- true

	if uiController != nil {
		uiController.Run(r.sigTerm)
	} else {
		<-r.sigTerm
	}

	r.log.Infof("Stopping ...")

	if uiController != nil {
		r.log.Infof("Stopping UI-Controller '%s' ...", uiController.Name())

		CheckError(uiController.Stop())

		r.log.Infof("UI-Controller '%s' stopped", uiController.Name())
	}

	for i := len(r.controllersOrder) - 1; i >= 0; i-- {
		name := r.controllersOrder[i]

		r.log.Infof("Stopping controller '%s' ...", name)

		r.ctx.GetController(name).Stop()

		r.log.Infof("Controller '%s' stopped", name)
	}

	for i := len(r.servicesOrder) - 1; i >= 0; i-- {
		name := r.servicesOrder[i]

		r.log.Infof("Stopping service '%s' ...", name)

		r.ctx.GetService(name).Stop()

		r.log.Infof("Service '%s' stopped", name)
	}
}

// AwaitReady is a blocking function waiting for the Runner to have started all
// services and controllers
func (r *Runner) AwaitReady() {
	<-r.sigReady
}

// Kill sends a SIGINt signal to the Runner causing it to stop
func (r *Runner) Kill() {
	r.sigTerm <- syscall.SIGINT
}

// NewRunner creates a new initialized instance of Runner, also initializing
// the config flags and the logger
func NewRunner(projectName string, version string) IRunner {
	sigTerm := make(chan os.Signal, 1)

	if !flag.Parsed() {
		flag.String(flag.DefaultConfigFlagname, "", "Path to config file")
		flag.Parse()
	}

	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM)

	InitLogger(projectName)

	log := GetLogger("main")

	log.Infof("%s %s", projectName, version)

	return &Runner{
		ctx:              NewContext(),
		log:              log,
		sigReady:         make(chan bool, 1),
		sigTerm:          sigTerm,
		servicesOrder:    []string{},
		controllersOrder: []string{},
		projectName:      projectName,
		version:          version,
	}
}
