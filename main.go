package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/* State of the current iteration of the gocat function. */
type LoopState struct {
	// Text value of the line
	line string
	// Is current line empty
	empty bool
	// Number of the current line
	lineNumber int
	// Number of the current line (if only counting non-empty lines)
	lineNumberNonEmpty int
	// Current number of empty lines in a row. TODO: cap at 2?
	emptyLines int
}

func transformLine(loopState *LoopState, config *Config) (text string, shouldBePrinted bool) {
	line := loopState.line

	if config.squeezeBlank {
		if loopState.empty && loopState.emptyLines > 2 {
			return "", false
		}
	}

	if config.numberNonBlank {
		if !loopState.empty {
			line = strconv.Itoa(loopState.lineNumberNonEmpty) + " " + line
		}
	} else if config.number {
		line = strconv.Itoa(loopState.lineNumber) + " " + line
	}

	if config.showEnds {
		line += "$"
	}

	if config.showTabs {
		line = strings.ReplaceAll(line, "\t", "^I")
	}

	if config.showNonPrinting {
		newLine := ""
		for _, r := range line {
			if unicode.IsGraphic(r) || r == '\t' {
				newLine += string(r)
			} else {
				newLine += "\\TODO"
			}
		}
		line = newLine
	}

	return line, true
}

func gocat(file *os.File, config *Config) {
	scanner := bufio.NewScanner(file)

	loopState := LoopState{}

	for scanner.Scan() {
		loopState.line = scanner.Text()

		loopState.empty = (loopState.line == "")

		loopState.lineNumber++

		if loopState.empty {
			loopState.emptyLines++
		} else {
			loopState.lineNumberNonEmpty++
			loopState.emptyLines = 0
		}

		text, shouldBePrinted := transformLine(&loopState, config)

		if shouldBePrinted {
			fmt.Println(text)
		}
	}
}

/* Parses command line arguments and executes gocat logic on the input. */
func main() {
	config := PrepareConfig(os.Args)

	if config.showError {
		PrintError(config.err)
		return
	}

	if config.showHelp {
		PrintHelp()
		return
	}

	if config.showVersion {
		PrintVersion()
		return
	}

	// TODO: if file / stdin

	file, err := os.Open(config.filepath)
	if err != nil {
		PrintError(err.Error())
		return
	}
	defer file.Close()

	gocat(file, &config)
}
