package util_test

import (
	"fmt"
	"os/exec"
	"sync"
	"testing"

	"github.com/ionrock/xenv/util"
)

func TestLineReader(t *testing.T) {
	cmd := exec.Command("echo", "hello world")
	stdout, _ := cmd.StdoutPipe()

	var output string
	expected := fmt.Sprintf("true | hello world")

	outHandler := func(line string) string {
		output = fmt.Sprintf("true | %s", line)
		return output
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go util.LineReader(wg, stdout, outHandler)

	cmd.Start()
	cmd.Wait()

	if output != expected {
		t.Errorf("wrong output: %s != %s", output, expected)
	}
}
