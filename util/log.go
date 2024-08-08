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

func LogError(tag string, a ...interface{}) {
	red := color.New(color.FgHiRed).SprintFunc()
	b := make([]interface{}, 0, len(a)+1)
	b = append(b, red(tag+":"))
	b = append(b, a...)
	fmt.Println(b...)
}

func LogOk(tag string, a ...interface{}) {
	red := color.New(color.FgHiGreen).SprintFunc()
	b := make([]interface{}, 0, len(a)+1)
	b = append(b, red(tag+":"))
	b = append(b, a...)
	fmt.Println(b...)
}

func LogSendmsg(tag string, a ...interface{}) {
	red := color.New(color.FgHiCyan).SprintFunc()
	b := make([]interface{}, 0, len(a)+1)
	b = append(b, red(tag+":"))
	b = append(b, a...)
	fmt.Println(b...)
}

func LogRevmsg(tag string, a ...interface{}) {
	red := color.New(color.FgHiMagenta).SprintFunc()
	b := make([]interface{}, 0, len(a)+1)
	b = append(b, red(tag+":"))
	b = append(b, a...)
	fmt.Println(b...)
}
