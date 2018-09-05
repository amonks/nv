package editor

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

func makeTempFile() (string, error) {
	tempFilePrefix := "tempFilePrefix"
	tmpDir := os.TempDir()

	tmpFile, err := ioutil.TempFile(tmpDir, tempFilePrefix)
	if err != nil {
		return "", errors.Wrapf(err, "Error creating tempFile")
	}

	path, err := filepath.Abs(filepath.Dir(tmpFile.Name()))
	if err != nil {
		return "", errors.Wrapf(err, "Error getting path to tmp file")
	}

	return path, nil
}

func editFile(editor, path string) error {
	editorWords, err := shellquote.Split(editor)
	if err != nil {
		return errors.Wrapf(err, "Error parsing editor %s", editor)
	}
	editorPath, err := exec.LookPath(editorWords[0])
	if err != nil {
		return errors.Wrapf(err, "Error %s while finding %s", err, editor)
	}

	words := append(editorWords[1:], path)
	cmd := exec.Command(editorPath, words...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "error starting external editor")
	}

	if err = cmd.Wait(); err != nil {
		return errors.Wrapf(err, "error with external editor")
	}

	return nil
}

func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "Error opening file")
	}

	info, err := file.Stat()
	if err != nil {
		return "", errors.Wrapf(err, "Error statting file")
	}

	bytes := make([]byte, info.Size())
	if _, err := file.Read(bytes); err != nil {
		return "", errors.Wrapf(err, "Error reading file")
	}

	return string(bytes), nil
}

func writeNewFile(path, contents string) error {
	if _, err := os.Stat(path); err == nil {
		return errors.Errorf("File '%s' exists", path)
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "Error opening file")
	}

	if _, err := file.WriteString(contents); err != nil {
		return errors.Wrapf(err, "Error writing to file")
	}

	return nil
}

func EditString(editor, str string) (string, error) {
	path, err := makeTempFile()
	if err != nil {
		return "", errors.Wrapf(err, "Error making temp file")
	}

	if err := writeNewFile(path, str); err != nil {
		return "", errors.Wrapf(err, "Error writing file")
	}

	if err := editFile(editor, path); err != nil {
		return "", errors.Wrapf(err, "Error editing file")
	}

	outStr, err := readFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "Error reading file")
	}

	return outStr, nil
}

func EditExistingFile(editor, path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", errors.Errorf("File '%s' does not exist", path)
	}

	if err := editFile(editor, path); err != nil {
		return "", errors.Wrapf(err, "Error editing file '%s'", path)
	}

	outStr, err := readFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "Error reading file '%s' after edit", path)
	}

	return outStr, nil
}

func CreateAndEditFile(editor, path string, initialContents string) (string, error) {
	if _, err := os.Stat(path); err == nil {
		return "", errors.Errorf("File '%s' exists", path)
	}

	if err := writeNewFile(path, initialContents); err != nil {
		return "", errors.Wrapf(err, "Error writing file '%s'", path)
	}

	return EditExistingFile(editor, path)
}

func CreateOrEditFile(editor, path, initialContentsIfCreate string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return CreateAndEditFile(editor, path, initialContentsIfCreate)
	}
	return EditExistingFile(editor, path)
}
