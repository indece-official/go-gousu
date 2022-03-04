package goususmtp

import (
	"testing"

	"github.com/indece-official/go-gousu/gousu"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	ctx := gousu.NewContext()

	service := NewService(ctx)

	assert.NotNil(t, service)
	assert.IsType(t, &Service{}, service)
}
