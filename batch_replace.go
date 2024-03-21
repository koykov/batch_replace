package batch_replace

import (
	"bytes"
	"strconv"

	"github.com/koykov/bytealg"
	"github.com/koykov/byteconv"
	"github.com/koykov/byteseq"
)

const (
	// Int base edges.
	baseLo = 2
	baseHi = 36
)

// BatchReplace is a tool to collect replacement pairs and apply them once..
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

// NewBatchReplace inits new replacer.
func NewBatchReplace[T byteseq.Byteseq](x T) *BatchReplace {
	var src []byte
	if b, ok := byteseq.ToBytes(x); ok {
		src = b
	}
	if s, ok := byteseq.ToString(x); ok {
		src = byteconv.S2B(s)
	}
	o := queue{queue: make([]byteptrn, 0)}
	n := queue{queue: make([]byteptrn, 0)}
	r := BatchReplace{
		old: o,
		new: n,
	}
	r.SetSource(src)
	return &r
}

// SetSource set the source as bytes.
//
// For use outside of pools.
func (r *BatchReplace) SetSource(src []byte) *BatchReplace {
	r.buf = append(r.buf[:0], src...)
	r.src.set(r.off, len(src))
	r.off = len(src)
	return r
}

// SetSourceString set the source as string.
func (r *BatchReplace) SetSourceString(src string) *BatchReplace {
	return r.SetSource(byteconv.S2B(src))
}

// SetSrcBytes set the source as bytes.
//
// For use outside of pools.
// Deprecated: use SetSource() instead.
func (r *BatchReplace) SetSrcBytes(src []byte) *BatchReplace {
	return r.SetSource(src)
}

// SetSrcStr set the source as string.
// Deprecated: use SetSourceString() instead.
func (r *BatchReplace) SetSrcStr(src string) *BatchReplace {
	return r.SetSourceString(src)
}

// BytesToBytes registers new bytes to bytes replacement.
func (r *BatchReplace) BytesToBytes(old []byte, new []byte) *BatchReplace {
	n := bytes.Count(r.indirect(r.src), old)
	if n == 0 {
		return r
	}
	r.old.add(r.alloc(old), n)
	r.new.add(r.alloc(new), n)
	return r
}

// BytesToStr registers new bytes to bytes to string replacement.
func (r *BatchReplace) BytesToStr(old []byte, new string) *BatchReplace {
	return r.BytesToBytes(old, byteconv.S2B(new))
}

// StrToStr registers new string to string replacement.
func (r *BatchReplace) StrToStr(old, new string) *BatchReplace {
	return r.BytesToBytes(byteconv.S2B(old), byteconv.S2B(new))
}

// StrToBytes registers new bytes to string to bytes replacement.
func (r *BatchReplace) StrToBytes(old string, new []byte) *BatchReplace {
	return r.BytesToBytes(byteconv.S2B(old), new)
}

// BytesToInt registers bytes to int replacement.
func (r *BatchReplace) BytesToInt(old []byte, new int64) *BatchReplace {
	return r.BytesToIntBase(old, new, 10)
}

// StrToInt segisters string to int replacement.
func (r *BatchReplace) StrToInt(old string, new int64) *BatchReplace {
	return r.StrToIntBase(old, new, 10)
}

// BytesToIntBase registers bytes to int replacement with given base.
func (r *BatchReplace) BytesToIntBase(old []byte, new int64, base int) *BatchReplace {
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

// StrToIntBase registers string to int replacement with given base.
func (r *BatchReplace) StrToIntBase(old string, new int64, base int) *BatchReplace {
	return r.BytesToIntBase(byteconv.S2B(old), new, base)
}

// BytesToUint registers bytes to uint replacement.
func (r *BatchReplace) BytesToUint(old []byte, new uint64) *BatchReplace {
	return r.BytesToUintBase(old, new, 10)
}

// StrToUint registers string to uint replacement.
func (r *BatchReplace) StrToUint(old string, new uint64) *BatchReplace {
	return r.StrToUintBase(old, new, 10)
}

// BytesToUintBase registers bytes to uint replacement with given base.
func (r *BatchReplace) BytesToUintBase(old []byte, new uint64, base int) *BatchReplace {
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

// StrToUintBase registers string to uint replacement with given base.
func (r *BatchReplace) StrToUintBase(old string, new uint64, base int) *BatchReplace {
	return r.BytesToUintBase(byteconv.S2B(old), new, base)
}

// BytesToFloat registers bytes to float replacement.
func (r *BatchReplace) BytesToFloat(old []byte, new float64) *BatchReplace {
	return r.BytesToFloatTunable(old, new, 'f', -1, 64)
}

// StrToFloat registers string to float replacement.
func (r *BatchReplace) StrToFloat(old string, new float64) *BatchReplace {
	return r.StrToFloatTunable(old, new, 'f', -1, 64)
}

// BytesToFloatTunable registers bytes to float replacement with params.
func (r *BatchReplace) BytesToFloatTunable(old []byte, new float64, fmt byte, prec, bitSize int) *BatchReplace {
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

// StrToFloatTunable registers string to float replacement with params.
func (r *BatchReplace) StrToFloatTunable(old string, new float64, fmt byte, prec, bitSize int) *BatchReplace {
	return r.BytesToFloatTunable(byteconv.S2B(old), new, fmt, prec, bitSize)
}

// Commit applies all registered replacement pairs.
func (r *BatchReplace) Commit() []byte {
	// Calculate final length.
	bl := r.src.len() + r.new.acc
	l := bl - r.old.acc

	r.buf = bytealg.GrowDelta(r.buf, r.off+bl*2)
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

// CommitCopy perform the replaces and return copy of result.
//
// Made to avoid data sharing.
// See Commit.
func (r *BatchReplace) CommitCopy() []byte {
	return append([]byte(nil), r.Commit()...)
}

// CommitStr is a string version of Commit().
func (r *BatchReplace) CommitStr() string {
	return byteconv.B2S(r.Commit())
}

// CommitCopyStr os a string version of CommitCopy().
func (r *BatchReplace) CommitCopyStr() string {
	return byteconv.B2S(r.CommitCopy())
}

// Reset clears the replacer with keeping of allocated space to reuse.
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

// Replace old to new in s and apply the result to dst.
func (r *BatchReplace) replaceTo(dst, s, old, new []byte, n int) []byte {
	start := 0
	for i := 0; i < n; i++ {
		j := start + bytes.Index(s[start:], old)
		if j == -1 {
			continue
		}
		dst = append(dst, s[start:j]...)
		dst = append(dst, new...)
		start = j + len(old)
	}
	dst = append(dst, s[start:]...)
	return dst
}

// Get byte slice according byte pointer.
func (r *BatchReplace) indirect(p byteptr) []byte {
	return r.buf[p.offset() : p.offset()+p.len()]
}

// Alloc more space (or use exiting) in buffer and return corresponding byte pointer.
func (r *BatchReplace) alloc(b []byte) (p byteptr) {
	c := r.off
	l := len(b)
	r.buf = bytealg.GrowDelta(r.buf, l)
	copy(r.buf[c:c+l], b)
	p.set(c, l)
	r.off += l
	return
}
