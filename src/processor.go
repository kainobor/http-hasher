package src

import (
	"context"
	"sync"
)

type (
	Processor interface {
		Start(context.Context) <-chan struct{}
		Process(string)
		Shutdown(context.Context)
	}

	processor struct {
		finished   chan struct{}
		forProcess chan string
		parallel   int
		writer     Writer
		httpGetter HttpGetter
		hasher     Hasher
		wg         sync.WaitGroup
	}
)

var _ Processor = &processor{}

func NewProcessor(parallel int, w Writer, h Hasher, g HttpGetter) *processor {
	return &processor{
		forProcess: make(chan string, parallel),
		finished:   make(chan struct{}),
		parallel:   parallel,
		writer:     w,
		hasher:     h,
		httpGetter: g,
	}
}

func (p *processor) Start(ctx context.Context) <-chan struct{} {
	p.initPool(ctx)

	return p.finished
}

func (p *processor) Process(input string) {
	p.forProcess <- input
}

func (p *processor) Shutdown(ctx context.Context) {
	close(p.forProcess)
	p.wg.Wait()
	select {
	case <-ctx.Done():
		p.writer.Write("Gracefully shut down")
	default:
	}
}

func (p *processor) initPool(ctx context.Context) {
	for i := 0; i < p.parallel; i++ {
		p.wg.Add(1)
		go p.processInput(ctx)
	}
}

func (p *processor) processInput(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.wg.Done()
			return
		case input, ok := <-p.forProcess:
			if !ok {
				p.wg.Done()
				return
			}

			url, resp, err := p.httpGetter.HttpGet(input)
			if err != nil {
				p.writer.Write(input, err.Error())
				continue
			}

			hash, err := p.hasher.Hash(resp)
			if err != nil {
				p.writer.Write(input, err.Error())
				continue
			}

			p.writer.Write(url, hash)
		}
	}
}
