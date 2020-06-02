package batch_replace

import (
	"strconv"
	"strings"
	"testing"
)

var (
	breplOriginS = "foo {tag0} bar {tag1} string {macro} with {cnt} tags"
	breplExpectS = "foo s0 bar long string string 1234567.0987654321 with 4 tags"
	brTag0S      = "{tag0}"
	brTag0ValS   = "s0"
	brTag1S      = "{tag1}"
	brTag1ValS   = "long string"
	brTag2S      = "{macro}"
	brTag3S      = "{cnt}"
)

func TestBatchReplaceStr_Replace(t *testing.T) {
	n := NewBatchReplaceStr("foo {tag0} bar {tag1} string {macro} with {cnt} tags").
		Replace("{tag0}", "s0").
		Replace("{tag1}", "long string").
		ReplaceFloat("{macro}", float64(1234567.0987654321)).
		ReplaceInt("{cnt}", int64(4)).
		Commit()
	if n != breplExpectS {
		t.Error("BatchReplaceStr: mismatch result and expectation")
	}
}

func BenchmarkBatchReplaceStr_Replace(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := BatchStrPool.Get(breplOriginS)
		n := r.Replace(brTag0S, brTag0ValS).
			Replace(brTag1S, brTag1ValS).
			ReplaceFloat(brTag2S, float64(1234567.0987654321)).
			ReplaceInt(brTag3S, int64(4)).
			Commit()
		if n != breplExpectS {
			b.Error("BatchReplaceStr: mismatch result and expectation")
		}
		BatchStrPool.Put(r)
	}
}

func BenchmarkBatchReplaceStrNative_Replace(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n := strings.Replace(breplOriginS, brTag0S, brTag0ValS, -1)
		n = strings.Replace(n, brTag1S, brTag1ValS, -1)
		n = strings.Replace(n, brTag2S, strconv.FormatFloat(1234567.0987654321, 'f', -1, 64), -1)
		n = strings.Replace(n, brTag3S, strconv.Itoa(4), -1)
		if n != breplExpectS {
			b.Error("BatchReplaceStr: mismatch result and expectation")
		}
	}
}
