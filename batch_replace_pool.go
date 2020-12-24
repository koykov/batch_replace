package batch_replace

import (
	"sync"

	"github.com/koykov/fastconv"
)

// Pool to store batch replacers.
type BatchReplacePool struct {
	p sync.Pool
}

// Default instance of the pool.
// Just use batch_replace.BatchPool.Get() and batch_replace.BatchPool.Put().
var BatchPool BatchReplacePool

// Get old or create new instance of the batch replacer.
func (p *BatchReplacePool) Get(s []byte) *BatchReplace {
	v := p.p.Get()
	if v != nil {
		if r, ok := v.(*BatchReplace); ok {
			r.src = append(r.src, s...)
			return r
		}
	}
	return NewBatchReplace(s)
}

// Put batch replacer to the pool.
func (p *BatchReplacePool) Put(r *BatchReplace) {
	r.Reset()
	p.p.Put(r)
}

// Get replacer from default pool.
func (p *BatchReplacePool) Acquire(s []byte) *BatchReplace {
	return BatchPool.Get(s)
}

// Get replacer from default pool and set string as a source.
func (p *BatchReplacePool) SAcquire(s string) *BatchReplace {
	return BatchPool.Get(fastconv.S2B(s))
}

// Put replacer back to default pool.
func (p *BatchReplacePool) Release(x *BatchReplace) {
	BatchPool.Put(x)
}
