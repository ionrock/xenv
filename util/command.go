package util

import (
	"log"
	"os"
	"strings"

	shlex "github.com/flynn/go-shlex"
)

func SplitCommand(cmd string) []string {
	parts, err := shlex.Split(strings.Trim(os.ExpandEnv(cmd), "`"))
	if err != nil {
		log.Fatal(err)
	}
	return parts
}
