package gousu

import (
	"fmt"
	"net/http"

	"github.com/namsral/flag"
)

var (
	actuatorHost = flag.String("actuator_host", "0.0.0.0", "")
	actuatorPort = flag.Int("actuator_port", 9000, "")
)

// ActuatorController is a controller running in a separate thread providing an health endpoint
type ActuatorController struct {
	services []IService
	log      *Log
	error    error
}

// ActuatorController implement IController
var _ IController = (*ActuatorController)(nil)

// Start starts a HTTP-server for health-checks
func (c *ActuatorController) Start() error {
	c.error = nil

	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			for i := range c.services {
				service := c.services[i]

				err := service.Health()
				if err != nil {
					w.WriteHeader(500)
					fmt.Fprintf(w, "Not healthy: %s", err)

					return
				}
			}

			fmt.Fprintf(w, "OK")
		})

		err := http.ListenAndServe(fmt.Sprintf("%s:%d", *actuatorHost, *actuatorPort), nil)
		if err != nil {
			c.error = c.log.ErrorfX("Can't start actuator server: %s", err)
		}
	}()

	c.log.Infof("Actuator server listening on %s:%d", *actuatorHost, *actuatorPort)

	return nil
}

// Health checks if the http server properly started
func (c *ActuatorController) Health() error {
	return c.error
}

// Stop currently does nothing
func (c *ActuatorController) Stop() error {
	return nil
}

// NewActuatorController creates a new initilized instance of ActuatorController
func NewActuatorController(services []IService) *ActuatorController {
	return &ActuatorController{
		log:      GetLogger("controller.actuator"),
		services: services,
	}
}
