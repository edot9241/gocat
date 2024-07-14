package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

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

	file, err := os.Open(config.filepath)
	if err != nil {
		PrintError(err.Error())
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0
	emptyLines := 0
	for scanner.Scan() {
		line := scanner.Text()

		if config.numberNonBlank {
			if line != "" {
				line = string(i) + " " + line
				i++
			}
		} else if config.number {
			line = string(i) + " " + line
			i++
		}

		if config.showEnds {
			line += "$"
		}

		if config.squeezeBlank && line == "" {
			if line != "" {
				emptyLines = 0
			} else {
				emptyLines++
			}

			if emptyLines > 2 {
				continue
			}
		}

		if config.showTabs {
			line = strings.ReplaceAll(line, "\t", "^I")
		}

		if config.showNonPrinting {
			for _, r := range line {
				if unicode.IsGraphic(r) || r == '\t' {
					continue
				}

			}
		}

		fmt.Println(line)
	}
}
