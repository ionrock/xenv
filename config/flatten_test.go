package config_test

import (
	"strings"
	"testing"

	"github.com/ionrock/xenv/config"
)

func TestLoadEnvFiles(t *testing.T) {
	paths := []string{
		"testdata/list_of_maps.json",
		"testdata/maps.json",
		"testdata/list_of_maps.yml",
		"testdata/maps.yml",
	}

	for _, path := range paths {
		env, err := config.NewFlatEnv(path)
		if err != nil {
			t.Fatalf("failed to load %q: %q", path, err)
		}

		v, ok := env["FOO"]

		if !ok {
			t.Errorf("key missing FOO in %q", env)
		}

		if v != "bar" {
			t.Errorf("value is wrong: %q != bar; %q", env["FOO"], env)
		}
	}
}

func TestFlattenKey(t *testing.T) {
	fe := &config.FlatEnv{
		Env: make(map[string]string),
	}

	prefix := []string{"FOO", "BAR", "BAZ"}
	err := fe.Load("hello world", prefix)
	if err != nil {
		t.Fatalf("error checking prefix size: %q", err)
	}

	expected := "FOO_BAR_BAZ"
	if _, ok := fe.Env[expected]; !ok {
		t.Errorf("flatKey failed: expected %s in %#v", expected, fe.Env)
	}

}

func TestNestedMaps(t *testing.T) {
	path := "testdata/nested_maps.yml"
	env, err := config.NewFlatEnv(path)
	if err != nil {
		t.Fatalf("failed to load %q: %q", path, err)
	}

	if _, ok := env["FOO_BAR_BAZ"]; !ok {
		t.Fatalf("error getting key: %#v", env)
	}
}

func TestListValues(t *testing.T) {
	path := "testdata/list_values.yml"
	env, err := config.NewFlatEnv(path)
	if err != nil {
		t.Fatalf("failed to load %q: %q", path, err)
	}
	val, ok := env["FOO"]
	if !ok {
		t.Fatalf("error getting key: %#v", env)
	}

	expected := []string{"one", "two", "three"}
	result := strings.Fields(val)
	if len(result) != len(expected) {
		t.Fatalf("error loading list values: %d != %d", len(result), len(expected))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("error in list value: %q != %q", result[i], expected[i])
		}
	}
}
