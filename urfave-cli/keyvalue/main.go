package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
)

const (
	_usageExitCode = 2
	_keyNotFoundExitCode = 3
)


func SetKeyValue(ctx *cli.Context) error {
	if ctx.NArg() != 2 {
		return cli.NewExitError("usage: key value", _usageExitCode)
	}

	kv, err := NewFileKeyValue(ctx.GlobalString("file"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if err = kv.Set(ctx.Args()[0], ctx.Args()[1]); err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}

func GetKeyValue(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return cli.NewExitError("usage: key", _usageExitCode)
	}

	kv, err := NewFileKeyValue(ctx.GlobalString("file"))
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	value, err := kv.Get(ctx.Args()[0])
	if err != nil && err == KeyNotFoundErr {
		return cli.NewExitError("key not found", _keyNotFoundExitCode)
	}

	fmt.Println(value)
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
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "don't show any output",
		},
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
		},
		{
			Name: "get",
			Usage: "show filue by given key",
			Action: GetKeyValue,
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