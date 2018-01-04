package config_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/ionrock/xenv/config"
)

func TestGetConfig(t *testing.T) {
	data := map[string]string{
		"FOO": "bar",
	}
	c := config.Config{Data: data}

	if v := c.GetConfig("FOO"); v != "bar" {
		t.Errorf("config didn't contain expected value: %q != %q", v, "bar")
	}

	// Ensure we look up in the env when it doesn't exist in the
	// config.
	varName := "XE_CONFIG_TEST_VAR_FOO"
	os.Setenv(varName, "bar")
	defer os.Setenv(varName, "")
	if v := c.GetConfig(varName); v != "bar" {
		t.Errorf("config didn't contain expected value: %q != %q", v, "bar")
	}
}

func TestConfigGetAndSet(t *testing.T) {
	data := map[string]string{
		"FOO": "bar",
	}
	c := config.Config{Data: data}

	c.Set("FOO", "bar")

	if v, ok := c.Get("FOO"); v != "bar" || !ok {
		t.Errorf("config wasn't set correctly: %q != %q", v, "bar")
	}
}

func TestToEnv(t *testing.T) {
	data := map[string]string{
		"FOO": "bar",
	}
	c := config.Config{Data: data}

	if c.ToEnv()[0] != "FOO=bar" {
		t.Errorf("missing env vars from ToEnv")
	}
}

// TestEnvScriptOutput is a helper test that prints simple JSON to
// stdout for an xenv to read.
func TestWithEnvironSet(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	data := map[string]string{
		// Override bar from environ
		"BAR": "bar",

		// Add baz
		"BAZ": "baz",
	}
	c := config.Config{Data: data}

	for _, line := range c.Environ() {
		fmt.Println(line)
	}
	os.Exit(0)
}

// TestEnviron uses the scriptHelper to set an environment and
// allowing the config to use it when constructing a new environment with
// the config data.
func TestEnviron(t *testing.T) {
	env := []string{
		"FOO=foo",
		"BAR=baz",
	}
	cmd := scriptHelper("TestWithEnvironSet", env)
	result, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed running helper: %s", err)
	}

	data := [][]byte{
		[]byte("FOO=foo"),
		[]byte("BAR=bar"),
		[]byte("BAZ=baz"),
	}

	for _, expected := range data {
		if !bytes.Contains(result, expected) {
			t.Errorf("missing %s from environment", expected)
		}
	}
}
