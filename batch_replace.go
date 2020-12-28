package batch_replace

import (
	"bytes"
	"strconv"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

const (
	// Int base edges.
	baseLo = 2
	baseHi = 36
)

// Replace queue.
type queue struct {
	// Queue items.
	queue []byteptrn
	// Current index.
	idx int
	// Max index.
	cap int
	// Accumulated items length.
	acc int
}

// Byteptr struct with number of replacements.
type byteptrn struct {
	p byteptr
	n int
}

// Replacer.
type BatchReplace struct {
	// Common buffer to store all bytes data.
	// Data is allocating the following:
	// [ Source string | Old substr0 | New substr0 | Old substr1 | New substr1 | ... | Old substrN | New substrN | Destination string | Replacement buffer ]
	// That way need to reduce amount of pointers.
	buf []byte
	// Offset of used bytes in the buffer. Actually is a length of the buffer.
	off int
	// Source byte pointer.
	src byteptr
	// Destination byte pointer.
	dst byteptr
	// Queue of old parts.
	old queue
	// Queue of replacements.
	new queue
}

// Init new replacer.
func NewBatchReplace(s []byte) *BatchReplace {
	o := queue{queue: make([]byteptrn, 0)}
	n := queue{queue: make([]byteptrn, 0)}
	r := BatchReplace{
		old: o,
		new: n,
	}
	r.SetSrc(s)
	return &r
}

// Set the source.
//
// For use outside of pools.
func (r *BatchReplace) SetSrc(src []byte) *BatchReplace {
	r.buf = append(r.buf[:0], src...)
	r.src.set(r.off, len(src))
	r.off = len(src)
	return r
}

// Set the source using string.
func (r *BatchReplace) SetSrcStr(src string) *BatchReplace {
	return r.SetSrc(fastconv.S2B(src))
}

// Register new bytes replacement.
func (r *BatchReplace) Replace(old []byte, new []byte) *BatchReplace {
	n := bytes.Count(r.indirect(r.src), old)
	if n == 0 {
		return r
	}
	r.old.add(r.alloc(old), n)
	r.new.add(r.alloc(new), n)
	return r
}

// Register new string replacement.
// todo remove it due to SReplace() method.
func (r *BatchReplace) ReplaceStr(old, new string) *BatchReplace {
	return r.Replace(fastconv.S2B(old), fastconv.S2B(new))
}

// Register new string replacement.
func (r *BatchReplace) SReplace(old, new string) *BatchReplace {
	return r.Replace(fastconv.S2B(old), fastconv.S2B(new))
}

// Register bytes-int replacement.
func (r *BatchReplace) ReplaceInt(old []byte, new int64) *BatchReplace {
	return r.ReplaceIntBase(old, new, 10)
}

// Register string-int replacement.
func (r *BatchReplace) SReplaceInt(old string, new int64) *BatchReplace {
	return r.SReplaceIntBase(old, new, 10)
}

// Register bytes-int replacement with given base.
func (r *BatchReplace) ReplaceIntBase(old []byte, new int64, base int) *BatchReplace {
	n := bytes.Count(r.indirect(r.src), old)
	if n == 0 || base < baseLo || base > baseHi {
		return r
	}
	r.old.add(r.alloc(old), n)

	c := r.off
	r.buf = strconv.AppendInt(r.buf, new, base)
	r.off = len(r.buf)
	np := byteptr{
		o: c,
		l: r.off - c,
	}
	r.new.add(np, n)
	return r
}

// Register string-int replacement with given base.
func (r *BatchReplace) SReplaceIntBase(old string, new int64, base int) *BatchReplace {
	return r.ReplaceIntBase(fastconv.S2B(old), new, base)
}

// Register bytes-uint replacement.
func (r *BatchReplace) ReplaceUint(old []byte, new uint64) *BatchReplace {
	return r.ReplaceUintBase(old, new, 10)
}

// Register string-uint replacement.
func (r *BatchReplace) SReplaceUint(old string, new uint64) *BatchReplace {
	return r.SReplaceUintBase(old, new, 10)
}

// Register bytes-uint replacement with given base.
func (r *BatchReplace) ReplaceUintBase(old []byte, new uint64, base int) *BatchReplace {
	n := bytes.Count(r.indirect(r.src), old)
	if n == 0 || base < baseLo || base > baseHi {
		return r
	}
	r.old.add(r.alloc(old), n)

	c := r.off
	r.buf = strconv.AppendUint(r.buf, new, base)
	r.off = len(r.buf)
	np := byteptr{
		o: c,
		l: r.off - c,
	}
	r.new.add(np, n)
	return r
}

// Register string-uint replacement with given base.
func (r *BatchReplace) SReplaceUintBase(old string, new uint64, base int) *BatchReplace {
	return r.ReplaceUintBase(fastconv.S2B(old), new, base)
}

// Register bytes-float replacement.
func (r *BatchReplace) ReplaceFloat(old []byte, new float64) *BatchReplace {
	return r.ReplaceFloatTunable(old, new, 'f', -1, 64)
}

// Register string-float replacement.
func (r *BatchReplace) SReplaceFloat(old string, new float64) *BatchReplace {
	return r.SReplaceFloatTunable(old, new, 'f', -1, 64)
}

// Register bytes-float replacement with params.
func (r *BatchReplace) ReplaceFloatTunable(old []byte, new float64, fmt byte, prec, bitSize int) *BatchReplace {
	n := bytes.Count(r.indirect(r.src), old)
	if n == 0 {
		return r
	}
	r.old.add(r.alloc(old), n)

	c := r.off
	r.buf = strconv.AppendFloat(r.buf, new, fmt, prec, bitSize)
	r.off = len(r.buf)
	np := byteptr{
		o: c,
		l: r.off - c,
	}
	r.new.add(np, n)
	return r
}

// Register string-float replacement with params.
func (r *BatchReplace) SReplaceFloatTunable(old string, new float64, fmt byte, prec, bitSize int) *BatchReplace {
	return r.ReplaceFloatTunable(fastconv.S2B(old), new, fmt, prec, bitSize)
}

// Perform the replaces.
func (r *BatchReplace) Commit() []byte {
	// Calculate final length.
	bl := r.src.len() + r.new.acc
	l := bl - r.old.acc

	r.buf = bytealg.GrowDelta(r.buf, bl*2)
	r.dst.set(r.off, bl)
	dst := r.indirect(r.dst)
	copy(dst, r.indirect(r.src))
	r.off += bl
	buf := r.buf[r.off:]
	// Walk over queue and replace.
	for i := 0; i < len(r.old.queue); i++ {
		o := r.old.queue[i]
		n := r.new.queue[i]
		buf = r.replaceTo(buf[:0], dst, r.indirect(o.p), r.indirect(n.p), o.n)
		dst = append(dst[:0], buf...)
	}

	return r.indirect(r.dst)[:l]
}

func (r *BatchReplace) replaceTo(dst, s, old, new []byte, n int) []byte {
	start := 0
	for i := 0; i < n; i++ {
		j := start + bytes.Index(s[start:], old)
		dst = append(dst, s[start:j]...)
		dst = append(dst, new...)
		start = j + len(old)
	}
	dst = append(dst, s[start:]...)
	return dst
}

// Perform the replaces and return copy of result.
//
// Made to avoid data sharing.
func (r *BatchReplace) CommitCopy() []byte {
	return append([]byte(nil), r.Commit()...)
}

// String version of Commit().
// todo remove it due to SCommit() method.
func (r *BatchReplace) CommitStr() string {
	return fastconv.B2S(r.Commit())
}

// String version of Commit().
func (r *BatchReplace) SCommit() string {
	return fastconv.B2S(r.Commit())
}

// String version of CommitCopy().
// todo remove it due to SCommitCopy() method.
func (r *BatchReplace) CommitCopyStr() string {
	return fastconv.B2S(r.CommitCopy())
}

// String version of CommitCopy().
func (r *BatchReplace) SCommitCopy() string {
	return fastconv.B2S(r.CommitCopy())
}

// Clear the replacer with keeping of allocated space to reuse.
func (r *BatchReplace) Reset() *BatchReplace {
	r.buf = r.buf[:0]
	r.off = 0
	r.src.reset()
	r.dst.reset()
	for i := 0; i < len(r.old.queue); i++ {
		r.old.queue[i].p.reset()
		r.old.queue[i].n = 0
		r.old.idx, r.old.acc = 0, 0
		r.new.queue[i].p.reset()
		r.new.queue[i].n = 0
		r.new.idx, r.new.acc = 0, 0
	}
	return r
}

func (r *BatchReplace) indirect(p byteptr) []byte {
	return r.buf[p.offset() : p.offset()+p.len()]
}

func (r *BatchReplace) alloc(b []byte) (p byteptr) {
	c := r.off
	l := len(b)
	r.buf = bytealg.GrowDelta(r.buf, l)
	copy(r.buf[c:c+l], b)
	p.set(c, l)
	r.off += l
	return
}

// Add new item to queue.
func (q *queue) add(p byteptr, n int) {
	if n == 0 {
		n = 1
	}
	q.acc += p.len() * n
	x := byteptrn{p: p, n: n}
	if q.idx < q.cap {
		q.queue[q.idx] = x
		q.idx++
		return
	}
	q.queue = append(q.queue, x)
	q.idx++
	q.cap++
}
