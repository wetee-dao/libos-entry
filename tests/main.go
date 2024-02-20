package main

import (
	"fmt"
	"net/http"

	"github.com/wetee-dao/libos-entry/lib/ego"
)

func main() {
	err := ego.InitEgo("")
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/", indexHandler)
	fmt.Println("Start http://0.0.0.0:8999 ...")
	http.ListenAndServe(":8999", nil)
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
