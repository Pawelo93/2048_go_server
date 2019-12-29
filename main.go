package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/load", onLoad)
	http.ListenAndServe(":8080", nil)
}

func onLoad(w http.ResponseWriter, r *http.Request) {

	answer := Answer{Direction: "LEFT"}
	b, err := json.Marshal(answer)

	if err != nil {
		fmt.Println("ERROR")
	} else {
		fmt.Fprintln(w, string(b))
		fmt.Println(string(b))
	}
}

type Answer struct {
	Direction string
}
