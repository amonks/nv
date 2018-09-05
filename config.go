package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	Directory      string
	Editor         string
	UseFrontmatter bool
	ExitAfter      bool
	FdCmd          string
	Sync           bool
}

func nonEmptyStringOr(s, alt string) string {
	if s != "" {
		return s
	}
	return alt
}

var config Config

func init() {
	editorF := flag.String("editor", nonEmptyStringOr(os.Getenv("EDITOR"), "vim +"), "editor command")
	directoryF := flag.String("directory", "", "directory of text files")
	exitAfterF := flag.Bool("exit", false, "exit after editing one file")
	fdCmdF := flag.String("find-command", "ls", "how to get a list of the files in a directory")
	syncF := flag.Bool("sync", false, "automatically pull, commit, and push")
	useFrontmatterF := flag.Bool("frontmatter", false, "create yaml frontmatter for new files")
	flag.Parse()
	config = Config{
		Directory:      *directoryF,
		Editor:         *editorF,
		ExitAfter:      *exitAfterF,
		FdCmd:          *fdCmdF,
		Sync:           *syncF,
		UseFrontmatter: *useFrontmatterF,
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if config.Directory != "" {
		if config.Directory == "~" {
			config.Directory = usr.HomeDir
		} else if strings.HasPrefix(config.Directory, "~/") {
			config.Directory = filepath.Join(usr.HomeDir, config.Directory[2:])
		}

		if _, err := os.Stat(config.Directory); os.IsNotExist(err) {
			log.Printf("Creating config directory %s\n", config.Directory)
			os.MkdirAll(config.Directory, os.ModePerm)
		}
	}
}

func getConfig() Config {
	return config
}
