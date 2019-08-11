package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

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

	tpl := template.Must(template.New("").Parse(storyTmpl))

	h := adventure.NewHandler(story,
		adventure.WithTemplate(tpl),
		adventure.WithPathFunc(pathFn),
	)
	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting the server on port %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

var storyTmpl = `
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<title>Choose Your Own Adventure</title>
		</head>
		<body>
			<section class="page">
				<h1>{{.Title}}</h1>
				{{range .Paragraphs}}
				<p>{{.}}</p>
				{{end}}
				<ul>
					{{range .Options}}
						<li>
							<a href="/story/{{.Chapter}}">{{.Text}}</a>
						</li>
					{{end}}
				</ul>
			</section>
			<style>
				body {
					font-family: helvetica, arial;
				}
				h1 {
					text-align: center;
					position: relative;
				}
				.page {
					width: 80%;
					max-width: 500px;
					margin: 40px auto;
					padding: 80px;
					background-color: #fffcf6;
					border: 1px solid #eee;
					box-shadow: 0 10px 6px -6px #777;
				}
				ul {
					border-top: 1px dotted #ccc;
					padding: 10px 0 0 0;
					-webkit-padding-start: 0;
				}
				li {
					padding-top: 10px;
				}
				a,
				a:visited {
					text-decoration: none;
					color: #6295b5;
				}
				a:active,
				a:hover {
					color: #779a2;
				}
				p {
					text-indent: 1em;
				}
			</style>
		</body>
	</html>
`
