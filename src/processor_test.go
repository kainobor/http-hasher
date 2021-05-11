package src

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestProcessInput_HttpErr(t *testing.T) {
	ctx := context.Background()
	w := &writerMock{t: t}
	p := NewProcessor(1, w, &hasherMock{}, &getterMock{})
	p.Start(ctx)

	w.Expect("bad-site.com", "error")
	p.forProcess <- "bad-site.com"

	close(p.forProcess)
	p.wg.Wait()
}

func TestProcessInput_HashErr(t *testing.T) {
	ctx := context.Background()
	w := &writerMock{t: t}
	p := NewProcessor(1, w, &hasherMock{}, &getterMock{})
	p.Start(ctx)

	w.Expect("bad-resp.com", "error")
	p.forProcess <- "bad-resp.com"

	close(p.forProcess)
	p.wg.Wait()
}

func TestProcessInput_Ok(t *testing.T) {
	ctx := context.Background()
	w := &writerMock{t: t}
	p := NewProcessor(1, w, &hasherMock{}, &getterMock{})
	p.Start(ctx)

	w.Expect("http://google.com", "7c4f29407893c334a6cb7a87bf045c0d")
	p.forProcess <- "google.com"

	close(p.forProcess)
	p.wg.Wait()
}

type getterMock struct{}

func (_ *getterMock) HttpGet(uri string) (resultURL string, resp []byte, respErr error) {
	switch uri {
	case "bad-site.com":
		return "", nil, fmt.Errorf("error")
	case "bad-resp.com":
		return "http://bad-resp.com", []byte("bad-resp"), nil
	default:
		return "http://google.com", []byte("some response"), nil
	}
}

type hasherMock struct{}

func (_ *hasherMock) Hash(input []byte) (string, error) {
	switch string(input) {
	case "bad-resp":
		return "", fmt.Errorf("error")
	default:
		return "7c4f29407893c334a6cb7a87bf045c0d", nil
	}
}

type writerMock struct {
	exp []string
	t   *testing.T
}

func (w *writerMock) Expect(output ...string) {
	w.exp = output
}

func (w *writerMock) Write(got ...string) {
	if !reflect.DeepEqual(w.exp, got) {
		w.t.Errorf("got: %#v\nwanted: %#v\n", got, w.exp)
	}
	w.exp = nil
}
