package batch_replace

import (
	"bytes"
	"strconv"
	"strings"
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

	breplOriginS = "foo {tag0} bar {tag1} string {macro} with {cnt} tags"
	breplExpectS = "foo s0 bar long string string 1234567.0987654321 with 4 tags"
	brTag0S      = "{tag0}"
	brTag0ValS   = "s0"
	brTag1S      = "{tag1}"
	brTag1ValS   = "long string"
	brTag2S      = "{macro}"
	brTag3S      = "{cnt}"
)

func TestBatchReplace_Replace(t *testing.T) {
	n := NewBatchReplace(breplOrigin).
		Replace(brTag0, brTag0Val).
		Replace(brTag1, brTag1Val).
		ReplaceFloat(brTag2, 1234567.0987654321).
		ReplaceInt(brTag3, int64(4)).
		Commit()
	if !bytes.Equal(n, breplExpect) {
		t.Error("BatchReplace: mismatch result and expectation")
	}
}

func TestBatchReplaceStr_Replace(t *testing.T) {
	n := NewBatchReplace(nil).
		SetSrcStr("foo {tag0} bar {tag1} string {macro} with {cnt} tags").
		SReplace("{tag0}", "s0").
		SReplace("{tag1}", "long string").
		SReplaceFloat("{macro}", float64(1234567.0987654321)).
		SReplaceInt("{cnt}", int64(4)).
		SCommit()
	if n != breplExpectS {
		t.Error("BatchReplace: mismatch string result and expectation")
	}
}

func BenchmarkBatchReplace_Replace(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := Acquire(breplOrigin)
		n := r.Replace(brTag0, brTag0Val).
			Replace(brTag1, brTag1Val).
			ReplaceFloat(brTag2, 1234567.0987654321).
			ReplaceInt(brTag3, int64(4)).
			Commit()
		if !bytes.Equal(n, breplExpect) {
			b.Error("BatchReplace: mismatch result and expectation")
		}
		Release(r)
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

func BenchmarkBatchReplaceStr_Replace(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := SAcquire(breplOriginS)
		n := r.SReplace(brTag0S, brTag0ValS).
			SReplace(brTag1S, brTag1ValS).
			SReplaceFloat(brTag2S, float64(1234567.0987654321)).
			SReplaceInt(brTag3S, int64(4)).
			CommitStr()
		if n != breplExpectS {
			b.Error("BatchReplace: mismatch string result and expectation")
		}
		Release(r)
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
