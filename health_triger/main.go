package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func notifySlack(env string) {
	http.Get("https://60cqrfceu4.execute-api.eu-west-1.amazonaws.com/development?env=" + env)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	env := r.FormValue("text")

	fmt.Fprintf(w, "Checking all apps for %v environment. This might take some time", env)

	go notifySlack(env)

	time.Sleep(1 * time.Second)
	fmt.Println("all done!")
}
