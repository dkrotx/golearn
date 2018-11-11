package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	err := cli.NewApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
