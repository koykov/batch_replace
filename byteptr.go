package batch_replace

type byteptr struct {
	o, l int
}

func (p *byteptr) set(o, l int) {
	p.o, p.l = o, l
}

func (p *byteptr) len() int {
	return p.l
}

func (p *byteptr) offset() int {
	return p.o
}

func (p *byteptr) reset() {
	p.o, p.l = 0, 0
}
