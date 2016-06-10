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
	Books []Book
}

type Book struct {
	PK             int
	Title          string
	Author         string
	Classification Classification
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
			return
		}

		p := Page{Books: []Book{}}
		rows, _ := db.Query("select pk,title,author,classification from books")
		for rows.Next() {
			var b Book
			rows.Scan(&b.PK, &b.Title, &b.Author, &b.Classification.MostPopular)
			p.Books = append(p.Books, b)
		}
		err = template.Execute(w, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		results, err := Search(r.FormValue("search"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/books/add", func(w http.ResponseWriter, r *http.Request) {
		book, err := Find(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result, err := db.Exec("insert into books (pk, title, author, id, classification) values (?, ?, ?, ?, ?)",
			nil, book.BookData.Title, book.BookData.Author, book.BookData.ID, book.Classification.MostPopular)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pk, _ := result.LastInsertId()
		b := Book{
			PK: int(pk),
			Title: book.BookData.Title,
			Author: book.BookData.Author,
			Classification: Classification{MostPopular: book.Classification.MostPopular},
		}

		err = json.NewEncoder(w).Encode(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/books/delete", func(w http.ResponseWriter, r *http.Request) {
		_, err := db.Exec("delete from books where pk = ?", r.FormValue("pk"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		file, err := ioutil.ReadFile(r.URL.Path[1:])
		if err != nil {
			http.NotFound(w, r)
			return
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
	Classification Classification `xml:"recommendations>ddc>mostPopular"`
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
