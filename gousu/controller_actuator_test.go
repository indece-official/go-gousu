package gousu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewActuatorController(t *testing.T) {
	ctx := NewContext()
	ctx.RegisterService(NewMockService())
	controller := NewActuatorController(ctx)

	assert.NotNil(t, controller)
}
