package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/amonks/nv/editor"
	"github.com/amonks/nv/fd"
	"github.com/amonks/nv/fzf"
	"github.com/amonks/nv/git"
	"github.com/pkg/errors"
)

func makeFrontmatter(title string) string {
	return strings.Join([]string{
		fmt.Sprintf("created-at: %s", time.Now().Format("2006-01-02 15:04")),
		fmt.Sprintf("title: %s", title),
		"---",
	}, "\n")
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	cfg := getConfig()
	for {
		if cfg.Sync {
			if _, err := git.Pull(cfg.Directory); !strings.Contains(err.Error(), "No remote repository specified.") {
				exitOnErr(errors.Wrapf(err, "Error pulling"))
			}
		}

		files, err := fd.Fd(cfg.FdCmd, cfg.Directory)
		exitOnErr(errors.Wrap(err, "Error listing directories"))

		filename, err := fzf.Fzf(files)
		exitOnErr(errors.Wrap(err, "Error finding path"))
		if filename == "" {
			os.Exit(1)
		}

		path := path.Join(cfg.Directory, filename)

		template := ""
		if cfg.UseFrontmatter {
			template = makeFrontmatter(filename) + "\n\n"
		}

		_, err = editor.CreateOrEditFile(cfg.Editor, path, template)
		exitOnErr(errors.Wrap(err, "Error editing file"))

		if cfg.Sync {
			_, err = git.Add(cfg.Directory, path)
			exitOnErr(errors.Wrapf(err, "Error adding files"))

			if _, err := git.Commit(cfg.Directory, "text edit "+path); !strings.Contains(err.Error(), "nothing to commit") {
				exitOnErr(errors.Wrapf(err, "Error committing"))
			}

			if _, err = git.Push(cfg.Directory); !strings.Contains(err.Error(), "No configured push destination.") {
				exitOnErr(errors.Wrapf(err, "Error pushing"))
			}
		}

		if cfg.ExitAfter {
			os.Exit(1)
		}
	}
}
