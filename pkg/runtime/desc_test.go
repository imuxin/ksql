package runtime

import (
	"fmt"
	"testing"

	// "github.com/alecthomas/chroma/v2/quick"

	"github.com/fatih/color"
)

func TestXxx(_ *testing.T) {
	// _ = quick.Highlight(os.Stdout, "package main", "go", "html", "monokai")
	color.NoColor = false
	fmt.Println("This", color.RedString("warning"), "should be not neglected.")
	// // Create SprintXxx functions to mix strings with other non-colorized strings:
	// yellow := color.New(color.FgYellow).SprintFunc()
	// red := color.New(color.FgRed).SprintFunc()
	// fmt.Printf("This is a %s and this is %s.\n", yellow("warning"), red("error"))

	// info := color.New(color.FgWhite, color.BgGreen).SprintFunc()
	// fmt.Printf("This %s rocks!\n", info("package"))

	// // Use helper functions
	// fmt.Println("This", color.RedString("warning"), "should be not neglected.")
	// fmt.Printf("%v %v\n", color.GreenString("Info:"), "an important message.")

	// // Windows supported too! Just don't forget to change the output to color.Output
	// fmt.Fprintf(color.Output, "Windows support: %s", color.GreenString("PASS"))
}
