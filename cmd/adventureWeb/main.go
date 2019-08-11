package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mrpineapples/adventure"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the Adventure web app on")
	filename := flag.String("file", "gopher.json", "the JSON file with the create your own adeventure story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	story, err := adventure.JSONStory(f)
	if err != nil {
		panic(err)
	}

	h := adventure.NewHandler(story)
	fmt.Printf("Starting the server on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
