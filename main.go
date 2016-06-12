package main

import (
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"

	"database/sql"
	"encoding/json"
	"encoding/xml"
	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	gmux "github.com/gorilla/mux"
	"github.com/yosssi/ace"
	"io/ioutil"
	"net/url"
)

type Page struct {
	Books  []Book
	Filter string
}

type Book struct {
	PK             int64  `db:"pk"`
	Title          string `db:"title"`
	Author         string `db:"author"`
	Classification string `db:"classification"`
	ID             string `db:"id"`
}

type SearchResult struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year   string `xml:"hyr,attr"`
	ID     string `xml:"owi,attr"`
}

var db *sql.DB
var dbmap *gorp.DbMap

func initDb() {
	db, _ = sql.Open("sqlite3", "dev.db")

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(Book{}, "books").SetKeys(true, "pk")
	dbmap.CreateTablesIfNotExists()
}

func verifyDatabase(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	next(w, r)
}

func getStringFromSession(r *http.Request, key string) string {
	var strValue string
	val := sessions.GetSession(r).Get(key)
	if val != nil {
		strValue = val.(string)
	}
	return strValue
}

func getBookCollection(books *[]Book, sortCol string, filterByClass string, w http.ResponseWriter) bool {
	if sortCol == "" {
		sortCol = "pk"
	}

	var where string
	if filterByClass == "fiction" {
		where = " where classification between '800' and '900'"
	} else if filterByClass == "nonfiction" {
		where = " where classification not between '800' and '900'"
	}

	_, err := dbmap.Select(books, "select * from books"+where+" order by "+sortCol)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	return true
}

func main() {
	fmt.Println("Go web development ( ͡° ͜ʖ ͡°)")

	initDb()
	mux := gmux.NewRouter()

	// Index page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template, err := ace.Load("templates/index", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p := Page{Books: []Book{}, Filter: getStringFromSession(r, "Filter")}
		if !getBookCollection(&p.Books, getStringFromSession(r, "SortBy"), getStringFromSession(r, "Filter"), w) {
			return
		}

		err = template.Execute(w, p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// Login page
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("register") != "" || r.FormValue("login") != "" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		template, err := ace.Load("templates/login", "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = template.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Serve static files
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))).Methods("GET")

	// Book search
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
	}).Methods("POST")

	// Add book
	mux.HandleFunc("/books/{pk}", func(w http.ResponseWriter, r *http.Request) {
		book, err := Find(gmux.Vars(r)["pk"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := Book{
			PK:             -1,
			Title:          book.BookData.Title,
			Author:         book.BookData.Author,
			Classification: book.Classification.MostPopular,
		}
		err = dbmap.Insert(&b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("PUT")

	// Delete book
	mux.HandleFunc("/books/{pk}", func(w http.ResponseWriter, r *http.Request) {
		pk, _ := strconv.ParseInt(gmux.Vars(r)["pk"], 10, 64)
		b := &Book{
			PK:             pk,
			Title:          "",
			Author:         "",
			ID:             "",
			Classification: "",
		}
		_, err := dbmap.Delete(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	// Sort books by column
	mux.HandleFunc("/books/{columnName:title|author|classification}", func(w http.ResponseWriter, r *http.Request) {
		columnName := gmux.Vars(r)["columnName"]

		var b []Book
		if !getBookCollection(&b, columnName, getStringFromSession(r, "Filter"), w) {
			return
		}

		sessions.GetSession(r).Set("SortBy", columnName)

		err := json.NewEncoder(w).Encode(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// Sort books by fiction
	mux.HandleFunc("/books/{filter:all|fiction|nonfiction}", func(w http.ResponseWriter, r *http.Request) {
		filter := gmux.Vars(r)["filter"]

		var b []Book
		if !getBookCollection(&b, getStringFromSession(r, "SortBy"), filter, w) {
			return
		}

		sessions.GetSession(r).Set("Filter", filter)

		err := json.NewEncoder(w).Encode(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	n := negroni.Classic()
	n.Use(sessions.Sessions("go-book-db", cookiestore.New([]byte("password123"))))
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
