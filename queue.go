package batch_replace

// Replace queue.
type queue struct {
	// Queue items.
	queue []byteptrn
	// Current index.
	idx int
	// Max index.
	cap int
	// Accumulated items length.
	acc int
}

// Add new item to queue.
func (q *queue) add(p byteptr, n int) {
	if n == 0 {
		n = 1
	}
	q.acc += p.len() * n
	x := byteptrn{p: p, n: n}
	if q.idx < q.cap {
		q.queue[q.idx] = x
		q.idx++
		return
	}
	q.queue = append(q.queue, x)
	q.idx++
	q.cap++
}
