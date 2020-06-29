package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const nbCells = 3
const templatesDir = "public/templates/"

var turn int

type gridPage struct {
	Title string
	Cells [][]int
}

type cell struct {
	Cell int `json:"cell"`
}

type token struct {
	Token string `json:"token"`
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := template.New("grid")
	t = template.Must(t.ParseFiles(templatesDir+"layout.tmpl", templatesDir+"/partials/grid.tmpl"))

	turn = 0
	cells := make([][]int, nbCells)

	for i := 0; i < nbCells; i++ {
		for j := 1; j <= nbCells; j++ {
			cells[i] = append(cells[i], i*nbCells+j)
		}
	}

	err := t.ExecuteTemplate(w, "layout", gridPage{"Tic Tac Goe", cells})

	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}

func play(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var cell cell
	err := json.NewDecoder(r.Body).Decode(&cell)

	if err != nil {
		panic(err)
	}

	var i int
	if turn < 1 {
		i = rand.Intn(2)
		turn = i
	} else {
		i = turn % 2
	}
	turn++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token{[]string{"X", "O"}[i]})
}

func main() {
	router := httprouter.New()

	router.GET("/", index)
	router.POST("/play", play)

	router.ServeFiles("/css/*filepath", http.Dir("public/css"))
	router.ServeFiles("/js/*filepath", http.Dir("public/js"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
