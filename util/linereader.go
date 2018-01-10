package util

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

// OutHandler provides a simple function for processing a line of
// output from a process.
type OutHandler func(string) string

// LineReader can be started in a go routine to watch stdout/stderr
// and print each line prefixed with the Task Name.
func LineReader(wg *sync.WaitGroup, r io.Reader, handler OutHandler) {
	defer wg.Done()

	reader := bufio.NewReader(r)
	var buffer bytes.Buffer

	for {
		buf := make([]byte, 1024)

		n, err := reader.Read(buf)
		if err != nil {
			return
		}

		buf = buf[:n]

		for {
			i := bytes.IndexByte(buf, '\n')
			if i < 0 {
				break
			}

			buffer.Write(buf[0:i])
			handler(buffer.String())
			buffer.Reset()
			buf = buf[i+1:]
		}
		buffer.Write(buf)
	}
}
