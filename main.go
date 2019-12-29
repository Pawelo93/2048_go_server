package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/load", onLoad)
	http.ListenAndServe(":8080", nil)
}

func onLoad(w http.ResponseWriter, r *http.Request) {
	// for i := 1; i < 5; i++ {
	// 	fmt.Fprintln(w, "Calling /load endpoint")
	// }
	var buffer = new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	var answer = Answer{direction: "LEFT"}
	fmt.Println(answer)
	// fmt.Println(encoder.Encode(answer))
	encoder.Encode(answer)
	io.Copy(os.Stdout, buffer)
	// fmt.Fprintln(w, buffr)
}

type Answer struct {
	direction string
}
