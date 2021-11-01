# BatchReplace

BR is a tool to make replaces in bytes/string handy and alloc-free.

In fact it isn't a replacement of [strings.Replacer](https://golang.org/pkg/strings/#Replacer) since vanilla replacer
made for concurrent use, whereas BatchReplacer made to reduce allocations for big lists of replacements.

Usage example:

```go
originalStr := "this IS a string that contains {tag0}, {tag1}, tag2 and #s"
expectStr := "this WAS a string that contains 'very long substring', 1234567890, 154.195628217573 and etc..."

// Use pool instead of direct using of NewBatchReplace() or NewBatchReplaceStr().
// Pool may help you to get zero allocations on long distance and under high load.
r := batch_replace.AcquireWithStrSrc(originalStr)
defer batch_replace.Release(r)
res := r.StrToStr("IS", "WAS").
    S2S("{tag0}", "'very long substring'").
    StrToInt("{tag1}", int64(1234567890)).
    S2F("tag2", float64(154.195628217573)).
    S2S("#s", "etc...").
    Commit()
fmt.Println(res == expectStr) // true
```

## Benchmarks

```
BenchmarkBatchReplace/b2x-8         	 1566085	       742.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkBatchReplace/b2x_native-8  	 1570418	       729.9 ns/op	     288 B/op	       6 allocs/op
BenchmarkBatchReplace/s2x-8         	 1543791	       776.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkBatchReplace/s2x_native-8  	 1707190	       721.6 ns/op	     288 B/op	       6 allocs/op
BenchmarkBatchReplace/no_alloc-8    	 2998296	       409.1 ns/op	       0 B/op	       0 allocs/op
```
