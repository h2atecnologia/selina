package text

import (
	"bufio"
	"context"
	"errors"
	"io"

	"github.com/licaonfee/selina"
)

var _ selina.Worker = (*Reader)(nil)

//ReaderOptions customize Reader
type ReaderOptions struct {
	//Reader from which data is readed
	Reader io.Reader
	//AutoClose if its true and Reader implements io.Closer
	//io.Reader.Close() method is called on Process finish
	AutoClose bool
}

//Reader a worker that read data from an io.Reader
type Reader struct {
	opts ReaderOptions
}

func (t *Reader) cleanup() error {
	if t.opts.Reader == nil {
		return nil
	}
	if c, ok := t.opts.Reader.(io.Closer); t.opts.AutoClose && ok {
		return c.Close()
	}
	return nil
}

//ErrNilReader is returned when a nil io.Reader interface is provided
var ErrNilReader = errors.New("nil io.Reader provided to TextReader")

//Process implements Worker interface
func (t *Reader) Process(ctx context.Context, args selina.ProcessArgs) (err error) {
	defer func() {
		close(args.Output)
		cerr := t.cleanup()
		if err == nil { //if an error occurred not override it
			err = cerr
		}
	}()
	if t.opts.Reader == nil {
		return ErrNilReader
	}
	sc := bufio.NewScanner(t.opts.Reader)
	for sc.Scan() {
		select {
		case _, ok := <-args.Input:
			if !ok {
				return nil
			}
		case args.Output <- []byte(sc.Text()):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

//NewReader create a new Reader with given options
func NewReader(opts ReaderOptions) *Reader {
	t := Reader{opts: opts}
	return &t
}
