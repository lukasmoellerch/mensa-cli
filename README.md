# mensa-cli
A CLI which fetches the current meals available at ETHZ mensas

## Installation

The CLI tool can be install using `go install`:
```
go install github.com/lukasmoellerch/mensa-cli
```

## Usage
```
A CLI tool which fetches the list of meals for a given date

Usage:
  mensa-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  get         Fetches the list of meals for a given date
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.mensa-cli.yaml)
  -h, --help            help for mensa-cli

Use "mensa-cli [command] --help" for more information about a command.

