package imgcat

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/pkg/errors"
)

type writer struct {
	pw   *io.PipeWriter
	done chan struct{}
}

func (w writer) Write(data []byte) (int, error) {
	return w.pw.Write(data)
}

func (w writer) Close() error {
	if err := w.pw.Close(); err != nil {
		return err
	}

	<-w.done
	return nil
}

// NewWriter receives a writer and returns a imgcat Writer
func NewWriter(w io.Writer) io.WriteCloser {
	pr, pw := io.Pipe()
	wc := writer{pw: pw, done: make(chan struct{})}

	go func() {
		defer close(wc.done)
		err := copy(w, pr)
		pr.CloseWithError(err)
	}()

	return wc
}

func copy(w io.Writer, r io.Reader) error {

	header := strings.NewReader("\033]1337;File=inline=1:")
	footer := strings.NewReader("\a\n")

	pr, pw := io.Pipe()

	go func() {
		wc := base64.NewEncoder(base64.StdEncoding, pw)
		_, err := io.Copy(wc, r)
		if err != nil {
			pw.CloseWithError(errors.Wrap(err, "error encoding image"))
			return
		}

		if err := wc.Close(); err != nil {
			pw.CloseWithError(errors.Wrap(err, "error closing wc"))
			return
		}
		pw.Close()
	}()

	_, err := io.Copy(w, io.MultiReader(header, pr, footer))
	return err
}
