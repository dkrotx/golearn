package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"os"
	"time"
)

const (
	_usageExitCode = 2
	_keyNotFoundExitCode = 3
)

func parseSetOptions(ctx *cli.Context) ([]Option, error) {
	var opts []Option

	ttl := ctx.String("ttl")

	if ttl != "" {
		duration, err := time.ParseDuration(ttl)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse ttl option")
		}

		opts = append(opts, TTL(duration))
	}

	return opts, nil
}

func SetKeyValue(ctx *cli.Context) error {
	if ctx.NArg() != 2 {
		return cli.NewExitError("usage: key value", _usageExitCode)
	}
	key, value := ctx.Args()[0], ctx.Args()[1]

	kv, err := NewFileKeyValue(ctx.GlobalString("file"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	opts, err := parseSetOptions(ctx)
	if err != nil {
		return err
	}

	if err = kv.Set(key, value, opts...); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func GetKeyValue(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return cli.NewExitError("usage: key", _usageExitCode)
	}

	quiet := ctx.Bool("quiet")

	kv, err := NewFileKeyValue(ctx.GlobalString("file"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	value, err := kv.Get(ctx.Args()[0])
	if err != nil && err == KeyNotFoundErr {
		if quiet {
			return cli.NewExitError("", _keyNotFoundExitCode)
		}
		return cli.NewExitError("key not found", _keyNotFoundExitCode)
	}

	if  !quiet {
		fmt.Println(value)
	}

	return nil
}

func DeleteKeyValue(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return cli.NewExitError("usage: key", _usageExitCode)
	}

	kv, err := NewFileKeyValue(ctx.GlobalString("file"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = kv.Delete(ctx.Args()[0]); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "/tmp/golearn_kv.txt",
			Usage: "use given `FILE` as storage",
			EnvVar: "KV_FILE",
		},
	}

	app.Commands = []cli.Command {
		{
			Name: "set",
			Aliases: []string{"update"},
			Usage: "set key-value",
			Action: SetKeyValue,
			Flags: []cli.Flag {
				cli.StringFlag {
					Name: "ttl",
					Usage: "ttl duration (examples: 5s, 1m, 1m25s, ...)",
				},
			},
		},
		{
			Name: "get",
			Usage: "show filue by given key",
			Action: GetKeyValue,
			Flags: []cli.Flag {
				cli.BoolFlag{
					Name:  "quiet, q",
					Usage: "don't show any output. Only provide exit-code. (useful for scripts)",
				},
			},
		},
		{
			Name: "delete",
			Aliases: []string{"remove", "rm", "del"},
			Usage: "remove given key and it's value from storage",
			Action: DeleteKeyValue,
		},
	}

	app.Run(os.Args)
}