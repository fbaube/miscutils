package miscutils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PressAnyKey waits until any key is pressed.
func PressAnyKey() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "icanon", "min", "1").Run()
	fmt.Printf("Press any key to continue...")
	var b = make([]byte, 1)
	os.Stdin.Read(b)
	fmt.Println("")
	// fmt.Println("PressAnyKey() got the byte", b, "("+string(b)+")")
}

// InteractiveInput displays the arg as a prompt and then
// reads a String entered by the user (ended by Enter).
func InteractiveInput(prompt string) string {
	if prompt != "" {
		fmt.Printf("%s ", prompt)
	}
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stdout, "InteractiveInput failed: "+err.Error())
		return ""
	}
	// sanitize input
	return SanitizeInput(input)
}

// GetKeypress is a nifty hack to read any key.
func GetKeypress() byte {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var b = make([]byte, 1)
	os.Stdin.Read(b)
	fmt.Println("I got the byte", b, "("+string(b)+")")
	return b[0]
}

// SanitizeInput applies the following rules iteratively
// until no further processing can be done:
// :ol:
// :: trim all the extra white spaces
// :: trim all return carriage chars
// :: trim leading / ending quotation marks (ex.: "my text")
// :: trim leading / ending spaces
// -ol-
func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return input
	}
	input = strings.TrimPrefix(input, "\"")
	input = strings.TrimSuffix(input, "\"")
	return strings.NewReplacer("  ", " ", "\n", " ", "\t", " ", "\r", " ").Replace(input)
}
