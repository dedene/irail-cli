# 🚃 irail-cli - iRail (NMBS/SNCB) in your terminal

[![CI](https://github.com/dedene/irail-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/dedene/irail-cli/actions/workflows/ci.yml)
[![Go 1.23+](https://img.shields.io/badge/go-1.23+-00ADD8.svg)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/dedene/irail-cli)](https://goreportcard.com/report/github.com/dedene/irail-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

CLI for Belgian railway (NMBS/SNCB) real-time data via [iRail API](https://api.irail.be/).

## Installation

### Homebrew (macOS/Linux)

```bash
brew install dedene/tap/irail
```

### Go Install

```bash
go install github.com/dedene/irail-cli/cmd/irail@latest
```

### Binary

Download from [Releases](https://github.com/dedene/irail-cli/releases).

## Usage

### Liveboard

Show departures from a station:

```bash
irail liveboard Brugge
irail liveboard "Brussel-Centraal"
irail liveboard Brugge --arrivals
irail liveboard Brugge --time 09:00 --date 2025-02-15
```

### Connections

Find routes between stations:

```bash
irail connections Brugge Leuven
irail connections Brugge Leuven --time 09:00
irail connections Brugge Leuven --arrive-by  # time is arrival time
irail connections Brugge Leuven --results 10
```

### Stations

List or search stations:

```bash
irail stations
irail stations --search bruss
```

### Vehicle

Show train information:

```bash
irail vehicle IC1832
irail vehicle IC1832 --stops  # show all stops
```

### Composition

Show train composition (seats, amenities):

```bash
irail composition S51507
```

### Disturbances

Show service disruptions:

```bash
irail disturbances
irail disturbances --type planned    # only planned works
irail disturbances --type disturbance  # only disruptions
```

## Options

| Flag         | Description              |
| ------------ | ------------------------ |
| `--json`     | Output JSON              |
| `--lang`     | Language: nl, fr, en, de |
| `--no-color` | Disable colors           |

### Environment Variables

- `IRAIL_LANG` - Default language
- `IRAIL_JSON` - Default to JSON output
- `NO_COLOR` - Disable colors

## Features

- Real-time departures and arrivals
- Connection planning with transfers
- Delay highlighting (red for delays)
- Platform change warnings (yellow ⚠️)
- Occupancy indicators
- Clickable train links (in supported terminals)
- Shell completions (bash, zsh, fish)

## Shell Completions

```bash
# Bash
irail completion bash > /etc/bash_completion.d/irail

# Zsh
irail completion zsh > "${fpath[1]}/_irail"

# Fish
irail completion fish > ~/.config/fish/completions/irail.fish
```

## Agent Skill

This CLI is available as an [open agent skill](https://skills.sh/) for AI assistants including [Claude Code](https://claude.ai/code), [OpenClaw](https://openclaw.ai/), Cursor, and GitHub Copilot:

```bash
npx skills add dedene/irail-cli
```

## License

MIT
