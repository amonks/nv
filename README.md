# nv

`nv` is a CLI note-taking application inspired by [Notational Velocity](http://notational.net) and [nvALT](http://brettterpstra.com/projects/nvalt/). It just shells out to other programs.

![gif](https://raw.githubusercontent.com/amonks/nv/master/demo.gif)

### features

- bring your own editor (vim, emacs, nano, kak, \*)
- bring your own file finder (fd, git-ls-files, ls, find)
- automatic sync over git

### options

Run `text --help` to get this list of options:

```
$ nv -help
Usage of nv:
  -directory string
        directory of text files
  -editor string
        editor command (default "vim +")
  -exit
        exit after editing one file
  -find-command string
        how to get a list of the files in a directory (default "ls")
  -frontmatter
        create yaml frontmatter for new files
  -sync
        automatically pull, commit, and push
```

### install

```bash
brew install fzf go

go get github.com/amonks/nv
go install github.com/amonks/nv
```

### use

I use nv via an alias:

```bash
alias t="nv --find-command 'fd --type f' --directory ~/txt --editor 'vim +' --frontmatter --sync"
```
