package main

import (
	"./api"
	"./dotenv_parser"
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	config, err := dotenv_parser.Parse("./.env")
	handleErr(err)

	scheme := config["SCHEME"]
	host := config["HOST"]
	port := config["PORT"]

	urls, err := api.New("mongo:27017", "urls", "urls")
	handleErr(err)

	defer urls.Session.Close()

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("./views/index.html")
		handleErr(err)

		w.Header().Add("Content-type", "text/html; charset=utf8")

		handleErr(t.Execute(w, map[string]string{
			"scheme": scheme,
			"host":   host,
			"port":   port,
		}))
	})

	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		link := r.URL.Query().Get("link")

		query := bson.M{"link": link}
		hasDocument, err := urls.Has(query)
		handleErr(err)

		var url api.URL
		if hasDocument {
			url, err = urls.Get(query, false, scheme, host, port)
			handleErr(err)
		} else {
			url, err = urls.Insert(link, scheme, host, port)
			handleErr(err)
		}

		responseJson, err := json.Marshal(url)
		handleErr(err)

		w.Header().Set("Content-type", "application/json")
		_, err = w.Write(responseJson)

		handleErr(err)
	}).Methods("GET")

	r.HandleFunc("/{hash:[0-9a-zA-Z]+}", func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"]
		query := bson.M{"hash": hash}

		has, err := urls.Has(query)
		handleErr(err)

		if has {
			url, err := urls.Get(query, false, scheme, host, port)
			handleErr(err)
			err = urls.IncTransitions(hash)
			handleErr(err)

			http.Redirect(w, r, url.Link, 302)
		} else {
			http.NotFound(w, r)
		}
	}).Methods("GET")

	log.Println("Start server at the port " + port)
	handleErr(http.ListenAndServe(":"+port, r))
}
