package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

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

	var err error = nil
	switch action {
	case "health":
		err = handleHealth(w, r, remainder)
	case "remember":
		err = handleRemember(w, r, remainder)
	case "greet":
		err = handleGreet(w, r, remainder)
	default:
		w.WriteHeader(404)
		fmt.Fprintf(w, "Unkown path")
		return
	}

	if err != nil {
		log.Printf("ERROR while executing action %s: %s", action, err)
		w.WriteHeader(500)
		fmt.Fprintf(w, "An error occured while action %s was performed", action)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request, remainderPath string) error {
	_, err := getDB()
	if err == nil {
		w.WriteHeader(200)
		fmt.Fprintf(w, "Working")
	}

	return err
}

func handleRemember(w http.ResponseWriter, r *http.Request, remainderPath string) error {
	rememberedPerson, greeting := splitOutFirstPathPart(remainderPath)

	err := saveGreeting(rememberedPerson, greeting)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Stored greeting for %s: %s", rememberedPerson, greeting)
	return nil
}

func handleGreet(w http.ResponseWriter, r *http.Request, remainderPath string) error {
	rememberedPerson, _ := splitOutFirstPathPart(remainderPath)

	greeting, err := getGreeting(rememberedPerson)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, *greeting)
	return nil
}

func main() {
	port := 8080

	var dataSource string
	if len(os.Args) > 1 {
		dataSource = os.Args[1]
	} else {
		dataSource = os.Getenv("GOGREETING_DATASOURCE")
	}
	if dataSource == "" {
		panic("No datasource specified, either hand it over as first argument to the cli or set GOGREETING_DATASOURCE")
	}

	setDataSource(dataSource)

	http.HandleFunc("/", handler)
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
