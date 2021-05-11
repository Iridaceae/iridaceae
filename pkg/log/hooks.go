package log

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/rs/zerolog"
)

var (
	_globalMapper = InitGlobalStorage()
	mapper        *Storage
)

type MapperHook struct{}

func (m MapperHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if len(_fields) == 0 {
		return
	}
	fields := make(map[string]interface{})
	for _, f := range _fields {
		fields[f] = Mapper().GetString(f)
	}
	e.Fields(fields)
}

// Storage defines our key, value mapping for our context.
// Storage is inherently safe for concurrency and goroutines.
type Storage struct {
	items map[string]interface{}
	sync.RWMutex
}

// InitGlobalStorage initializes a Storage objects with a key,value map.
func InitGlobalStorage() *Storage {
	if mapper == nil {
		mapper = &Storage{items: make(map[string]interface{})}
	}
	return mapper
}

// ResetGlobalStorage sets an empty default map.
func ResetGlobalStorage() {
	mapper.RLock()
	defer mapper.RUnlock()
	_globalMapper.items = make(map[string]interface{})
}

// Mapper initialize a global mapper for our logger.
func Mapper() *Storage {
	mapper.RLock()
	defer mapper.RUnlock()
	s := _globalMapper
	return s
}

// Set adds key-value pair to Storage map.
func (s *Storage) Set(key string, value interface{}) {
	unique := s.getUnique(key)
	s.Lock()
	s.items[unique] = value
	s.Unlock()
}

// SetMap adds given map to mapper's items.
func (s *Storage) SetMap(mp map[string]interface{}) {
	for k, v := range mp {
		unique := s.getUnique(k)
		s.Lock()
		s.items[unique] = v
		s.Unlock()
	}
}

// SetAbsent returns false if value is already exists in given map, true otherwise and insert into
// give map. This is safe concurrently.
func (s *Storage) SetAbsent(key string, value interface{}) bool {
	unique := s.getUnique(key)
	s.Lock()
	_, ok := s.items[unique]
	if !ok {
		s.items[unique] = value
	}
	s.Unlock()
	return !ok
}

// Get returns a given value in our mapper context key,value pair.
func (s *Storage) Get(key string) (interface{}, bool) {
	unique := s.getUnique(key)
	s.RLock()
	defer s.RUnlock()
	v, ok := s.items[unique]
	return v, ok
}

// GetString returns a string representation of value.
func (s *Storage) GetString(key string) string {
	unique := s.getUnique(key)
	s.RLock()
	defer s.RUnlock()
	v, ok := s.items[unique]
	if !ok {
		v = ""
	}
	return fmt.Sprintf("%v", v)
}

// Count returns # of k,v in our mapper mapping.
func (s *Storage) Count() int {
	s.Lock()
	defer s.Unlock()
	count := len(s.items)
	return count
}

// Has check if key is already exists in map.
func (s *Storage) Has(key string) bool {
	unique := s.getUnique(key)
	s.RLock()
	defer s.RUnlock()
	_, ok := s.items[unique]
	return ok
}

// IsEmpty checks if map is empty.
func (s *Storage) IsEmpty() bool {
	s.RLock()
	defer s.RUnlock()
	return len(s.items) == 0
}

// Remove a given key from our map.
func (s *Storage) Remove(key string) {
	unique := s.getUnique(key)
	s.Lock()
	defer s.Unlock()
	delete(s.items, unique)
}

// Keys returns a list of string.
func (s *Storage) Keys() []string {
	count := s.Count()
	ch := make(chan string, count)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(a *Storage) {
			a.RLock()
			for key := range a.items {
				ch <- key
			}
			a.RUnlock()
			wg.Done()
		}(s)
		wg.Wait()
		close(ch)
	}()

	keys := make([]string, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

func (s *Storage) getUnique(key string) string {
	if key == "" {
		panic("key cannot be empty")
	}
	return key + "-" + strconv.FormatUint(Goid(), 10)
}
