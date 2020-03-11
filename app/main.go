package main

import (
	"fmt"
	"strings"
    "log"
    "net/http"
)

var greetings = make(map[string]string)

func splitOutFirstPathPart(path string) (string, string) {
	parts := strings.SplitN(path, "/", 2)

	secondPart := ""	
	if len(parts) > 1 {
		secondPart = parts[1]
	}

	return parts[0], secondPart
}

func handler(w http.ResponseWriter, r *http.Request) {
	action, remainder := splitOutFirstPathPart(r.URL.Path[1:])
	log.Printf("Handling request to path %s, action is %s, remainder is %s", r.URL.Path, action, remainder)

	switch action {
		case "remember":
			handleRemember(w, r, remainder)
		case "greet":
			handleGreet(w, r, remainder)
		default:			
			w.WriteHeader(404)
			fmt.Fprintf(w, "Unkown path")
	}
}

func handleRemember(w http.ResponseWriter, r *http.Request, remainderPath string) {
	rememberedPerson, greeting := splitOutFirstPathPart(remainderPath)

	greetings[rememberedPerson] = greeting
	fmt.Fprintf(w, "Stored greeting for %s: %s", rememberedPerson, greeting)
}

func handleGreet(w http.ResponseWriter, r *http.Request, remainderPath string) {
	rememberedPerson, _ := splitOutFirstPathPart(remainderPath)

	greeting := greetings[rememberedPerson]
	
	fmt.Fprintf(w, greeting)
}

func main() {
	port := 8080
	http.HandleFunc("/", handler)
	log.Printf("Listening on port %d", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
