package shellout

import (
	"fmt"
	"os/exec"
	"strings"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

func Run(command string, formatArgs ...interface{}) (string, error) {
	return RunIn("", command, formatArgs...)
}

func RunIn(path string, rawCommand string, formatArgs ...interface{}) (string, error) {
	command := fmt.Sprintf(rawCommand, formatArgs...)
	words, err := shellquote.Split(command)
	if err != nil {
		return "", errors.Errorf("Error splitting command '%s'", command)
	}
	if len(words) == 0 {
		return "", errors.Errorf("No command given")
	}

	resolvedCmdPath, err := exec.LookPath(words[0])
	if err != nil {
		return "", errors.Wrapf(err, "Error while finding command '%s'", words[0])
	}

	cmd := exec.Command(resolvedCmdPath, words[1:]...)
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Errorf("Error running command (%s)\n\n%s", command, strings.Trim(string(output), "\n"))
	}
	return string(output), nil
}
