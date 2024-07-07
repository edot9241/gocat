package main

import "fmt"

// TODO: prevent redefinition of the fields?
type Config struct {
	// Display error and exit
	showError bool
	err       string

	showHelp    bool
	showVersion bool

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
