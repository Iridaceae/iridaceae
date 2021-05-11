package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertGetUniquePanic(t *testing.T) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%+v doesn't panic when it should be", Mapper().getUnique(""))
		}
	}()
	Mapper().getUnique("")
}

func TestInitGlobalStorage(t *testing.T) {
	s := Mapper()
	assert.Equal(t, &_globalMapper, &s)
}

// TODO: greedy reset global storage
// Could implement a DI container as a sidecar.

func TestStorage_Set(t *testing.T) {
	ResetGlobalStorage()
	Mapper().Set("test", 1)
	assert.Equal(t, 1, Mapper().Count())
}

func TestStorage_SetMap(t *testing.T) {
	ResetGlobalStorage()
	test := map[string]interface{}{
		"k1": 1,
		"k2": "2",
		"k3": uint64(2),
	}
	Mapper().SetMap(test)
	assert.Equal(t, 3, Mapper().Count())
}

func TestStorage_SetAbsent(t *testing.T) {
	ResetGlobalStorage()
	Mapper().SetAbsent("test", 1)
	ok := Mapper().SetAbsent("test", 1)
	assert.False(t, ok)
}

func TestStorage_Get(t *testing.T) {
	k, ok := Mapper().Get("name")
	assert.Nil(t, k)
	assert.False(t, ok)
}

func TestStorage_Count(t *testing.T) {
	ResetGlobalStorage()
	assert.Len(t, Mapper().items, 0)
}

func TestStorage_GetString(t *testing.T) {
	ResetGlobalStorage()
	Mapper().Set("ver", "v1")
	assert.Equal(t, "v1", Mapper().GetString("ver"))
	assert.Equal(t, "", Mapper().GetString("this doesn't exists"))
}

func TestStorage_Has(t *testing.T) {
	assert.False(t, Mapper().Has("test"))
}

func TestStorage_IsEmpty(t *testing.T) {
	ResetGlobalStorage()
	assert.True(t, Mapper().IsEmpty())
	assertGetUniquePanic(t)
}

func TestStorage_Keys(t *testing.T) {
	tests := []struct {
		name string
		key  string
		id   string
	}{
		{"k1 id", "k1", Mapper().getUnique("k1")},
		{"k2 id", "k2", Mapper().getUnique("k2")},
		{"k3 id", "k3", Mapper().getUnique("k3")},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResetGlobalStorage()
			Mapper().Set(tt.key, fmt.Sprintf("value-%d", i))
			assert.Equal(t, Mapper().getUnique(tt.key), Mapper().Keys()[0])
		})
	}
}

func TestStorage_Remove(t *testing.T) {
	ResetGlobalStorage()
	Mapper().Set("test", 1)
	Mapper().Remove("test")
	assert.Equal(t, 0, Mapper().Count())
}
