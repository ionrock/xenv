package templates_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/ionrock/xenv/templates"
)

func TestApplyTemplate(t *testing.T) {
	tmpl := "testdata/my.cfg.tmpl"

	env := map[string]string{
		"LISTEN":        "10.0.0.1:8900",
		"CLUSTER_HOSTS": "10.0.0.2 10.0.0.3 10.0.0.4",
	}

	var b bytes.Buffer

	err := templates.ApplyTemplate(tmpl, &b, env)
	if err != nil {
		t.Fatalf("failed to write template: %q", err)
	}

	contents := b.String()
	expected := `[service]
listen = 10.0.0.1:8900
workers =
   - 10.0.0.2
   - 10.0.0.3
   - 10.0.0.4
`
	if contents != expected {
		t.Errorf("wrong content: \n%q\n !=\n%q", contents, expected)
	}
}

type tmplPathTest struct {
	path   string
	tmpl   string
	target string
}

func TestConfigTmplFileInfo(t *testing.T) {
	conf := templates.Renderer{
		Target:   "testdata/foo.cfg",
		FileMode: "0645",
	}

	err := conf.SetPermissions()
	if err != nil {
		t.Fatalf("error settings permissions: %q", err)
	}

	info, err := os.Stat(conf.Target)
	if err != nil {
		t.Fatalf("failed to get fileinfo: %q", err)
	}

	mode := os.FileMode(int(0645))
	if info.Mode() != mode {
		t.Errorf("failed to set mode: %q != %q", info.Mode(), mode)
	}
}
