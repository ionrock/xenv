package process_test

import (
	"os"
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/ionrock/xenv/process"
)

func TestNewCmd(t *testing.T) {
	newCmdTests := []struct {
		cmd   *process.Proc
		parts []string
	}{
		// Make sure we split more than one part
		{process.New("/bin/echo 'hello world'"), []string{"/bin/echo", "hello world"}},

		// Make sure we only call with a path argument
		{process.New("/bin/echo"), []string{"/bin/echo"}},
	}

	for _, cmdTest := range newCmdTests {
		for i, part := range cmdTest.parts {
			Expect(t, cmdTest.cmd.Cmd.Args[i]).To(Equal(part))
		}
	}
}

func TestEnv(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) (*testing.T, string) {
		envvar := "PROC_ENVVAR_FALLBACK_TEST"
		os.Setenv("PROC_ENVVAR_FALLBACK_TEST", "bar")
		return t, envvar
	})

	o.Spec("include env", func(t *testing.T, envvar string) {
		env := map[string]string{"foo": "hello world"}
		environ := process.ParseEnv(process.Env(env, true))

		Expect(t, environ).To(HaveKey(envvar))
		Expect(t, environ).To(HaveKey("foo"))
	})
}
