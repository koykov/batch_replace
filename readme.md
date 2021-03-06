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
r := batch_replace.SAcquire(originalStr)
defer batch_replace.Release(r)
res := r.Replace("IS", "WAS").
    Replace("{tag0}", "'very long substring'").
    ReplaceInt("{tag1}", int64(1234567890)).
    ReplaceFloat("tag2", float64(154.195628217573)).
    Replace("#s", "etc...").
    Commit()
fmt.Println(res == expectStr) // true
```

## Benchmarks

```
BenchmarkBatchReplace_Replace-8             1613643       755 ns/op       0 B/op       0 allocs/op
BenchmarkBatchReplaceNative_Replace-8       1584399       713 ns/op     304 B/op       6 allocs/op
BenchmarkBatchReplaceStr_Replace-8          1488924       791 ns/op       0 B/op       0 allocs/op
BenchmarkBatchReplaceStrNative_Replace-8    1363946       840 ns/op     544 B/op      10 allocs/op
```
