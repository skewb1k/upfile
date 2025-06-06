# UpFile

> ⚠️ **Work in Progress**: This project is under active development.
> The CLI is not stable and may change at any time until v1.0.0 release.
> If you encounter any issues on beta versions, try to delete $UPFILE_DIR directory, its layout may change

**UpFile** is a CLI tool for syncing files across multiple projects.

It's designed to help you manage shared configuration files like .prettierrc, .golangci.yml,
or any other files.

UpFile operates on a simple but powerful principles:

- Each file identifies by unique `filename`.
- Each filename has associated with it `upstream` version, which acts as the source of truth.
- You can add multiple instances of the same file across your projects. These are called `entries`.
- Entry can be `pushed` to the upstream and be `pulled` from it.

# Installation

```bash
go install github.com/skewb1k/upfile/cmd/upfile@latest
```

# Basic Usage

Suppose you have two projects: project-a and project-b, and both have the same config file .prettierrc.

```
~/project-a/.prettierrc
~/project-b/.prettierrc
```

And you want to keep them in sync, edit one file and easily spread changes to other places in other projects.

## Add files to tracking

Use `add` to register each file:

```bash
$ upfile add ~/project-a/.prettierrc
$ upfile add ~/project-b/.prettierrc
```

Check tracked state:

```bash
$ upfile ls
test.txt:
  /home/user/project-a/.prettierrc  Up-to-date
  /home/user/project-b/.prettierrc  Up-to-date
```

## Make a change and push it

Make a change in one of tracked files and push to the upstream.

```bash
$ echo 'change' >> ~/project-a/.prettierrc
$ upfile push ~/project-a/.prettierrc
```

List again to see status:

```bash
$ upfile ls
test.txt:
  /home/user/project-a/.prettierrc  Up-to-date
  /home/user/project-b/.prettierrc  Modified
```

## Sync other copies

```bash
$ upfile sync test.txt
The following tracked files will be updated:
 - /home/user/project-b/.prettierrc

Proceed? [Y/n]: y
```

Now both files are consistent again.

# Environment Variables

- `$UPFILE_DIR` - Path to the directory where UpFile stores metadata and upstream file versions. By default `$XDG_DATA_HOME/upfile`

# Shell Completion

UpFile can generate Bash, fish, PowerShell, and Zsh completion files.

See the instructions on `upfile completion <YOUR_SHELL> -h`
