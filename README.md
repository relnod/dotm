# dotm - Dotfile Manager

[![Build Status](https://travis-ci.org/relnod/dotm.svg?branch=master)](https://travis-ci.org/relnod/dotm)
[![Godoc](https://godoc.org/github.com/relnod/dotm?status.svg)](https://godoc.org/github.com/relnod/dotm)

dotm is a dotfile manager. It works by symlink the dotfiles from multiple
profiles to the `$HOME` directory.

## Installation

### Using installer

The `install.sh` script will automatically install the latest binary from the
github release artifacts.

```shell
sh <(curl -s https://raw.githubusercontent.com/relnod/dotm/master/install.sh)
```

By adding the `--user` flag the install directory will be `$HOME/.local/bin`
instead of `/usr/local/bin`.

### Using Go

If you have a working go environment, you can simply install it via `go get`.

```shell
go get github.com/relnod/dotm/cmd/dotm
```

#### Completion

```shell
# Bash completions can be generated with the following command
dotm --genCompletions > /etc/bash_completion.d/dotm

# A cd helper can be generated with the following command
dotm --genChangeDirectory >> ~/.bashrc

dcd myprofile
```

## Usage

`dotm` works by symlinking the files from the dotfile folder to its
corresponding place under the home directory of the user.

### Configuration file

The configuration file is located at `$HOME/.config/dotm/config.toml`. It can
hold multiple profiles. Each profile consists of a path to the local dotfile
location and an optional remote path to a git repository.

**Example**:

```toml
# $HOME/.config/dotm/config.toml

# You can define multiple profiles
[profiles.default]

# Upstream git repository
remote = "github.com/relnod/dotm"

# Path to local git repository
path = ".dotfiles/default/"

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

A Dotfile folder consists of multiple top level directories to group similar
configuration files (e.g. "vim" or "tmux"). The file structure below those top
level directories are directly mapped to the `$HOME` directory.

**Example**:

```
tmux/.tmux.conf             -> $HOME/.tmux.conf
bash/.bashrc                -> $HOME/.bashrc
nvim/.config/nvim/init.vim  -> $HOME/.config/nvim/init.vim
```

### New (empty) dotfile profile

```shell
dotm new myprofile
```

### New profile from an existing dotfile folder

```shell
dotm init <path-to-existing-dotfile-folder>
dotm init --profile=myprofile <path-to-existing-dotfile-folder>
```

### New profile from a remote git repository

```shell
dotm install <url-to-remote-repository>
dotm install --profile=myprofile <url-to-remote-repository>
```

### Updating a profile

```shell
dotm update myprofile
dotm update myprofile --fromRemote
```

### Hooks

Update hooks can be applied via global config, at profile root and per top level
directory. For hooks at profile root and top level directory you can create a
hooks.toml. Note: This file won't be symlinked.

**Example**:

```toml
# $HOME/.config/dotm/profiles/myprofile/nvim/hooks.toml
pre_update = [
    "nvim +PlugInstall +qall"
]
```

## Development

There is a Makefile at the repository root to help with development. Start by
running `make` to execute tests and check for lint failures.
