package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	config := PrepareConfig(os.Args)

	if config.err != "" {
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
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
