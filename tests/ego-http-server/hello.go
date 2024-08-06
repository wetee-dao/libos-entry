package main

import (
	"fmt"
	"net/http"

	"github.com/wetee-dao/libos-entry/lib/ego"
)

func main() {
	err := ego.InitEgo()
	if err != nil {
		fmt.Println(err)
		return
	}
	http.HandleFunc("/", indexHandler)
	err = http.ListenAndServe(":8999", nil)
	fmt.Println(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(
		`<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
				<title>hello world</title>
			</head>
			<body>
				<h1>hello world</h1>
			</body>
		</html>`),
	)
}
