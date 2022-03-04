package gousu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsString(t *testing.T) {
	arr := []string{}
	assert.False(t, ContainsString(arr, "test"))

	arr = []string{"test0"}
	assert.True(t, ContainsString(arr, "test0"))

	arr = []string{"test0", "test1", "test2"}
	assert.True(t, ContainsString(arr, "test2"))

	arr = []string{"test0", "test1", "test2"}
	assert.False(t, ContainsString(arr, "test3"))
}
