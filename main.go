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

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	name := r.FormValue("name")
	fmt.Fprintf(w, "Name = %s\n", name)

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
