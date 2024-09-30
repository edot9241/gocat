package internal

import (
	"strings"
)

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

	Filepaths []string
}

const (
	FilepathStdin = "-"
)

func isCompoundOption(option string) bool {
	return option[0] == '-' && len(option) > 2 && option[1] != '-'
}

func decoupleCompound(compoundOption string) []string {
	decoupledOptions := make([]string, 0)
	for _, char := range compoundOption[1:] {
		decoupledOptions = append(decoupledOptions, string("-")+string(char))
	}
	return decoupledOptions
}

func prepareArgs(args []string) []string {
	newArgs := make([]string, 0)
	for _, arg := range args[1:] {
		if isCompoundOption(arg) {
			newArgs = append(newArgs, decoupleCompound(arg)...)
		} else {
			newArgs = append(newArgs, arg)
		}
	}
	return newArgs
}

func PrepareConfig(osArgs []string) Config {
	config := Config{
		Filepaths: make([]string, 0),
	}

	args := prepareArgs(osArgs)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		default:
			if strings.HasPrefix(arg, "--") {
				config.ShowError = true
				config.Err = "gocat: unrecognized option '" + arg + "'"
				return config
			}
			if strings.HasPrefix(arg, "-") {
				config.ShowError = true
				config.Err = "gocat: invalid option -- '" + arg + "'"
				return config
			}
			config.Filepaths = append(config.Filepaths, arg)
		case "-", "--":
			config.Filepaths = append(config.Filepaths, FilepathStdin)
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

	// if no files passed, use stdin
	if len(config.Filepaths) == 0 {
		config.Filepaths = append(config.Filepaths, FilepathStdin)
	}

	return config
}
