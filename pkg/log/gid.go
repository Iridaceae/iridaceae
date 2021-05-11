package log

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"
)

// we introduce buffer as a sync.Pool to provide a safe way to access our goroutine.
// this should only be used during the logging context debugging purpose.
var buffer = sync.Pool{New: func() interface{} {
	b := make([]byte, 64)
	return &b
}}

// Goid returns a goroutine id of given stack trace.
// As the word of the wise Dave Cheney you might as well go to hell by using this.
// This is a hacky pure Go version of dave cheney's implementation.
func Goid() uint64 {
	p, _ := buffer.Get().(*[]byte)
	defer buffer.Put(p)
	b := *p
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
