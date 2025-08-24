package model

import (
	"encoding/json"
	"fmt"
)

func PrintJson(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println("------------------ json --------------------------")
	fmt.Println(string(b))
	fmt.Println("------------------ json --------------------------")
}
