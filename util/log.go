package util

import (
	"fmt"

	"github.com/fatih/color"
)

func LogWithRed(tag string, a ...interface{}) {
	red := color.New(color.FgRed).SprintFunc()
	b := make([]interface{}, 0, len(a)+1)
	b = append(b, red(tag+":"))
	b = append(b, a...)
	fmt.Println(b...)
}
