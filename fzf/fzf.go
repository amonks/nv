package fzf

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func Fzf(list string) (string, error) {
	fzfPath, err := exec.LookPath("fzf")
	if err != nil {
		return "", errors.Wrapf(err, "Error while finding fzf")
	}

	cmd := exec.Command(fzfPath, "--print-query")

	var stdoutBuf, stderrBuf bytes.Buffer

	stdoutCloser, _ := cmd.StdoutPipe()
	stderrCloser, _ := cmd.StderrPipe()

	cmd.Stdin = io.MultiReader(
		strings.NewReader(list),
	)

	var errStdout, errStderr error

	stdout := io.MultiWriter(&stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	if err := cmd.Start(); err != nil {
		return "", errors.Wrapf(err, "Error starting fzf")
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutCloser)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrCloser)
	}()

	foundExistingFile := true
	if err := cmd.Wait(); err != nil {
		foundExistingFile = false
	}

	if errStdout != nil || errStderr != nil {
		return "", errors.Errorf("Failed to capture stdout or stderr")
	}

	outStr := string(stdoutBuf.Bytes())
	lines := strings.Split(outStr, "\n")

	if foundExistingFile {
		return lines[1], nil
	}
	return lines[0], nil
}
