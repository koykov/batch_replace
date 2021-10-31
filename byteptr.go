package batch_replace

// Byte pointer.
type byteptr struct {
	// Offset and length of byte data in array.
	o, l int
}

// Byteptr struct with number of replacements.
type byteptrn struct {
	p byteptr
	n int
}

// Set offset and length.
func (p *byteptr) set(o, l int) {
	p.o, p.l = o, l
}

// Get length of data.
func (p *byteptr) len() int {
	return p.l
}

// Get offset in byte array.
func (p *byteptr) offset() int {
	return p.o
}

// Reset the pointer.
func (p *byteptr) reset() {
	p.o, p.l = 0, 0
}
