package batch_replace

import "bytes"

// Make byte slice copy of given string.
func scopy(s string) []byte {
	return append([]byte(nil), s...)
}

// Generic replace in destination slice.
//
// Destination should has enough space (capacity).
func replaceTo(dst, s, old, new []byte, n int) []byte {
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
