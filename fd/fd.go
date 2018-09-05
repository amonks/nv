package fd

import (
	"fmt"
	"strings"

	"github.com/amonks/nv/shellout"
	"github.com/pkg/errors"
)

func Fd(cmd, path string) (string, error) {
	output, err := shellout.RunIn(path, cmd)
	if err != nil {
		return "", errors.Wrapf(err, "Error searching for files")
	}

	list := strings.Split(string(output), "\n")
	for i, str := range list {
		fmt.Println(str)
		list[i] = strings.Replace(str, path+"/", "", 1)
	}

	return strings.Join(list, "\n"), nil
}
