package repl

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/scanner"

	"github.com/peterh/liner"
	"k8s.io/klog/v2"
)

const (
	promptDefault  = "ksql> "
	promptContinue = "..... "
	indent         = "    "
)

type contLiner struct {
	*liner.State
	buffer string
	depth  int
}

func newContLiner() *contLiner {
	rl := liner.NewLiner()
	rl.SetCtrlCAborts(true)
	return &contLiner{State: rl}
}

func (cl *contLiner) promptString() string {
	if cl.buffer != "" {
		return promptContinue + strings.Repeat(indent, cl.depth)
	}

	return promptDefault
}

func (cl *contLiner) Prompt() (string, error) {
	line, err := cl.State.Prompt(cl.promptString())
	switch err {
	case io.EOF:
		if cl.buffer != "" {
			// cancel line continuation
			cl.Accepted()
			fmt.Println()
			err = nil
		}
	case liner.ErrPromptAborted:
		err = nil
		if cl.buffer != "" {
			cl.Accepted()
		} else {
			fmt.Println("(^D to quit)")
		}
	case nil:
		if cl.buffer != "" {
			cl.buffer = cl.buffer + "\n" + line
		} else {
			cl.buffer = line
		}
	}

	return cl.buffer, err
}

func (cl *contLiner) Accepted() {
	cl.State.AppendHistory(cl.buffer)
	cl.buffer = ""
}

func (cl *contLiner) Clear() {
	cl.buffer = ""
	cl.depth = 0
}

var errUnmatchedBraces = fmt.Errorf("unmatched braces")

func (cl *contLiner) Reindent() error {
	oldDepth := cl.depth
	cl.depth = cl.countDepth()

	if cl.depth < 0 {
		return errUnmatchedBraces
	}

	if cl.depth < oldDepth {
		lines := strings.Split(cl.buffer, "\n")
		if len(lines) > 1 {
			lastLine := lines[len(lines)-1]

			cursorUp()
			fmt.Printf("\r%s%s", cl.promptString(), lastLine)
			eraseInLine()
			fmt.Print("\n")
		}
	}

	return nil
}

func (cl *contLiner) countDepth() int {
	reader := bytes.NewBufferString(cl.buffer)
	sc := new(scanner.Scanner)
	sc.Init(reader)
	sc.Error = func(_ *scanner.Scanner, msg string) {
		klog.V(9).Infof("scanner: %s", msg)
	}

	depth := 0
	for {
		switch sc.Scan() {
		case '{', '(':
			depth++
		case '}', ')':
			depth--
		case scanner.EOF:
			return depth
		}
	}
}

func cursorUp() {
	fmt.Print("\x1b[1A")
}

func eraseInLine() {
	fmt.Print("\x1b[0K")
}
