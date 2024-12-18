package internal

import (
	"bufio"
	"fmt"
	"io"
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

func ProcessLine(loopState *LoopState, config *Config) (text string, shouldBePrinted bool) {
	line := loopState.line

	if config.SqueezeBlank {
		if loopState.empty && loopState.emptyLines > 1 {
			return "", false
		}
	}

	if config.NumberNonBlank {
		if !loopState.empty {
			line = fmt.Sprintf("%6s", strconv.Itoa(loopState.lineNumberNonEmpty)) + "\t" + line
		}
	} else if config.Number {
		line = fmt.Sprintf("%6s", strconv.Itoa(loopState.lineNumber)) + "\t" + line
	}

	if config.ShowEnds {
		line += "$"
	}

	if config.ShowTabs {
		line = strings.ReplaceAll(line, "\t", "^I")
	}

	if config.ShowNonPrinting {
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

func ProcessInput(input *os.File, config *Config, output io.Writer) {
	scanner := bufio.NewScanner(input)

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

		text, shouldBePrinted := ProcessLine(&loopState, config)

		if shouldBePrinted {
			fmt.Fprintln(output, text)
		}
	}
}
