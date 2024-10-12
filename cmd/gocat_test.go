package gocat

// TODO: move into gocat_test

import (
	"bytes"
	"os"
	"path"
	"testing"
)

const (
	testString1 string = `
	hello,
	world 1!

       


	hey 1

	!@)(*$&)*!#%^%!_)+_




	` + "\n"

	testString2 string = `
	hello,
	world 2!

       


	hey 2

	!@)(*$&)*!#%^%!_)+_




	` + "\n"
)

func TestStdin(t *testing.T) {
	osArgs := [][]string{
		{"gocat.exe"},
		{"gocat.exe", "-"},
		{"gocat.exe", "--"},
	}

	expectedOutput := testString1

	for _, args := range osArgs {
		fakeStdinWrite(t, testString1)
		testOutput(t, args, expectedOutput)
	}
}

func TestFiles(t *testing.T) {
	testDir := t.TempDir()

	testPath1, testPath2 := path.Join(testDir, "test1.txt"), path.Join(testDir, "test2.txt")
	testContent1, testContent2 := testString1, testString2

	createTestFile(t, testPath1, testContent1)
	createTestFile(t, testPath2, testContent2)

	expectedOutputToArgs := map[string][]string{
		testContent1:                {"gocat.exe", testPath1},
		testContent1 + testContent1: {"gocat.exe", testPath1, testPath1},
		testContent2 + testContent2: {"gocat.exe", testPath2, testPath2},
		testContent1 + testContent2: {"gocat.exe", testPath1, testPath2},
		testContent2 + testContent1: {"gocat.exe", testPath2, testPath1},
	}

	for expectedOutput, args := range expectedOutputToArgs {
		testOutput(t, args, expectedOutput)
	}
}

func TestFailUnrecognized(t *testing.T) {
	args := []string{"gocat.exe", "--nonexistent-option"}
	expectedOutput := "gocat: unrecognized option '" + args[1] + "'"

	testOutput(t, args, expectedOutput)
}

func TestFailInvalidOption(t *testing.T) {
	args := []string{"gocat.exe", "-p"}
	expectedOutput := "gocat: invalid option -- '" + args[1] + "'"

	testOutput(t, args, expectedOutput)
}

func TestFailFileNotFound(t *testing.T) {
	args := []string{"gocat.exe", "non-existent-file.txt"}
	expectedOutput := "gocat: " + args[1] + ": No such file or directory"

	testOutput(t, args, expectedOutput)
}

func testOutput(t *testing.T, args []string, expectedOutput string) {
	var buffer bytes.Buffer
	Run(args, &buffer)

	output := buffer.String()
	if output != expectedOutput {
		t.Fatal(
			"Args: ", args, "\n",
			"Expected: ", expectedOutput, "\n",
			"Received: ", output, "\n",
		)
	}
}

// Replaces os.Stdin with a pipe and writes a string into it.
// Every place in code reading from os.Stdin will read that string.
func fakeStdinWrite(t *testing.T, s string) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdin = r

	_, err = w.WriteString(s)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()
}

func createTestFile(t *testing.T, path string, fileContent string) {
	err := os.WriteFile(path, []byte(fileContent), 0777)
	if err != nil {
		t.Fatal("Failed to write to", path)
	}
}
