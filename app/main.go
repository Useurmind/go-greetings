package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Useurmind/go-greetings/pkg/db"
)

var dbType string = ""
var dataSource string = ""

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
	ctx, err := db.NewDBContext(dbType, dataSource)
	if err != nil {
		log.Printf("ERROR while creating db context: %v", err)
		w.WriteHeader(500)
		return
	}

	switch action {
	case "health":
		err = handleHealth(ctx, w, r, remainder)
	case "remember":
		err = handleRemember(ctx, w, r, remainder)
	case "greet":
		err = handleGreet(ctx, w, r, remainder)
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

func handleHealth(ctx db.DBContext, w http.ResponseWriter, r *http.Request, remainderPath string) error {
	w.WriteHeader(200)
	fmt.Fprintf(w, "Working")

	return nil
}

func handleRemember(ctx db.DBContext, w http.ResponseWriter, r *http.Request, remainderPath string) error {
	rememberedPerson, greeting := splitOutFirstPathPart(remainderPath)

	err := ctx.SaveGreeting(rememberedPerson, greeting)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Stored greeting for %s: %s", rememberedPerson, greeting)
	return nil
}

func handleGreet(ctx db.DBContext, w http.ResponseWriter, r *http.Request, remainderPath string) error {
	rememberedPerson, _ := splitOutFirstPathPart(remainderPath)

	greeting, err := ctx.GetGreeting(rememberedPerson)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, *greeting)
	return nil
}

func main() {
	port := 8080

	dbType = os.Getenv("GOGREETING_DBTYPE")
	dataSource = os.Getenv("GOGREETING_DATASOURCE")
	if dbType == "" {
		panic("No db type specified, set GOGREETING_DBTYPE")
	}
	if dataSource == "" {
		panic("No datasource specified, set GOGREETING_DATASOURCE")
	}

	fmt.Printf("Using %s\r\n", dbType)

	http.HandleFunc("/", handler)
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
