package main

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"encoding/json"
	"encoding/xml"
	"github.com/codegangsta/negroni"
	"github.com/yosssi/ace"
	"io/ioutil"
	"net/url"
)

type Page struct {
	Name     string
	DBStatus bool
}

type SearchResult struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year   string `xml:"hyr,attr"`
	ID     string `xml:"owi,attr"`
}

var db *sql.DB

func verifyDatabase(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	next(w, r)
}

func main() {
	fmt.Println("Go web development ( ͡° ͜ʖ ͡°)")

	db, _ = sql.Open("sqlite3", "dev.db")
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template, err := ace.Load("templates/index", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = template.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		results, err := Search(r.FormValue("search"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/books/add", func(w http.ResponseWriter, r *http.Request) {
		book, err := Find(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = db.Exec("insert into books (pk, title, author, id, classification) values (?, ?, ?, ?, ?)",
			nil, book.BookData.Title, book.BookData.Author, book.BookData.ID, book.Classification.MostPopular)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		file, err := ioutil.ReadFile("static/" + r.URL.Path[8:])
		if err != nil {
			http.NotFound(w, r)
		}

		w.Write(file)
	})

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(verifyDatabase))
	n.UseHandler(mux)
	n.Run(":80")
}

type ClassifySearchResponse struct {
	Results []SearchResult `xml:"works>work"`
}

type BookData struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	ID     string `xml:"owi,attr"`
}

type Classification struct {
	MostPopular string `xml:"sfa,attr"`
}
type ClassifyBookResponse struct {
	BookData       BookData       `xml:"work"`
	Classification Classification `xml:"recommendations>dcc>mostPopular"`
}

func Find(id string) (ClassifyBookResponse, error) {
	body, err := classifyApi("http://classify.oclc.org/classify2/Classify?summary=true&owi=" + url.QueryEscape(id))
	if err != nil {
		return ClassifyBookResponse{}, err
	}

	var c ClassifyBookResponse
	err = xml.Unmarshal(body, &c)
	return c, err
}

func Search(query string) ([]SearchResult, error) {
	query_url := "http://classify.oclc.org/classify2/Classify?summary=true&title=" + url.QueryEscape(query)

	body, err := classifyApi(query_url)
	if err != nil {
		return []SearchResult{}, err
	}

	var c ClassifySearchResponse
	err = xml.Unmarshal(body, &c)
	return c.Results, err
}

func classifyApi(req_url string) ([]byte, error) {
	resp, err := http.Get(req_url)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
