package batch_replace

import (
	"bytes"
	"strconv"

	"github.com/koykov/fastconv"
)

const (
	// Int base edges.
	baseLo = 2
	baseHi = 36

	// Buffers length.
	intBufLen = 64 // 64 bit + 1 for sign for base 2
	fltBufLen = 24
)

// Replace queue.
type batchReplaceQueue struct {
	// Queue items.
	queue [][]byte
	// Current index.
	idx int
	// Max index.
	cap int
	// Accumulated items length.
	acc int
}

// Replacer.
type BatchReplace struct {
	// Replacer source.
	src []byte
	// Destination slice.
	dst []byte
	// Internal buffer.
	buf []byte
	// Queue of old parts.
	old batchReplaceQueue
	// Queue of replacements.
	new batchReplaceQueue
}

// Init new replacer.
func NewBatchReplace(s []byte) *BatchReplace {
	o := batchReplaceQueue{queue: make([][]byte, 0)}
	n := batchReplaceQueue{queue: make([][]byte, 0)}
	r := BatchReplace{
		src: s,
		old: o,
		new: n,
	}
	return &r
}

// Set the source.
//
// For use outside of pools.
func (r *BatchReplace) SetSrc(src []byte) *BatchReplace {
	r.src = append(r.src[:0], src...)
	return r
}

// Set the source using string.
func (r *BatchReplace) SetSrcStr(src string) *BatchReplace {
	r.src = append(r.src[:0], src...)
	return r
}

// Register new byte slice replacement.
func (r *BatchReplace) Replace(old []byte, new []byte) *BatchReplace {
	n := bytes.Count(r.src, old)
	if n == 0 {
		return r
	}
	r.old.add(old, n)
	r.new.add(new, n)
	return r
}

// Register new string replacement.
func (r *BatchReplace) ReplaceStr(old, new string) *BatchReplace {
	return r.Replace(fastconv.S2B(old), fastconv.S2B(new))
}

// Register int replacement.
func (r *BatchReplace) ReplaceInt(old []byte, new int64) *BatchReplace {
	return r.ReplaceIntBase(old, new, 10)
}

// Register int replacement with given base.
func (r *BatchReplace) ReplaceIntBase(old []byte, new int64, base int) *BatchReplace {
	n := bytes.Count(r.src, old)
	if n == 0 || base < baseLo || base > baseHi {
		return r
	}
	r.old.add(old, n)

	nb := r.new.next(intBufLen)
	nb = strconv.AppendInt(nb, new, base)
	r.new.set(nb, n)
	return r
}

// Register uint replacement.
func (r *BatchReplace) ReplaceUint(old []byte, new uint64) *BatchReplace {
	return r.ReplaceUintBase(old, new, 10)
}

// Register uint replacement with given base.
func (r *BatchReplace) ReplaceUintBase(old []byte, new uint64, base int) *BatchReplace {
	n := bytes.Count(r.src, old)
	if n == 0 || base < baseLo || base > baseHi {
		return r
	}
	r.old.add(old, n)

	nb := r.new.next(intBufLen)
	nb = strconv.AppendUint(nb, new, base)
	r.new.set(nb, n)
	return r
}

// Register float replacement.
func (r *BatchReplace) ReplaceFloat(old []byte, new float64) *BatchReplace {
	return r.ReplaceFloatTunable(old, new, 'f', -1, 64)
}

// Register float replacement with params.
func (r *BatchReplace) ReplaceFloatTunable(old []byte, new float64, fmt byte, prec, bitSize int) *BatchReplace {
	n := bytes.Count(r.src, old)
	if n == 0 {
		return r
	}
	r.old.add(old, n)

	nb := r.new.next(max(prec+4, fltBufLen))
	nb = strconv.AppendFloat(nb, new, fmt, prec, bitSize)
	r.new.set(nb, n)
	return r
}

// Perform the replaces.
func (r *BatchReplace) Commit() []byte {
	// Calculate final length.
	l := len(r.src) + r.new.acc - r.old.acc

	r.dst = append(r.dst[:0], r.src...)
	// Walk over queue and replace.
	for i := 0; i < len(r.old.queue); i++ {
		o := r.old.queue[i]
		n := r.new.queue[i]
		c := bytes.Count(r.src, o)
		r.buf = replaceTo(r.buf[:0], r.dst, o, n, c)
		r.dst = append(r.dst[:0], r.buf...)
	}

	return r.dst[:l]
}

// Perform the replaces and return copy of result.
//
// Made to avoid data sharing.
func (r *BatchReplace) CommitCopy() []byte {
	return append([]byte(nil), r.Commit()...)
}

// String version of Commit().
func (r *BatchReplace) CommitStr() string {
	return fastconv.B2S(r.Commit())
}

// String version of CommitCopy().
func (r *BatchReplace) CommitCopyStr() string {
	return fastconv.B2S(r.CommitCopy())
}

// Clear the replacer with keeping of allocated space to reuse.
func (r *BatchReplace) Reset() *BatchReplace {
	r.src = r.src[:0]
	for i := 0; i < len(r.old.queue); i++ {
		r.old.queue[i] = r.old.queue[i][:0]
		r.old.idx, r.old.acc = 0, 0
		r.new.queue[i] = r.new.queue[i][:0]
		r.new.idx, r.new.acc = 0, 0
	}
	return r
}

// Add new item to queue.
func (q *batchReplaceQueue) add(p []byte, n int) {
	if n == 0 {
		n = 1
	}
	q.acc += len(p) * n
	if q.idx < q.cap {
		q.queue[q.idx] = append(q.queue[q.idx][:0], p...)
		q.idx++
		return
	}
	q.queue = append(q.queue, append([]byte(nil), p...))
	q.idx++
	q.cap++
}

// Update last item in queue.
func (q *batchReplaceQueue) set(p []byte, n int) {
	q.queue[q.idx-1] = append(q.queue[q.idx-1], p...)
	q.acc += len(p) * n
}

// Shift to the next item.
func (q *batchReplaceQueue) next(max int) []byte {
	if q.idx < q.cap {
		b := q.queue[q.idx]
		q.idx++
		return b
	}
	b := make([]byte, 0, max)
	q.queue = append(q.queue, b)
	q.idx++
	q.cap++
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
