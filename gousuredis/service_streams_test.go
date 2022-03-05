package gousuredis

import (
	"testing"
	"time"

	"github.com/indece-official/go-gousu/v2/gousu"
	"github.com/stretchr/testify/assert"
)

func TestXAdd(t *testing.T) {
	ctx := gousu.NewContext()

	service := NewService(ctx).(IService)
	assert.NoError(t, service.Start())

	id, err := service.XAdd(
		"teststream01",
		map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestXGroupCreate(t *testing.T) {
	ctx := gousu.NewContext()

	service := NewService(ctx).(IService)
	assert.NoError(t, service.Start())

	err := service.XGroupCreate(
		"testgroup01",
		"teststream01",
		XGroupCreateOffsetFirst,
		true,
		true,
	)

	assert.NoError(t, err)
}

func TestXReadGroup(t *testing.T) {
	ctx := gousu.NewContext()

	service := NewService(ctx).(IService)
	assert.NoError(t, service.Start())

	go func() {
		time.Sleep(100 * time.Millisecond)

		_, err := service.XAdd(
			"teststream01",
			map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		)
		assert.NoError(t, err)
	}()

	evt, err := service.XReadGroup(
		"testgroup01",
		"testconsumer01",
		"teststream01",
		1*time.Second,
		XReadGroupIDStreamNew,
	)

	assert.NoError(t, err)
	assert.NotEmpty(t, evt)
	assert.NotEmpty(t, evt.ID)
	assert.NotEmpty(t, evt.Data)
	assert.Equal(t, "teststream01", evt.Key)

	count, err := service.XAck(
		"testgroup01",
		"teststream01",
		evt.ID,
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
