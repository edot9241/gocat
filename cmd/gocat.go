package gocat

import (
	"fmt"
	"io"
	"os"

	"github.com/edot9241/gocat/v1/internal"
)

// Parse command line arguments, execute gocat logic on the input
// and print to output.
func Run(args []string, output io.Writer) {
	config := internal.PrepareConfig(args)

	if config.ShowError {
		fmt.Fprintf(output, config.Err)
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

	for _, filepath := range config.Filepaths {
		var file *os.File
		var err error

		if filepath == internal.FilepathStdin {
			file = os.Stdin
		} else {
			file, err = os.Open(filepath)
		}

		if err != nil {
			fmt.Fprintf(output, err.Error(), "\"", filepath, "\"")
			continue
		}

		internal.ProcessInput(file, &config, output)
		file.Close()
	}
}
