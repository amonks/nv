package git

import (
	"github.com/amonks/nv/shellout"
	"github.com/pkg/errors"
)

func Add(repo, path string) (string, error) {
	out, err := shellout.RunIn(repo, "git add %s", path)
	return out, errors.Wrapf(err, "Error adding files")
}

func Commit(repo, message string) (string, error) {
	out, err := shellout.RunIn(repo, "git commit -m '%s'", message)
	return out, errors.Wrapf(err, "Error committing")
}

func Pull(repo string) (string, error) {
	out, err := shellout.RunIn(repo, "git pull")
	return out, errors.Wrapf(err, "Error Pulling")
}

func Push(repo string) (string, error) {
	out, err := shellout.RunIn(repo, "git push")
	return out, errors.Wrapf(err, "Error pushing")
}
