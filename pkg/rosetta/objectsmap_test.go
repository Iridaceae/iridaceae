package rosetta

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

const count = 1000

func init() {
	newTestObjectsMap()
}

func newTestObjectsMap() {
	TestObjectsMap = newObjectsMap()
	TestObjectsMap.mutex = sync.RWMutex{}
}

func loop(f func(i int, key, value string)) {
	wg := new(sync.WaitGroup)
	wg.Add(count)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)

		go func(i int, key, value string) {
			defer wg.Done()
			f(i, key, value)
		}(i, key, value)
	}
	wg.Wait()
}

func setup() {
	loop(func(i int, key, value string) {
		TestObjectsMap.Set(key, value)
	})
}

func TestObjectsMap_Set(t *testing.T) {
	t.Run("set one value", func(t *testing.T) {
		TestObjectsMap.Set("rosetta_testObjectsMap", 10)
		assert.Equal(t, 10, TestObjectsMap.inner["rosetta_testObjectsMap"])
		assert.Equal(t, 1, len(TestObjectsMap.inner))
	})

	t.Run("concurrent 1000 times", func(t *testing.T) {
		newTestObjectsMap()
		setup()
		assert.Equal(t, count, len(TestObjectsMap.inner))
	})
}

func TestObjectsMap_Count(t *testing.T) {
	setup()
	if TestObjectsMap.Count() != count {
		t.Fatalf("expected %d elements in `TestObjectsMap`, got %d instead", count, TestObjectsMap.Count())
	}
}

func TestObjectsMap_Get(t *testing.T) {
	setup()
	loop(func(i int, key, value string) {
		got, ok := TestObjectsMap.Get(key)
		assert.Equal(t, value, got)
		assert.True(t, ok)
	})
}

func TestObjectsMap_GetValue(t *testing.T) {
	setup()
	loop(func(i int, key, value string) {
		got := TestObjectsMap.GetValue(key)
		assert.Equal(t, value, got)
	})
}

func TestObjectsMap_Delete(t *testing.T) {
	setup()
	loop(func(i int, key, value string) {
		TestObjectsMap.Delete(key)
	})
	assert.Equal(t, 0, TestObjectsMap.Count())
}
