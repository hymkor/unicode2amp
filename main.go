package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func filter(reader io.Reader, bw *bufio.Writer) error {
	br := bufio.NewReader(reader)
	for {
		r, _, err := br.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if r > 0xFFFF {
			fmt.Fprintf(bw, "&#x%X;", int(r))
		} else {
			bw.WriteRune(r)
		}
	}
	return nil
}

func mains(args []string) error {
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	if len(args) <= 0 {
		return filter(os.Stdin, bw)
	}
	for _, arg1 := range args {
		fd, err := os.Open(arg1)
		if err != nil {
			return err
		}
		err = filter(fd, bw)
		closeErr := fd.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
