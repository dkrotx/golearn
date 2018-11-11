package main

/* Features:
   - Flags (including Hidden option)
   - Shorten flags (see "lang" option)
   - FilePath (see "password")
   - Environment (see "color")

   - Use cli.NewExitError instead of manual exit and printing error
     (see mainImpl)
*/

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func coloredValue(s string, colored bool) string {
	if colored {
		return color.RedString(s)
	}
	return s
}

func mainImpl(c *cli.Context) error {
	name := "Jan"

	obsoleteName := c.String("name")
	if obsoleteName != "" {
		name = obsoleteName
	}

	if c.NArg() > 0 {
		name = c.Args().Get(0)
	}

	helloString := map[string]string{
		"spanish": "Hola",
		"english": "Hello",
	}

	lang := strings.ToLower(c.String("lang"))
	greeting, found := helloString[lang]
	if !found {
		return cli.NewExitError("Unknown language: "+lang, 1)
	}

	newLine := "\n"
	if c.Bool("n") {
		newLine = ""
	}

	fmt.Printf("%s, %s (password: %s) %s",
		greeting,
		coloredValue(name, c.Bool("color")),
		c.String("password"),
		newLine)
	return nil
}

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang, l",
			Value: "english",
			Usage: "language for the greeting (also try spanish)",
		},
		cli.StringFlag{
			Name:   "name",
			Hidden: true,
			Usage:  "name flag (obsolete, now using positional argument)",
		},
		cli.StringFlag{
			Name: "password, p",
			// it will take value from this file unless -p|--password given
			FilePath: "data/password",
			Usage:    "your secret password",
		},
		cli.BoolFlag{
			Name:  "n",
			Usage: "do not print the trailing newline character",
		},
		cli.BoolFlag{
			Name:   "color, c",
			Usage:  "use color for greeting",
			EnvVar: "COLOR_OUTPUT", // launch with COLOR_OUTPUT=1 to try
		},
	}

	app.Action = mainImpl
	app.Run(os.Args)
}
