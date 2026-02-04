package cmd

import (
	"fmt"
	"os"
)

type CompletionCmd struct {
	Bash BashCompletionCmd `cmd:"" help:"Generate bash completion script"`
	Zsh  ZshCompletionCmd  `cmd:"" help:"Generate zsh completion script"`
	Fish FishCompletionCmd `cmd:"" help:"Generate fish completion script"`
}

type BashCompletionCmd struct{}
type ZshCompletionCmd struct{}
type FishCompletionCmd struct{}

func (c *BashCompletionCmd) Run() error {
	fmt.Fprintln(os.Stdout, bashCompletion)

	return nil
}

func (c *ZshCompletionCmd) Run() error {
	fmt.Fprintln(os.Stdout, zshCompletion)

	return nil
}

func (c *FishCompletionCmd) Run() error {
	fmt.Fprintln(os.Stdout, fishCompletion)

	return nil
}

const bashCompletion = `# irail bash completion
_irail() {
    local cur prev words cword
    _init_completion || return

    local commands="version stations liveboard connections vehicle composition disturbances completion"

    if [[ ${cword} -eq 1 ]]; then
        COMPREPLY=($(compgen -W "${commands}" -- "${cur}"))
        return
    fi

    case ${words[1]} in
        completion)
            COMPREPLY=($(compgen -W "bash zsh fish" -- "${cur}"))
            ;;
    esac
}

complete -F _irail irail`

const zshCompletion = `#compdef irail

_irail() {
    local -a commands
    commands=(
        'version:Print version information'
        'stations:List or search stations'
        'liveboard:Show departures or arrivals for a station'
        'connections:Find connections between stations'
        'vehicle:Show vehicle/train information'
        'composition:Show train composition'
        'disturbances:Show service disturbances'
        'completion:Generate shell completions'
    )

    _arguments -C \
        '--json[Output JSON to stdout]' \
        '--lang[Language (nl, fr, en, de)]:language:(nl fr en de)' \
        '--no-color[Disable colors]' \
        '--help[Show help]' \
        '1: :->cmds' \
        '*::arg:->args'

    case "$state" in
        cmds)
            _describe -t commands 'command' commands
            ;;
        args)
            case $words[1] in
                completion)
                    _values 'shell' bash zsh fish
                    ;;
                liveboard)
                    _arguments \
                        '1:station:' \
                        '--arrivals[Show arrivals instead of departures]' \
                        '--time[Time (HH:MM)]:time:' \
                        '--date[Date (YYYY-MM-DD)]:date:'
                    ;;
                connections)
                    _arguments \
                        '1:from station:' \
                        '2:to station:' \
                        '--time[Departure time (HH:MM)]:time:' \
                        '--date[Date (YYYY-MM-DD)]:date:' \
                        '--arrive-by[Time is arrival time]' \
                        '--results[Number of results]:count:'
                    ;;
            esac
            ;;
    esac
}

_irail`

const fishCompletion = `# irail fish completion

complete -c irail -n "__fish_use_subcommand" -a version -d "Print version information"
complete -c irail -n "__fish_use_subcommand" -a stations -d "List or search stations"
complete -c irail -n "__fish_use_subcommand" -a liveboard -d "Show departures or arrivals"
complete -c irail -n "__fish_use_subcommand" -a connections -d "Find connections"
complete -c irail -n "__fish_use_subcommand" -a vehicle -d "Show vehicle information"
complete -c irail -n "__fish_use_subcommand" -a composition -d "Show train composition"
complete -c irail -n "__fish_use_subcommand" -a disturbances -d "Show disturbances"
complete -c irail -n "__fish_use_subcommand" -a completion -d "Generate completions"

complete -c irail -l json -d "Output JSON to stdout"
complete -c irail -l lang -d "Language" -xa "nl fr en de"
complete -c irail -l no-color -d "Disable colors"

complete -c irail -n "__fish_seen_subcommand_from completion" -a "bash zsh fish"
complete -c irail -n "__fish_seen_subcommand_from liveboard" -s a -l arrivals -d "Show arrivals"
complete -c irail -n "__fish_seen_subcommand_from liveboard" -s t -l time -d "Time (HH:MM)"
complete -c irail -n "__fish_seen_subcommand_from liveboard" -s d -l date -d "Date (YYYY-MM-DD)"`
