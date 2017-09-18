package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "I'm on it! Will inform you ASAP!")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}

	env := r.Form.Get("text")

	http.Get("https://60cqrfceu4.execute-api.eu-west-1.amazonaws.com/development?env=" + env)
}
