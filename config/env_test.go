package config_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ionrock/xenv/config"
)

func TestDataHandler(t *testing.T) {
	e := config.NewEnvironment()

	e.DataHandler(&config.XeConfig{Env: map[string]string{"FOO": "foo"}})

	if _, ok := e.Config.Get("FOO"); !ok {
		t.Error("error mapping EnvScript to data")
	}
}

func scriptCmd(name string, s ...string) []string {
	cmd := []string{os.Args[0], fmt.Sprintf("-test.run=%s", name)}
	if len(s) == 0 {
		return cmd
	}
	cmd = append(cmd, "--")
	cmd = append(cmd, s...)
	return cmd
}

func scriptHelper(name string, env []string) *exec.Cmd {
	parts := scriptCmd(name)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	cmd.Env = append(cmd.Env, env...)
	return cmd
}

func TestScriptEchoBar(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Println(os.Getenv("BAR"))
	os.Exit(0)
}

func TestScriptEnvInjected(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	if os.Getenv("ENV_INJECTED") != "true" {
		os.Exit(1)
	}

	os.Exit(0)
}

func TestSetEnv(t *testing.T) {
	e := config.NewEnvironment()

	// Set bar to a value we'll use in our script to ensure that our
	// config is linearly applied throughout the processing.
	e.SetEnv("BAR", "hello")

	// Set our process envvar to skip the test
	e.SetEnv("GO_WANT_HELPER_PROCESS", "1")

	// Get our command using our test executable
	val := fmt.Sprintf("`%s`", strings.Join(scriptCmd("TestScriptEchoBar"), " "))
	fmt.Println(val)

	// Set the result to FOO
	err := e.SetEnv("FOO", val)

	if err != nil {
		t.Fatalf("error running script: %s", err)
	}

	val, ok := e.Config.Get("FOO")
	if !ok {
		t.Errorf("error setting env var FOO from script")
	}

	if val != "hello" {
		t.Errorf("error with script result: %s != hello", val)
	}
}

func TestTaskGetsEnv(t *testing.T) {
	e := config.NewEnvironment()
	e.SetEnv("GO_WANT_HELPER_PROCESS", "1")
	e.SetEnv("ENV_INJECTED", "true")

	cmd := scriptCmd("TestScriptEnvInjected")

	err := e.RunTask("test_get_env", strings.Join(cmd, " "), "")
	if err != nil {
		t.Errorf("env not injected into task")
	}
}

func TestSetEnvFromScriptExpandsVars(t *testing.T) {
	e := config.NewEnvironment()
	e.SetEnv("GREETING", "world")
	err := e.SetEnvFromScript("cat testdata/script_out_with_vars.yml", ".")
	if err != nil {
		t.Errorf("error running script to update env: %s", err)
	}

	result, _ := e.Config.Get("HELLO")
	if result != "world" {
		t.Errorf("error setting value from existing env: %s", result)
	}
}
