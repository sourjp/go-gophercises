package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sourjp/go-gophercises/url_shortener/urlshort"
)

var (
	flagYAMLFileName string
	flagJSONFileName string
	flagDB           bool
)

func init() {
	flag.StringVar(&flagYAMLFileName, "yaml", "", "accept yaml file to build short url")
	flag.StringVar(&flagJSONFileName, "json", "", "accept json file to build short url")
	flag.BoolVar(&flagDB, "db", false, "get info from db to build short url")
	flag.Parse()
}

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	if len(flagYAMLFileName) > 0 {
		yaml, err := ioutil.ReadFile(flagYAMLFileName)
		if err != nil {
			log.Fatalln(err)
		}
		yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
		if err != nil {
			log.Fatalln(err)
		}
		mapHandler = yamlHandler
	}

	// Build the JSONHandler using the mapHandler as the
	// fallback
	if len(flagJSONFileName) > 0 {
		json, err := ioutil.ReadFile(flagJSONFileName)
		if err != nil {
			log.Fatalln(err)
		}
		jsonHandler, err := urlshort.JSONHandler([]byte(json), mapHandler)
		if err != nil {
			log.Fatalln(err)
		}
		mapHandler = jsonHandler
	}

	// Build the DBHandler using the mapHandler as the
	// fallback
	if flagDB {
		conn, err := urlshort.NewDB()
		if err != nil {
			log.Fatalln(err)
		}
		dbHandler, err := urlshort.DBHandler(conn, mapHandler)
		if err != nil {
			log.Fatalln(err)
		}
		mapHandler = dbHandler
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", mapHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
