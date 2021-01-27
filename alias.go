package batch_replace

// @file
// Shorthand alias of replace methods.

// Alias of BytesToBytes.
func (r *BatchReplace) B2B(old []byte, new []byte) *BatchReplace {
	return r.BytesToBytes(old, new)
}

// Alias of BytesToStr.
func (r *BatchReplace) B2S(old []byte, new string) *BatchReplace {
	return r.BytesToStr(old, new)
}

// Alias of StrToStr.
func (r *BatchReplace) S2S(old, new string) *BatchReplace {
	return r.StrToStr(old, new)
}

// Alias of StrToBytes.
func (r *BatchReplace) S2B(old string, new []byte) *BatchReplace {
	return r.StrToBytes(old, new)
}

// Alias of BytesToInt.
func (r *BatchReplace) B2I(old []byte, new int64) *BatchReplace {
	return r.BytesToInt(old, new)
}

// Alias of StrToInt.
func (r *BatchReplace) S2I(old string, new int64) *BatchReplace {
	return r.StrToInt(old, new)
}

// Alias of BytesToIntBase.
func (r *BatchReplace) B2IB(old []byte, new int64, base int) *BatchReplace {
	return r.BytesToIntBase(old, new, base)
}

// Alias of StrToIntBase.
func (r *BatchReplace) S2IB(old string, new int64, base int) *BatchReplace {
	return r.StrToIntBase(old, new, base)
}

// Alias of BytesToUint.
func (r *BatchReplace) B2U(old []byte, new uint64) *BatchReplace {
	return r.BytesToUint(old, new)
}

// Alias of StrToUint.
func (r *BatchReplace) S2U(old string, new uint64) *BatchReplace {
	return r.StrToUint(old, new)
}

// Alias of BytesToUintBase.
func (r *BatchReplace) B2UB(old []byte, new uint64, base int) *BatchReplace {
	return r.BytesToUintBase(old, new, base)
}

// Alias of StrToUintBase.
func (r *BatchReplace) S2UB(old string, new uint64, base int) *BatchReplace {
	return r.StrToUintBase(old, new, base)
}

// Alias of BytesToFloat.
func (r *BatchReplace) B2F(old []byte, new float64) *BatchReplace {
	return r.BytesToFloat(old, new)
}

// Alias of StrToFloat.
func (r *BatchReplace) S2F(old string, new float64) *BatchReplace {
	return r.StrToFloat(old, new)
}

// Alias of BytesToFloatTunable.
func (r *BatchReplace) B2FT(old []byte, new float64, fmt byte, prec, bitSize int) *BatchReplace {
	return r.BytesToFloatTunable(old, new, fmt, prec, bitSize)
}

// Alias of StrToFloatTunable.
func (r *BatchReplace) S2FT(old string, new float64, fmt byte, prec, bitSize int) *BatchReplace {
	return r.StrToFloatTunable(old, new, fmt, prec, bitSize)
}
