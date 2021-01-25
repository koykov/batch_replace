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
	// Just use batch_replace.P.Get() and batch_replace.P.Put().
	P Pool
)

// Get old or create new instance of the batch replacer.
func (p *Pool) Get(s []byte) *BatchReplace {
	v := p.p.Get()
	if v != nil {
		if r, ok := v.(*BatchReplace); ok {
			r.SetSrcBytes(s)
			return r
		}
	}
	return NewBatchReplace(s)
}

// Put batch replacer to the pool.
func (p *Pool) Put(r *BatchReplace) {
	r.Reset()
	p.p.Put(r)
}

// Get replacer from default pool.
func Acquire(s []byte) *BatchReplace {
	return P.Get(s)
}

// Get replacer from default pool and set string as a source.
func AcquireStr(s string) *BatchReplace {
	return P.Get(fastconv.S2B(s))
}

// Put replacer back to default pool.
func Release(x *BatchReplace) {
	P.Put(x)
}
