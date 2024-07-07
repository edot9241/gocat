package main

import "fmt"

const (
	COLOR_RESET = "\033[0m"
	COLOR_ERROR = "\033[31m"
)

const version = "1"

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
	Copyright © 2024 Free Software Foundation, Inc.  License GPLv3+:
	GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.
	This is free software: you are free to change and redistribute
	it.  There is NO WARRANTY, to the extent permitted by law.
`

func PrintVersion() {
	fmt.Print(version)
}

func PrintHelp() {
	fmt.Print(helpString)
}

func PrintError(strings ...string) {
	fmt.Println(COLOR_ERROR+"ERROR: ", strings, COLOR_RESET)
}
