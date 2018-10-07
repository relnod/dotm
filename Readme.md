# dotm - Dotfile Manager

[![CircleCI](https://circleci.com/gh/relnod/dotm.svg?style=svg)](https://circleci.com/gh/relnod/dotm)
[![Godoc](https://godoc.org/github.com/relnod/dotm?status.svg)](https://godoc.org/github.com/relnod/dotm)

This Project is still WIP!

## Installation

### Using Go
If you have a working go environment, you can simply install it via `go get`.
```
go get github.com/relnod/dotm/cmd/dotm
```

### Docker
If you have a working docker environment, you can run dotm with the following alias:

```
alias dotm="docker run -v /home/$USER:/home/$USER --env USER=$USER reldod/dotm:latest"
```
NOTE: Hooks might not work, if they require additional programs. If you want to call extra programs, you can create a new container based on relnod/dotm:latest and add those.

## Usage

`dotm` works by symlinking the files from the dotfile folder to its corresponding place under the home directory of the user.

### Configuration file
The configuration file is located at `$HOME/.dotfiles.toml` (can be changed with the --config flag). It can hold multiple profiles. Each profile consists of a path to the local dotfile location and an optional remote path to a git repository.


Example:
```toml
# You can define multiple profiles
[profiles.default]

# Upstream git repository
remote = "github.com/relnod/dotm"

# Path to local git repository
path = ".dotfiles/"

# Configs to be included
includes = [
    "bash",
    "nvim",
    "tmux"
]

# Configs to be excluded
excludes = [
    "bash",
    "nvim",
    "tmux"
]

pre_update = [
    "echo 'pre update'"
]

post_update = [
    "echo 'post update'"
]
```

### Dotfiles folder
A Dotfile folder consists of multiple top level directories to group similar configuration files (e.g. "vim" or "tmux"). The file structure below those top level directories are directly mapped to `$HOME`.

Example:
```
tmux/.tmux.conf             -> $HOME/.tmux.conf
bash/.bashrc                -> $HOME/.bashrc
nvim/.config/nvim/init.vim  -> $HOME/.config/nvim/init.vim
```

### From an existing dotfile folder
```
dotm init <path-to-existing-dotfile-folder>
```

### From a remote git repository
```
dotm install <url-to-remote-repository>
```

### Hooks
Update hooks can be applied at profile level (via global config) and per top level directory with a hooks.toml in the top level directory. Note: this file won't be symlinked.

Example:
```toml
pre_update = [
    "nvim +PlugInstall +qall"
]
```

## Development

There is a Makefile to help development. You can run `make watch` to start a file watcher, that runs tests on file change plus some other helpfull stuff. Type `make help` for a list of all commands.
