# Deltarune Save Manager (`drsm`)

A cross-platform Command Line Interface (CLI) tool written in Go that lets you store, view, name, copy, and link an unlimited number of Deltarune saves. It is designed to be fully compatible with future game chapters across Windows, macOS, and Linux.

## Features

- Store an unlimited number of saves and chapter completions.
- Attach or detach saves from the game's official save slots.
- Create instant copies of your saves.
- Link multiple game slots to a single save file to sync progress.
- Switch saves seamlessly while the game is running.
- View the detailed internal contents of your save files.

## Installation

### Releases

Download executable for your OS and architecture from the latest release.

### Using Go

Ensure you have [Go](https://go.dev) installed, then run:

```bash
go install github.com/LammoGit/Deltarune-Save-Manager/cmd/drsm@latest
```

## Usage

Get a full list of commands and options directly from the CLI:

```bash
drsm --help
```

### Command Overview

```text
NAME:
   drsm.exe - A CLI application for managing Deltarune saves

USAGE:
   drsm.exe [global options] [command [command options]]

COMMANDS:
   save     Commands for managing save files
   slot     Commands for managing save slots
   saves    List all stored save files
   slots    List all active save slots
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  Show help
```

### Examples

#### Save and link a save file in a slot
```bash
drsm slot save Pacifist 3 0
```
*Creates and links a new save named `Pacifist` using the contents of the first save slot from Chapter 3.*

#### Link a specific save slot to a save
```bash
drsm slot set Pacifist 3 3
```
*Links the save named `Pacifist` for Chapter 3 to the first completion slot of that same chapter.*

#### Clear a save slot without removing a save
```bash
drsm slot unset 4 2
```
*Clears the third save slot of Chapter 4 while keeping its save intact.*

## How It Works


`drsm` acts as a management layer between your local Deltarune save directory and drsm-managed save files' folder.

1. Creates a new folder named `save-manager` inside the official Deltarune save directory to store your backups.
2. Hardlinks save files from the drsm-managed folder directly to active game slots, ensuring that data stays synchronized.

## Roadmap

- [ ] Optimize file parser by removing `reflection` to improve performance.
- [ ] Refactor and improve global flags to string conversion.
- [ ] Add save file editing capabilities (modify stats, items, and flags).
- [ ] Develop a Graphical User Interface (GUI) wrapper for non-CLI users.

## License

This project is licensed under the [MIT License](LICENSE).
