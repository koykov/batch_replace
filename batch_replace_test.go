package batch_replace

import (
	"bytes"
	"strconv"
	"testing"
)

var (
	breplOrigin = []byte("foo {tag0} bar {tag1} string {macro} with {cnt} tags")
	breplExpect = []byte("foo s0 bar long string string 1234567.0987654321 with 4 tags")
	brTag0      = []byte("{tag0}")
	brTag0Val   = []byte("s0")
	brTag1      = []byte("{tag1}")
	brTag1Val   = []byte("long string")
	brTag2      = []byte("{macro}")
	brTag3      = []byte("{cnt}")
)

func TestBatchReplace_Replace(t *testing.T) {
	n := NewBatchReplace(breplOrigin).
		Replace(brTag0, brTag0Val).
		Replace(brTag1, brTag1Val).
		ReplaceFloat(brTag2, float64(1234567.0987654321)).
		ReplaceInt(brTag3, int64(4)).
		Commit()
	if !bytes.Equal(n, breplExpect) {
		t.Error("BatchReplace: mismatch result and expectation")
	}
}

func BenchmarkBatchReplace_Replace(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := BatchPool.Get(breplOrigin)
		n := r.Replace(brTag0, brTag0Val).
			Replace(brTag1, brTag1Val).
			ReplaceFloat(brTag2, float64(1234567.0987654321)).
			ReplaceInt(brTag3, int64(4)).
			Commit()
		if !bytes.Equal(n, breplExpect) {
			b.Error("BatchReplace: mismatch result and expectation")
		}
		BatchPool.Put(r)
	}
}

func BenchmarkBatchReplaceNative_Replace(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n := bytes.Replace(breplOrigin, brTag0, brTag0Val, -1)
		n = bytes.Replace(n, brTag1, brTag1Val, -1)
		n = bytes.Replace(n, brTag2, []byte(strconv.FormatFloat(1234567.0987654321, 'f', -1, 64)), -1)
		n = bytes.Replace(n, brTag3, []byte(strconv.Itoa(4)), -1)
		if !bytes.Equal(n, breplExpect) {
			b.Error("BatchReplace: mismatch result and expectation")
		}
	}
}
