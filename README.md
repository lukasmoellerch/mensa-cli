# mensa-cli
A CLI which fetches the current meals available at ETHZ / UZH canteens

## Installation

The CLI tool can be install using `go install`:
```
go install github.com/lukasmoellerch/mensa-cli
```

## Usage

`mensa-cli` supports the concept of groups: a group is a collection of canteens. The default which the `mensa-cli meals` command chooses is the group `default`. You can add groups by using the `mensa-cli group add` command, it will open an editor where you can select a list of canteens to add to the group.

The main command `mensa-cli meals` fetches the list of meals for a given date. It can be used in two ways:
- `mensa-cli meals --group [group]`: fetches the list of meals for the given group
- `mensa-cli meals --filter [filter]`: fetches the list of meals from the canteens which match the given filter.

The daytime (dinner / lunch) is chosen automatically, but can be overridden by using the `--daytime` flag or using the shorthand versions `--dinner` or `--lunch`.

Check out the documentation located at [docs/mensa-cli.md](docs/mensa-cli.md) for more information.