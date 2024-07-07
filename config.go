package main

import "fmt"

// TODO: prevent redefinition of the fields?
type Config struct {
	// Display error and exit
	showError bool
	err       string

	// Display help and exit
	showHelp bool
	// Output version information and exit
	showVersion bool

	// Number nonempty output lines, overrides 'number'
	numberNonBlank bool
	// Display $ at end of each line
	showEnds bool
	// Number all output lines
	number bool
	// Suppress repeated empty output lines
	squeezeBlank bool
	// Display TAB characters as ^I
	showTabs bool
	// Use ^ and M- notation, except for LFD and TAB
	showNonPrinting bool

	filepath string
}

func PrepareConfig(args []string) Config {
	config := Config{}

	for i := 1; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "--help":
			config.showHelp = true
			return config
		case "--version":
			config.showVersion = true
			return config
		default:
			if arg[0] == '-' {
				config.showError = true
				config.err = fmt.Sprint("Unknown switch:", arg)
				return config
			} else {
				config.filepath = arg
			}
		}
	}

	if config.filepath == "" {
		config.err = fmt.Sprint("No filepath specified")
		return config
	}

	return config
}
