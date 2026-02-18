# UpFile

**UpFile** is a CLI tool for syncing files across multiple projects.

It's designed to help you manage shared configuration files.

UpFile operates on a simple but powerful principles:

- Each file identifies by unique `filename`.
- Each filename is associated with an `upstream` version that acts as the source of truth.
- You can add multiple instances of the same file across your projects. These
  are called `entries`.
- Entry can be `pushed` to the upstream and be `pulled` from it.

## Basic Usage

Suppose you have two projects: project-a and project-b, and both have the same
config file .prettierrc.

```
~/project-a/.prettierrc
~/project-b/.prettierrc
```

And you want to keep them in sync, edit one file and easily spread changes to
other places in other projects.

### Add files to tracking

Use `add` to track each file:

```bash
$ upfile add ~/project-a/.prettierrc
$ upfile add ~/project-b/.prettierrc
```

Check tracked status:

```bash
$ upfile ls
test.txt:
  /home/user/project-a/.prettierrc  Up-to-date
  /home/user/project-b/.prettierrc  Up-to-date
```

### Make a change and push it

Make a change in one entry and push it to the upstream.

```bash
$ echo 'change' >> ~/project-a/.prettierrc
$ upfile push ~/project-a/.prettierrc
```

List to see status:

```bash
$ upfile ls
test.txt:
  /home/user/project-a/.prettierrc  Up-to-date
  /home/user/project-b/.prettierrc  Modified
```

### Sync with other entries

```bash
$ upfile sync test.txt
The following tracked files will be updated:
 - /home/user/project-b/.prettierrc

Proceed? [Y/n]: y
```

Now both files are consistent again.

## Docs

For a complete list of commands and detailed usage instructions, run:

```bash
upfile -h
```

## Installation

### Package Manager

```bash
yay -S upfile-bin
```

```bash
brew install skewb1k/homebrew-tap/upfile
```

Support for other package managers is planned.

Or download a binary from the [releases](https://github.com/skewb1k/upfile/releases) page.

### Go

```bash
go install github.com/skewb1k/upfile
```

## Environment Variables

- `$UPFILE_DIR` - Path to the directory where UpFile stores metadata and
  upstream file versions. By default `$XDG_DATA_HOME/upfile`

## Shell Completion

UpFile supports advanced and interactive shell completion for Bash, Zsh, fish, and PowerShell.

See the instructions on `upfile completion <YOUR_SHELL> -h`
