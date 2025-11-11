package batch_replace

import "github.com/koykov/byteconv"

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

// BytesToStr registers new bytes-to-bytes to string replacement.
// Deprecated: use BytesToString() instead.
func (r *BatchReplace) BytesToStr(old []byte, new string) *BatchReplace {
	return r.BytesToString(old, new)
}

// StrToStr registers new string-to-string replacement.
// Deprecated: use StringToString() instead.
func (r *BatchReplace) StrToStr(old, new string) *BatchReplace {
	return r.StringToString(old, new)
}

// StrToBytes registers new bytes to string to bytes replacement.
// Deprecated: use StringToBytes() instead.
func (r *BatchReplace) StrToBytes(old string, new []byte) *BatchReplace {
	return r.StringToBytes(old, new)
}

// StrToInt registers string to int replacement.
// Deprecated: use StringToInt() instead.
func (r *BatchReplace) StrToInt(old string, new int64) *BatchReplace {
	return r.StringToInt(old, new)
}

// StrToIntBase registers string to int replacement with given base.
// Deprecated: use StringToIntBase() instead.
func (r *BatchReplace) StrToIntBase(old string, new int64, base int) *BatchReplace {
	return r.StringToIntBase(old, new, base)
}

// StrToUint registers string to uint replacement.
// Deprecated: use StringToUint() instead.
func (r *BatchReplace) StrToUint(old string, new uint64) *BatchReplace {
	return r.StringToUint(old, new)
}

// StrToUintBase registers string to uint replacement with given base.
// Deprecated: use StringToUintBase() instead.
func (r *BatchReplace) StrToUintBase(old string, new uint64, base int) *BatchReplace {
	return r.StringToUintBase(old, new, base)
}

// StrToFloat registers string to float replacement.
// Deprecated: use StringToFloat() instead.
func (r *BatchReplace) StrToFloat(old string, new float64) *BatchReplace {
	return r.StringToFloatTunable(old, new, 'f', -1, 64)
}

// StrToFloatTunable registers string to float replacement with params.
// Deprecated: use StringToFloatTunable() instead.
func (r *BatchReplace) StrToFloatTunable(old string, new float64, fmt byte, prec, bitSize int) *BatchReplace {
	return r.StringToFloatTunable(old, new, fmt, prec, bitSize)
}

// CommitStr is a string version of Commit().
// Deprecated: use CommitCopy() instead.
func (r *BatchReplace) CommitStr() string {
	return byteconv.B2S(r.Commit())
}

// CommitCopyStr os a string version of CommitCopy().
// Deprecated: use CommitCopyString() instead.
func (r *BatchReplace) CommitCopyStr() string {
	return byteconv.B2S(r.CommitCopy())
}
