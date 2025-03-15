package console

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	PASS    = "PASS"
	FAIL    = "FAIL"
	SKIP    = "SKIP"
	UNKNOWN = "UNKNOWN"
)

var (
	length    = 150
	frameLine = strings.Repeat("-", length)
)

func PrintPass(out string) {
	pass := color.New(color.FgHiWhite, color.BgGreen).SprintfFunc()
	out = strings.Replace(out, PASS, pass(PASS), 1)
	fmt.Print(out)
	color.Blue(frameLine)
}

func PrintFail(out string) {
	fail := color.New(color.FgHiWhite, color.BgRed).SprintfFunc()
	out = strings.Replace(out, FAIL, fail(FAIL), 1)
	fmt.Print(out)
	color.Blue(frameLine)
}

func PrintSkip(out string) {
	skip := color.New(color.FgHiWhite, color.BgYellow).SprintfFunc()
	out = strings.Replace(out, SKIP, skip(SKIP), 1)
	fmt.Print(out)
	color.Blue(frameLine)
}

func PrintUnknown(out string) {
	unknown := color.New(color.FgHiWhite, color.BgMagenta).SprintfFunc()
	out = strings.Replace(out, UNKNOWN, unknown(UNKNOWN), 1)
	fmt.Print(out)
	color.Blue(frameLine)
}

func PrintDiffWithColors(diff []byte) {
	green := color.New(color.FgHiWhite).Add(color.BgGreen)
	red := color.New(color.FgHiWhite).Add(color.BgRed)

	lines := string(diff)
	for _, line := range strings.Split(lines, "\n") {
		switch {
		case strings.HasPrefix(line, "+++"):
			green.Println(line)
		case strings.HasPrefix(line, "---"):
			red.Println(line)
		case strings.HasPrefix(line, "+"):
			green.Println(line)
		case strings.HasPrefix(line, "-"):
			red.Println(line)
		default:
			fmt.Println(line)
		}
	}
}
