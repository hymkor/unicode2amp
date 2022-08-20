package main

import (
	"bufio"
	"flag"
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

var flagReplace = flag.Bool("i", false, "replace file")

func mainsReplace(args []string) error {
	for _, arg1 := range args {
		r, err := os.Open(arg1)
		if err != nil {
			return err
		}
		newFile := arg1 + ".tmp"
		w, err := os.Create(newFile)
		if err != nil {
			return err
		}
		bw := bufio.NewWriter(w)
		filterErr := filter(r, bw)
		bw.Flush()
		rCloseErr := r.Close()
		wCloseErr := w.Close()

		if filterErr != nil {
			return filterErr
		}
		if rCloseErr != nil {
			return rCloseErr
		}
		if wCloseErr != nil {
			return wCloseErr
		}
		if err = os.Rename(arg1, arg1+"~"); err != nil {
			return err
		}
		if err = os.Rename(newFile, arg1); err != nil {
			return err
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
	var err error
	flag.Parse()
	if *flagReplace {
		err = mainsReplace(flag.Args())
	} else {
		err = mains(flag.Args())
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
