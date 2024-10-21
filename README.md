# Crev CLI

`crev` is a command-line tool that allows you to easily bundle your codebase into a single file for use with LLMs.

This project is a fork of [vossenwout/crev](https://github.com/vossenwout/crev), with improvements to the command
syntax to make the `bundle` command more intuitive and customizable for various use cases. The `bundle` command follows
familiar patterns from other Linux command-line tools, taking a directory as its primary argument and supporting include
and exclude flags for granular control. It also supports globbing for include and exclude flags, using the excellent
[doublestar](https://github.com/bmatcuk/doublestar) library.


## Features

- Bundle your codebase into a single .txt file.
- Select files and directories to include or exclude.
- Cross-platform support (Linux, macOS, Windows).
- Customizable (ignore / include specific files, directories, etc).
- Written in Go.

## Installation / Documentation

For installation instructions and documentation, go to the [official docs](https://crevcli.com/docs).

## Important Commands

* **Bundle your codebase (saved locally as a .txt file)**:

   ```bash
   crev bundle
   ```

For full details on usage and configuration, visit the [official docs](https://crevcli.com/docs).


## Contributing

We welcome contributions!
