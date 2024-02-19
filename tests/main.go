package main

import (
	"fmt"

	"github.com/wetee-dao/libos-entry/lib/ego"
)

func main() {
	err := ego.InitEgo("")
	if err != nil {
		fmt.Println(err)
	}
}
