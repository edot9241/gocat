package gocat

import (
	"io"
	"os"

	"github.com/edot9241/gocat/v1/internal"
)

// Parse command line arguments, execute gocat logic on the input
// and print to output.
func Run(args []string, output io.Writer) {
	config := internal.PrepareConfig(args)

	if config.ShowError {
		internal.PrintError(output, config.Err)
		return
	}

	if config.ShowHelp {
		internal.PrintHelp(output)
		return
	}

	if config.ShowVersion {
		internal.PrintVersion(output)
		return
	}

	file, err := os.Open(config.Filepath)
	if err != nil {
		internal.PrintError(output, err.Error(), "\"", config.Filepath, "\"")
		return
	}
	defer file.Close()

	internal.ProcessFile(file, &config, output)
}
