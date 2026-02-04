package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

type RootFlags struct {
	JSON    bool   `help:"Output JSON to stdout"`
	Lang    string `help:"Language (nl, fr, en, de)" default:"" env:"IRAIL_LANG"`
	NoColor bool   `help:"Disable colors" env:"NO_COLOR"`
}

// DetectLang returns the language from IRAIL_LANG, LANG, or LC_ALL env vars.
func (f *RootFlags) DetectLang() string {
	if f.Lang != "" {
		return f.Lang
	}

	// Check LANG/LC_ALL for language
	for _, env := range []string{"LANG", "LC_ALL"} {
		if v := os.Getenv(env); v != "" {
			v = strings.ToLower(v)
			switch {
			case strings.HasPrefix(v, "nl"):
				return "nl"
			case strings.HasPrefix(v, "fr"):
				return "fr"
			case strings.HasPrefix(v, "de"):
				return "de"
			case strings.HasPrefix(v, "en"):
				return "en"
			}
		}
	}

	return "en"
}

type CLI struct {
	RootFlags `embed:""`

	Version      kong.VersionFlag `help:"Print version and exit"`
	VersionCmd   VersionCmd       `cmd:"" name:"version" help:"Print version information"`
	Stations     StationsCmd      `cmd:"" help:"List or search stations"`
	Liveboard    LiveboardCmd     `cmd:"" help:"Show departures or arrivals for a station"`
	Connections  ConnectionsCmd   `cmd:"" help:"Find connections between stations"`
	Vehicle      VehicleCmd       `cmd:"" help:"Show vehicle/train information"`
	Composition  CompositionCmd   `cmd:"" help:"Show train composition"`
	Disturbances DisturbancesCmd  `cmd:"" help:"Show service disturbances"`
	Completion   CompletionCmd    `cmd:"" help:"Generate shell completions"`
}

type exitPanic struct{ code int }

func Execute(args []string) (err error) {
	parser, err := newParser()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if ep, ok := r.(exitPanic); ok {
				if ep.code == 0 {
					err = nil

					return
				}

				err = &ExitError{Code: ep.code, Err: errors.New("exited")}

				return
			}

			panic(r)
		}
	}()

	if len(args) == 0 {
		args = []string{"--help"}
	}

	kctx, err := parser.Parse(args)
	if err != nil {
		parsedErr := wrapParseError(err)
		_, _ = fmt.Fprintln(os.Stderr, parsedErr)

		return parsedErr
	}

	err = kctx.Run()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)

		return err
	}

	return nil
}

func wrapParseError(err error) error {
	if err == nil {
		return nil
	}

	var parseErr *kong.ParseError
	if errors.As(err, &parseErr) {
		return &ExitError{Code: ExitInvalidArgs, Err: parseErr}
	}

	return err
}

func newParser() (*kong.Kong, error) {
	vars := kong.Vars{
		"version": VersionString(),
	}

	cli := &CLI{}
	parser, err := kong.New(
		cli,
		kong.Name("irail"),
		kong.Description("CLI for Belgian railway (NMBS/SNCB) schedules"),
		vars,
		kong.Writers(os.Stdout, os.Stderr),
		kong.Exit(func(code int) { panic(exitPanic{code: code}) }),
		kong.Bind(&cli.RootFlags),
		kong.Help(helpPrinter),
		kong.ConfigureHelp(helpOptions()),
	)
	if err != nil {
		return nil, err
	}

	return parser, nil
}
