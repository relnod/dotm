# dotm - Dotfile Manager

[![CircleCI](https://circleci.com/gh/relnod/dotm.svg?style=svg)](https://circleci.com/gh/relnod/dotm)
[![Godoc](https://godoc.org/github.com/relnod/dotm?status.svg)](https://godoc.org/github.com/relnod/dotm)

This Project is still WIP!

## Installation

Currently the only supported installation method is via the `go` command.
```
go install github.com/relnod/dotm/cmd/dotm
```

## Usage

### From an existing dotfile folder

```
dotm init <path-to-existing-dotfile-folder>
```

### From a remote git repository

```
dotm get <url-to-remote-repository>
```

## Development

There is a Makefile to help development. You can run `make watch` to start a file watcher, that runs tests on file change plus some other helpfull stuff. Type `make help` for a list of all commands.
