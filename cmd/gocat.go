package gocat

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/edot9241/gocat/v1/internal"
)

const versionString = `gocat 1.00
Copyright (C) E.#9241 <https://github.com/edot9241>.
License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

Written by E.#9241 <https://github.com/edot9241>`

const helpString = `NAME
       gocat - concatenate files and print on the standard output

SYNOPSIS
       gocat [OPTION]... [FILE]...

DESCRIPTION

       Concatenate FILE(s) to standard output.

       With no FILE, or when FILE is -, read standard input.

       -A, --show-all
              equivalent to -vET

       -b, --number-nonblank
              number nonempty output lines, overrides -n

       -e     equivalent to -vE

       -E, --show-ends
              display $ at end of each line

       -n, --number
              number all output lines

       -s, --squeeze-blank
              suppress repeated empty output lines

       -t     equivalent to -vT

       -T, --show-tabs
              display TAB characters as ^I

       -u     (ignored)

       -v, --show-nonprinting
              use ^ and M- notation, except for LFD and TAB

       --help display this help and exit

       --version
              output version information and exit
	
EXAMPLES
       gocat f - g
              Output f's contents, then standard input, then g's
              contents.

       gocat    Copy standard input to standard output.

AUTHOR
	Written by E.#9241 <https://github.com/edot9241>

REPORTING BUGS
	Github repository:
	<https://github.com/edot9241/gocat>

COPYRIGHT
	Copyright © 2024 E.#9241 <https://github.com/edot9241>.  License GPLv3+:
	GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.
	This is free software: you are free to change and redistribute
	it.  There is NO WARRANTY, to the extent permitted by law.
`

// Parse command line arguments, execute gocat logic on the input
// and print to output.
func Run(args []string, output io.Writer) {
	config := internal.PrepareConfig(args)

	if config.ShowError {
		fmt.Fprintf(output, config.Err)
		return
	}

	if config.ShowHelp {
		fmt.Fprint(output, helpString)
		return
	}

	if config.ShowVersion {
		fmt.Fprint(output, versionString)
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
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprint(output, "gocat: ", filepath, ": No such file or directory")
				continue
			} else {
				fmt.Fprintf(output, err.Error(), "\"", filepath, "\"")
				return
			}
		}

		internal.ProcessInput(file, &config, output)
		file.Close()
	}
}
