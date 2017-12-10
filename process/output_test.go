package process_test

import (
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/ionrock/xenv/process"
)

func TestOutputFactory(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.Group("capture", func() {
		o.BeforeEach(func(t *testing.T) (*testing.T, *process.Output) {
			return t, &process.Output{Padding: 5, Capture: true}
		})

		o.Spec("writing a line stores the output", func(t *testing.T, of *process.Output) {
			of.WriteLine("foo", "hello world", false)
			Expect(t, of.Output()).To(Equal("hello world"))
		})
	})
}
