package batch_replace

import (
	"sync"

	"github.com/koykov/byteconv"
	"github.com/koykov/byteseq"
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
				r.SetSource(src)
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

// Acquire gets replacer from default pool.
//
// Please note, this method doesn't provide source for replacer and you need to specify it manually by calling
// SetSource() method.
func Acquire() *BatchReplace {
	return p.Get(nil)
}

// AcquireWithSource gets replace from default pool and sets source for replacements.
func AcquireWithSource[T byteseq.Byteseq](x T) *BatchReplace {
	var src []byte
	if b, ok := byteseq.ToBytes(x); ok {
		src = b
	}
	if s, ok := byteseq.ToString(x); ok {
		src = byteconv.S2B(s)
	}
	br := p.Get(src)
	return br
}

// AcquireWithBytesSrc gets replacer from default pool and set byte array as a source.
// Deprecated: use AcquireWithSource() instead.
func AcquireWithBytesSrc(src []byte) *BatchReplace {
	return p.Get(src)
}

// AcquireWithStrSrc gets replacer from default pool and set string as a source.
// Deprecated: use AcquireWithSource() instead.
func AcquireWithStrSrc(src string) *BatchReplace {
	return p.Get(byteconv.S2B(src))
}

// Release puts replacer back to default pool.
func Release(x *BatchReplace) {
	p.Put(x)
}
