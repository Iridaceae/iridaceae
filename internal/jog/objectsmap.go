package jog

import "sync"

// ObjectsMap wraps around map[string]interface to provide thread-safe access endpoints.
// This is a implementation of *sync.Map.
type ObjectsMap struct {
	mutex sync.RWMutex
	inner map[string]interface{}
}

// newObjectsMap initializes a new ObjectsMap instance.
func newObjectsMap() *ObjectsMap {
	return &ObjectsMap{inner: make(map[string]interface{})}
}

// Get a value from a map.
// If a value was found, the value and true returns or otherwise.
func (o *ObjectsMap) Get(key string) (interface{}, bool) {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	v, ok := o.inner[key]
	return v, ok
}

// GetValue wraps Get and returns a value, nil otherwise.
func (o *ObjectsMap) GetValue(key string) interface{} {
	v, ok := o.Get(key)
	if !ok {
		return nil
	}
	return v
}

// Set a value from a map by key.
func (o *ObjectsMap) Set(key string, val interface{}) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.inner[key] = val
}

// Delete removes a key from key-value pair from the map.
func (o *ObjectsMap) Delete(key string) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	delete(o.inner, key)
}

// Count returns number of elements inside map.
func (o *ObjectsMap) Count() int {
	return len(o.inner)
}
