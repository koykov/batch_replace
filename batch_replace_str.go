package batch_replace

import (
	fc "github.com/koykov/fastconv"
)

// String replacer.
type BatchReplaceStr struct {
	r *BatchReplace
}

// Init new string replacer.
func NewBatchReplaceStr(s string) *BatchReplaceStr {
	r := BatchReplaceStr{
		r: NewBatchReplace(scopy(s)),
	}
	return &r
}

// Set the source.
//
// For use outside of pools.
func (r *BatchReplaceStr) SetSrc(s string) {
	r.r.SetSrc(fc.S2B(s))
}

// Register new string replacement.
func (r *BatchReplaceStr) Replace(old string, new string) *BatchReplaceStr {
	r.r.Replace(fc.S2B(old), fc.S2B(new))
	return r
}

// Register int replacement.
func (r *BatchReplaceStr) ReplaceInt(old string, new int64) *BatchReplaceStr {
	return r.ReplaceIntBase(old, new, 10)
}

// Register int replacement with given base.
func (r *BatchReplaceStr) ReplaceIntBase(old string, new int64, base int) *BatchReplaceStr {
	r.r.ReplaceIntBase(fc.S2B(old), new, base)
	return r
}

// Register uint replacement.
func (r *BatchReplaceStr) ReplaceUint(old string, new uint64) *BatchReplaceStr {
	return r.ReplaceUintBase(old, new, 10)
}

// Register uint replacement with given base.
func (r *BatchReplaceStr) ReplaceUintBase(old string, new uint64, base int) *BatchReplaceStr {
	r.r.ReplaceUintBase(fc.S2B(old), new, base)
	return r
}

// Register float replacement.
func (r *BatchReplaceStr) ReplaceFloat(old string, new float64) *BatchReplaceStr {
	return r.ReplaceFloatTunable(old, new, 'f', -1, 64)
}

// Register float replacement with params.
func (r *BatchReplaceStr) ReplaceFloatTunable(old string, new float64, fmt byte, prec, bitSize int) *BatchReplaceStr {
	r.r.ReplaceFloatTunable(fc.S2B(old), new, fmt, prec, bitSize)
	return r
}

// Perform the replaces.
func (r *BatchReplaceStr) Commit() string {
	return fc.B2S(r.r.Commit())
}

// Clear the replacer with keeping of allocated space to reuse.
func (r *BatchReplaceStr) Reset() {
	r.r.Reset()
}
