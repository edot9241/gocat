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

		// TODO: compound switches. E.g. vE
		switch arg {
		// Either a filepath or an unknown switch
		default:
			if arg[0] == '-' {
				config.showError = true
				config.err = fmt.Sprint("Unknown switch:", arg)
				return config
			} else {
				config.filepath = arg
			}
		case "-A", "--show-all":
			config.showNonPrinting = true
			config.showEnds = true
			config.showTabs = true
		case "-b", "--number-nonblank":
			config.numberNonBlank = true
		case "-e":
			config.showNonPrinting = true
			config.showEnds = true
		case "-E", "--show-ends":
			config.showEnds = true
		case "-n", "--number":
			config.number = true
		case "-s", "--squeeze-blank":
			config.squeezeBlank = true
		case "-t":
			config.showNonPrinting = true
			config.showTabs = true
		case "-T", "--show-tabs":
			config.showTabs = true
		// TODO: case "-u": ignored, but does it show something extra if you use it?
		case "-v", "--show-nonprinting":
			config.showNonPrinting = true
		case "--help":
			config.showHelp = true
			return config
		case "--version":
			config.showVersion = true
			return config
		}
	}

	if config.filepath == "" {
		config.err = fmt.Sprint("No filepath specified")
		return config
	}

	return config
}
