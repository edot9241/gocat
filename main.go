package main

import (
	"os"

	gocat "github.com/edot9241/gocat/v1/cmd"
)

func main() {
	gocat.Run(os.Args, os.Stdout)
}
