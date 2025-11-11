package batch_replace

// @file
// Shorthand alias of replace methods.

// B2B is an alias of BytesToBytes.
func (r *BatchReplace) B2B(old []byte, new []byte) *BatchReplace {
	return r.BytesToBytes(old, new)
}

// B2S is an alias of BytesToString.
func (r *BatchReplace) B2S(old []byte, new string) *BatchReplace {
	return r.BytesToString(old, new)
}

// S2S is an alias of StringToString.
func (r *BatchReplace) S2S(old, new string) *BatchReplace {
	return r.StringToString(old, new)
}

// S2B is an alias of StringToBytes.
func (r *BatchReplace) S2B(old string, new []byte) *BatchReplace {
	return r.StringToBytes(old, new)
}

// B2I is an alias of BytesToInt.
func (r *BatchReplace) B2I(old []byte, new int64) *BatchReplace {
	return r.BytesToInt(old, new)
}

// S2I is an alias of StringToInt.
func (r *BatchReplace) S2I(old string, new int64) *BatchReplace {
	return r.StringToInt(old, new)
}

// B2IB is an alias of BytesToIntBase.
func (r *BatchReplace) B2IB(old []byte, new int64, base int) *BatchReplace {
	return r.BytesToIntBase(old, new, base)
}

// S2IB is an alias of StringToIntBase.
func (r *BatchReplace) S2IB(old string, new int64, base int) *BatchReplace {
	return r.StringToIntBase(old, new, base)
}

// B2U is an alias of BytesToUint.
func (r *BatchReplace) B2U(old []byte, new uint64) *BatchReplace {
	return r.BytesToUint(old, new)
}

// S2U is an alias of StringToUint.
func (r *BatchReplace) S2U(old string, new uint64) *BatchReplace {
	return r.StringToUint(old, new)
}

// B2UB is an alias of BytesToUintBase.
func (r *BatchReplace) B2UB(old []byte, new uint64, base int) *BatchReplace {
	return r.BytesToUintBase(old, new, base)
}

// S2UB is an alias of StringToUintBase.
func (r *BatchReplace) S2UB(old string, new uint64, base int) *BatchReplace {
	return r.StringToUintBase(old, new, base)
}

// B2F is an alias of BytesToFloat.
func (r *BatchReplace) B2F(old []byte, new float64) *BatchReplace {
	return r.BytesToFloat(old, new)
}

// S2F is an alias of StringToFloat.
func (r *BatchReplace) S2F(old string, new float64) *BatchReplace {
	return r.StringToFloat(old, new)
}

// B2FT is an alias of BytesToFloatTunable.
func (r *BatchReplace) B2FT(old []byte, new float64, fmt byte, prec, bitSize int) *BatchReplace {
	return r.BytesToFloatTunable(old, new, fmt, prec, bitSize)
}

// S2FT is an alias of StringToFloatTunable.
func (r *BatchReplace) S2FT(old string, new float64, fmt byte, prec, bitSize int) *BatchReplace {
	return r.StringToFloatTunable(old, new, fmt, prec, bitSize)
}
