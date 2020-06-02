package batch_replace

import (
	"sync"
)

// Pool to store batch string replacers.
type BatchReplaceStrPool struct {
	p sync.Pool
}

// Default instance of the pool.
// Just use cbytealg.BatchStrPool.Get() and cbytealg.BatchStrPool.Put().
var BatchStrPool BatchReplaceStrPool

// Get old or create new instance of the batch string replacer.
func (p *BatchReplaceStrPool) Get(s string) *BatchReplaceStr {
	v := p.p.Get()
	if v != nil {
		if r, ok := v.(*BatchReplaceStr); ok {
			r.r.src = append(r.r.src, s...)
			return r
		}
	}
	return NewBatchReplaceStr(s)
}

// Put batch string replacer to the pool.
func (p *BatchReplaceStrPool) Put(r *BatchReplaceStr) {
	r.Reset()
	p.p.Put(r)
}
