package templates

import (
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig"
)

// Renderer provides the ability to write a template using the
// environment as input.
type Renderer struct {
	Template string `json:"template"`
	Target   string `json:"target"`
	Owner    string `json:"owner"`
	Group    string `json:"group"`
	FileMode string `json:"mode"`
	Env      map[string]string
}

// Execute renders the template to the specified target.
func (conf *Renderer) Execute() error {
	fh, err := os.Create(conf.Target)
	if err != nil {
		return err
	}

	defer fh.Close()

	err = ApplyTemplate(conf.Template, fh, conf.Env)
	if err != nil {
		return err
	}

	err = conf.SetPermissions()
	if err != nil {
		return err
	}

	return nil
}

// ApplyTemplate will takea template and write the output to the
// provided io.Writer adding the sprig helpers and using the provided env
// for data.
func ApplyTemplate(t string, fh io.Writer, env map[string]string) error {
	// Find the abs path of the template
	t, err := filepath.Abs(t)
	if err != nil {
		return err
	}
	name := filepath.Base(t)

	// Parse the template adding our sprig helpers
	tmpl, err := template.New(name).Funcs(sprig.TxtFuncMap()).ParseFiles(t)
	if err != nil {
		return err
	}

	// Execute the template with our environment
	err = tmpl.Execute(fh, env)
	if err != nil {
		return err
	}
	return nil
}

// SetPermissions ensures the user, group and file mode are set on the target file.
func (conf *Renderer) SetPermissions() error {

	var uid int
	var gid int

	// Set some defaults based on the current user
	u, err := user.Current()
	if err != nil {
		return err
	}

	g, err := user.LookupGroupId(u.Gid)
	if err != nil {
		return err
	}

	// If we have an owner, update the user var
	if conf.Owner != "" {
		u, err = user.Lookup(conf.Owner)
		if err != nil {
			return err
		}
	}

	// if we have a group, update the group var
	if conf.Group != "" {
		g, err = user.LookupGroup(conf.Group)
		if err != nil {
			return nil
		}
	}

	// set the user and group id vars
	uid, err = strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	gid, err = strconv.Atoi(g.Gid)
	if err != nil {
		return nil
	}

	// chown the file
	err = os.Chown(conf.Target, uid, gid)
	if err != nil {
		return err
	}

	if conf.FileMode != "" {
		mode, err := strconv.ParseUint(conf.FileMode, 0, 32)
		if err != nil {
			return err
		}

		os.Chmod(conf.Target, os.FileMode(mode))
	}

	return nil
}
