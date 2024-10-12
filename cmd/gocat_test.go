package gocat

// TODO: use a docker container with cat and diff instead?
// TODO: maybe move parameter test's expected output into their own files at least?
// TODO: move into gocat_test

import (
	"bytes"
	"os"
	"path"
	"testing"
)

// Note: cat writes "\n" after its output, so for clarity keeping '+ "\n"'
// and also prepending it for every test string. TODO: check if this comment is correct.
const (
	testString1 string = "test string 1" + "\n"
	testString2 string = "test string 2" + "\n"
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

func TestParameters(t *testing.T) {
	var expectedOutput string
	testContent := `
	hello,
	world 1!

       


	hey 1

	!@)(*$&)*!#%^%!_)+_




	` + "\n"

	/*args = []string{"gocat.exe", "-A"}
	args = []string{"gocat.exe", "--show-all"}
	args = []string{"gocat.exe", "-vET"}*/

	// -b, --number-nonblank, (-b, -n), (-n, -b) | number nonempty output lines, overrides -n
	expectedOutput = `
     1		hello,
     2		world 1!

     3	       


     4		hey 1

     5		!@)(*$&)*!#%^%!_)+_




     6		` + "\n"

	testParameters(t, []string{"-b"}, testContent, expectedOutput)
	testParameters(t, []string{"--number-nonblank"}, testContent, expectedOutput)
	testParameters(t, []string{"-b", "-n"}, testContent, expectedOutput)
	testParameters(t, []string{"-n", "-b"}, testContent, expectedOutput)

	/*args = []string{"gocat.exe", "-e"}
	args = []string{"gocat.exe", "-vE"}
	args = []string{"gocat.exe", "-Ev"}*/

	// -E, --show-ends | display $ at end of each line
	expectedOutput = `$
	hello,$
	world 1!$
$
       $
$
$
	hey 1$
$
	!@)(*$&)*!#%^%!_)+_$
$
$
$
$
	$` + "\n"
	testParameters(t, []string{"-E"}, testContent, expectedOutput)
	testParameters(t, []string{"--show-ends"}, testContent, expectedOutput)

	// -n, --number | number all output lines
	expectedOutput = `     1	
     2		hello,
     3		world 1!
     4	
     5	       
     6	
     7	
     8		hey 1
     9	
    10		!@)(*$&)*!#%^%!_)+_
    11	
    12	
    13	
    14	
    15		` + "\n"

	testParameters(t, []string{"-n"}, testContent, expectedOutput)
	testParameters(t, []string{"--number"}, testContent, expectedOutput)

	// -s, --squeeze-blank | suppress repeated empty output lines
	expectedOutput = `
	hello,
	world 1!

       

	hey 1

	!@)(*$&)*!#%^%!_)+_

	` + "\n"

	testParameters(t, []string{"-s"}, testContent, expectedOutput)
	testParameters(t, []string{"--squeeze-blank"}, testContent, expectedOutput)

	/*args = []string{"gocat.exe", "-t"}
	args = []string{"gocat.exe", "-vT"}
	args = []string{"gocat.exe", "-Tv"}*/

	// -T, --show-tabs | display TAB characters as ^I
	expectedOutput = `
^Ihello,
^Iworld 1!

       


^Ihey 1

^I!@)(*$&)*!#%^%!_)+_




^I` + "\n"
	testParameters(t, []string{"-T"}, testContent, expectedOutput)
	testParameters(t, []string{"--show-tabs"}, testContent, expectedOutput)

	// -u | ignored
	expectedOutput = testContent
	testParameters(t, []string{"-u"}, testContent, expectedOutput)

	/*args = []string{"gocat.exe", "-v"}
	args = []string{"gocat.exe", "--show-nonprinting"}*/

	// TODO: reuse strings in gocat.go with go:linkname
	//expectedOutput =
	//testParameter(t, "--help", testContent, expectedOutput)
	//expectedOutput =
	//testParameter(t, "--version", testContent, expectedOutput)
}

func testParameters(t *testing.T, params []string, testString string, expectedOutput string) {
	fakeStdinWrite(t, testString)
	args := append([]string{"gocat.exe"}, params...)
	testOutput(t, args, expectedOutput)
}

func testOutput(t *testing.T, args []string, expectedOutput string) {
	var buffer bytes.Buffer
	Run(args, &buffer)

	output := buffer.String()
	if output != expectedOutput {
		t.Fatal(
			"\nArgs: ", args,
			"\nExpected:\n", expectedOutput,
			"\nReceived:\n", output,
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
