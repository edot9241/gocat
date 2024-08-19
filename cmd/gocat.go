package gocat

import (
	"os"

	"github.com/edot9241/gocat/v1/internal"
)

/* Parses command line arguments and executes gocat logic on the input. */
func Run(args []string) {
	config := internal.PrepareConfig(args)

	if config.ShowError {
		internal.PrintError(config.Err)
		return
	}

	if config.ShowHelp {
		internal.PrintHelp()
		return
	}

	if config.ShowVersion {
		internal.PrintVersion()
		return
	}

	// TODO: if file / stdin

	file, err := os.Open(config.Filepath)
	if err != nil {
		internal.PrintError(err.Error(), "\"", config.Filepath, "\"")
		return
	}
	defer file.Close()

	internal.ProcessFile(file, &config)
}
