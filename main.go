package main

import (
	"fmt"
	"io"
	"os"
	"github.com/pkg/errors"
	"github.com/vsabreu/imgcat/imgcat"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enought arguments to cat.")
		os.Exit(1)
	}

	for _, img := range os.Args[1:] {
		if err := cat(img); err != nil {
			fmt.Fprintf(os.Stderr, "Could not cat %s: %v\n", img, err)
		}
	}
	fmt.Println("finished")
}

func cat(img string) error {
	f, err := os.Open(img)
	if err != nil {
		return errors.Wrap(err, "Failed opening file.")
	}
	defer f.Close()

	wc := imgcat.NewWriter(os.Stdout)
	if _, err := io.Copy(wc, f); err != nil {
		return err
	}

	return wc.Close()
}

// use a badWriter in imgcat.NewWriter to simulate an error
type badWriter struct{}

func (badWriter) Write([]byte) (int, error) {
	return 0, errors.New("bad writer")
}
