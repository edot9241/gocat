package internal

import "fmt"

// TODO: prevent redefinition of the fields?
type Config struct {
	// Display error and exit
	ShowError bool
	Err       string

	// Display help and exit
	ShowHelp bool
	// Output version information and exit
	ShowVersion bool

	// Number nonempty output lines, overrides 'number'
	NumberNonBlank bool
	// Display $ at end of each line
	ShowEnds bool
	// Number all output lines
	Number bool
	// Suppress repeated empty output lines
	SqueezeBlank bool
	// Display TAB characters as ^I
	ShowTabs bool
	// Use ^ and M- notation, except for LFD (newline?) and TAB
	ShowNonPrinting bool

	Filepath string
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
				config.ShowError = true
				config.Err = fmt.Sprint("Unknown switch:", arg)
				return config
			} else {
				config.Filepath = arg
			}
		case "-A", "--show-all":
			config.ShowNonPrinting = true
			config.ShowEnds = true
			config.ShowTabs = true
		case "-b", "--number-nonblank":
			config.NumberNonBlank = true
		case "-e":
			config.ShowNonPrinting = true
			config.ShowEnds = true
		case "-E", "--show-ends":
			config.ShowEnds = true
		case "-n", "--number":
			config.Number = true
		case "-s", "--squeeze-blank":
			config.SqueezeBlank = true
		case "-t":
			config.ShowNonPrinting = true
			config.ShowTabs = true
		case "-T", "--show-tabs":
			config.ShowTabs = true
		case "-v", "--show-nonprinting":
			config.ShowNonPrinting = true
		case "--help":
			config.ShowHelp = true
			return config
		case "--version":
			config.ShowVersion = true
			return config
		}
	}

	if config.Filepath == "" {
		config.Err = fmt.Sprint("No filepath specified")
		return config
	}

	return config
}
