package adventure

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
)

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

var tpl *template.Template

// Story contains a list of Chapter
type Story map[string]Chapter

// Chapter for each section in adventure
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option for each decision in adventure
type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

var defaultHandlerTmpl = `
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
							<a href="/{{.Chapter}}">{{.Text}}</a>
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

// NewHandler returns a handler containing a story
func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}

// JSONStory returns the decoded story
func JSONStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}
