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

	breplAllocS = "Balance of !user: !val !cur"
)

func TestBatchReplace(t *testing.T) {
	t.Run("b2x", func(t *testing.T) {
		n := NewBatchReplace(breplOrigin).
			BytesToBytes(brTag0, brTag0Val).
			BytesToBytes(brTag1, brTag1Val).
			BytesToFloat(brTag2, 1234567.0987654321).
			BytesToInt(brTag3, int64(4)).
			Commit()
		if !bytes.Equal(n, breplExpect) {
			t.Error("BatchReplace: mismatch result and expectation")
		}
	})
	t.Run("s2x", func(t *testing.T) {
		n := NewBatchReplace("foo {tag0} bar {tag1} string {macro} with {cnt} tags").
			StrToStr("{tag0}", "s0").
			StrToStr("{tag1}", "long string").
			StrToFloat("{macro}", 1234567.0987654321).
			StrToInt("{cnt}", int64(4)).
			CommitStr()
		if n != breplExpectS {
			t.Error("BatchReplace: mismatch string result and expectation")
		}
	})
}

func BenchmarkBatchReplace(b *testing.B) {
	b.Run("b2x", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			r := AcquireWithSource(breplOrigin)
			n := r.BytesToBytes(brTag0, brTag0Val).
				BytesToBytes(brTag1, brTag1Val).
				BytesToFloat(brTag2, 1234567.0987654321).
				BytesToInt(brTag3, int64(4)).
				Commit()
			if !bytes.Equal(n, breplExpect) {
				b.Error("BatchReplace: mismatch result and expectation")
			}
			Release(r)
		}
	})
	b.Run("b2x native", func(b *testing.B) {
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
	})
	b.Run("s2x", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			r := AcquireWithSource(breplOriginS)
			n := r.StrToStr(brTag0S, brTag0ValS).
				StrToStr(brTag1S, brTag1ValS).
				StrToFloat(brTag2S, 1234567.0987654321).
				StrToInt(brTag3S, int64(4)).
				CommitStr()
			if n != breplExpectS {
				b.Error("BatchReplace: mismatch string result and expectation")
			}
			Release(r)
		}
	})
	b.Run("s2x native", func(b *testing.B) {
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
	})
	b.Run("no alloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			r := AcquireWithSource(breplAllocS)
			n := r.StrToStr("!user", "John Ruth").
				StrToStr("!val", "8000").
				StrToStr("!cur", "USD").
				CommitStr()
			if n != "Balance of John Ruth: 8000 USD" {
				b.Error("BatchReplace + avoid alloc: mismatch string result and expectation")
			}
			Release(r)
		}
	})
}
