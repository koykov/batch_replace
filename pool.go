package batch_replace

import (
	"sync"

	"github.com/koykov/fastconv"
)

// Pool to store batch replacers.
type Pool struct {
	p sync.Pool
}

var (
	// Default instance of the pool.
	// Just use batch_replace.Acquire() and batch_replace.Release().
	p Pool
	// Suppress go vet warnings.
	_ = Acquire
)

// Get old or create new instance of the batch replacer.
func (p *Pool) Get(src []byte) *BatchReplace {
	v := p.p.Get()
	if v != nil {
		if r, ok := v.(*BatchReplace); ok {
			if len(src) > 0 {
				r.SetSrcBytes(src)
			}
			return r
		}
	}
	return NewBatchReplace(src)
}

// Put batch replacer to the pool.
func (p *Pool) Put(r *BatchReplace) {
	r.Reset()
	p.p.Put(r)
}

// Get replacer from default pool.
//
// Please note, this method doesn't provide source for replacer and you need to specify it manually by calling
// SetSrcBytes() and SetSrcStr() methods.
func Acquire() *BatchReplace {
	return p.Get(nil)
}

// Get replacer from default pool and set byte array as a source.
func AcquireWithBytesSrc(src []byte) *BatchReplace {
	return p.Get(src)
}

// Get replacer from default pool and set string as a source.
func AcquireWithStrSrc(src string) *BatchReplace {
	return p.Get(fastconv.S2B(src))
}

// Put replacer back to default pool.
func Release(x *BatchReplace) {
	p.Put(x)
}
