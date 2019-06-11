# dotm - Dotfile Manager

[![Build Status](https://travis-ci.org/relnod/dotm.svg?branch=master)](https://travis-ci.org/relnod/dotm)
[![Godoc](https://godoc.org/github.com/relnod/dotm?status.svg)](https://godoc.org/github.com/relnod/dotm)

dotm is a dotfile manager. It works by symlinking the dotfiles from multiple
profiles to the `$HOME` directory. It expects the dotfile profile to be under
source controll by git. This makes it easy to share dotfiles.

## Installation

### Using installer

The `install.sh` script will automatically install the latest binary from the
github release artifacts.

```shell
$ sh <(curl -s https://raw.githubusercontent.com/relnod/dotm/master/install.sh)
```

By adding the `--user` flag the install directory will be `$HOME/.local/bin`
instead of `/usr/local/bin`.

### Using Go

If you have a working go environment, you can simply install it via `go get`.

```shell
$ go get github.com/relnod/dotm/cmd/dotm
```

#### Completion

```shell
# Bash completions can be generated with the following command
$ dotm --genCompletions > /etc/bash_completion.d/dotm

# A cd helper can be generated with the following command
$ dotm --genChangeDirectory >> ~/.bashrc

$ dcd myprofile
```

## Usage

`dotm` works by symlinking the files from a dotfile folder to its
corresponding place under the `$HOME` directory of the user.

```shell
$ dotm
Usage:
  dotm [flags]
  dotm [command]

Available Commands:
  add         Add a new/existing file to the profile
  config      Sets/Gets values from the configuration file.
  fix         Tries to fix the configuration file
  help        Help about any command
  init        Initialize a new dotfile profile from the given path.
  install     Install dotfiles from a remote git repository
  list        Lists all profiles
  new         Create a new dotfile profile
  uninstall   Uninstall the profile
  update      Updates the symlinks for a given profile.

Flags:
      --genChangeDirectory   generate bash change directory command
      --genCompletions       generate bash completions
  -h, --help                 help for dotm
      --version              version for dotm

Use "dotm [command] --help" for more information about a command.
```

**Examples**:

```shell
# New (empty) dotfile profile
$ dotm new myprofile

# New profile from an existing dotfile folder
$ dotm init <path-to-existing-dotfile-folder>
$ dotm init --profile=myprofile <path-to-existing-dotfile-folder>

# New profile from a remote git repository
$ dotm install <url-to-remote-repository>
$ dotm install --profile=myprofile <url-to-remote-repository>

# Updating a profile
$ dotm update myprofile
$ dotm update myprofile --fromRemote
$ dotm update myprofile --no-hooks

# Get/Set a configuration value
$ dotm config ignore_prefix "_"
$ dotm config profile.default.path
$ dotm config profile.default.path "/my/path"
```

### Configuration file

The configuration file is located at `$HOME/.config/dotm/config.toml`. It can
hold multiple profiles. Each profile consists of a path to the local dotfile
location and an optional path to a remote git repository. It is also possible to
specify top level directories, that get included/excluded or add pre/post update
hooks.

**Example**:

```toml
# $HOME/.config/dotm/config.toml

# You can define multiple profiles
[profiles.default]

# Upstream git repository
remote = "github.com/relnod/dotm"

# Path to the local dotfile folder
path = "$HOME/.config/dotm/profiles/default"

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

# Pre update hooks
pre_update = [
    "echo 'pre update'"
]

# Post update hooks
post_update = [
    "echo 'post update'"
]

# Map of variables used for template processing
[profiles.default.vars]
foo = "bar"
bar = "foo"
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

### Template files

Template files can be used to dynamically add user identifying information inside
configuration files. All files with a `.tpl` file ending are treated as template
files. Templates support the syntax from the go
[text/template](https://golang.org/pkg/text/template/) package. Variables
can be configured per profile. When a template gets processed, a temporary file
with the same name plus a `.out` ending gets generated. This file will be
symlinked to the destination file without the `.tpl` suffix. Make sure to add
`*.tpl.out` to your `.gitignore` when using templates to prevent adding those
to git.

**Example**:

```text
# $HOME/.config/dotm/profiles/myprofile/git/.gitconfig.tpl

[user]
    name = {{ .GitUser }}
    email = {{ .GitEmail }}
```

### Hooks

Update hooks can be applied via global config, at profile root or per top level
directory. For hooks at profile root and top level directory you can create a
hooks.toml. Note: This file won't be symlinked.

**Example**:

```toml
# $HOME/.config/dotm/profiles/myprofile/nvim/hooks.toml
pre_update = [
    "nvim +PlugInstall +qall"
]
```

## Breaking Changes

Although dotm is considered somewhat stable, some breaking changes are expected
until a 1.0 release. When a breaking change is introduced try to run the fix
command. This tries to restore the original behaviour by modifying the
configuration file.

```shell
# Restore old behaviour by modifying the configuration file
$ dotm fix
```

## Development

There is a Makefile at the repository root to help with development. Start by
running `make` to execute tests and check for lint failures.

Tests are mostly wrtten using the [go internal testscript](https://github.com/rogpeppe/go-internal/tree/master/testscript)
package. The testscripts are located at `cmd/dotm/testdata`. For more
information on how to use `testscript` see this [README](https://github.com/golang/go/blob/master/src/cmd/go/testdata/script/README)
located in the go repository.

## License

dotm is licensed under the MIT License. See the [LICENSE](./LICENSE) details.
