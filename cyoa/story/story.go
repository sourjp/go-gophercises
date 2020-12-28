package story

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"
)

var (
	storyFileName  string = "gopher.json"
	storiesMap     Stories
	defaultPageKey string = "intro"
	tplFileName    string = "tmpl.html"
	tpl            *template.Template
)

// Stories compose Chapter title and Story
type Stories map[string]Story

// Story descrive chapter
type Story struct {
	Title   string
	Story   []string
	Options []Option
}

// Option describe selectioin and next chapter
type Option struct {
	Text string
	Arc  string
}

func init() {
	tpl = template.Must(template.ParseFiles(tplFileName))
	storiesData, err := ioutil.ReadFile(storyFileName)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(storiesData, &storiesMap)
}

// ViewPageHandler return template page with story
func ViewPageHandler(w http.ResponseWriter, r *http.Request) {
	// Path includes "/"(e.g. /path), so get the key after that
	key := r.URL.Path[1:]
	story := storiesMap[key]
	if len(story.Title) <= 0 {
		http.Redirect(w, r, "/intro", http.StatusPermanentRedirect)
	}
	tpl.Execute(w, story)
}
