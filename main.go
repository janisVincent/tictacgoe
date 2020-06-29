package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type gridPage struct {
	Title string
	Cells [][]int
}

type dtoInCell struct {
	Cell int
}

type dtoOutSymbol struct {
	Symbol string
}

var turn int

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := template.New("grid")
	t = template.Must(t.ParseFiles("public/tmpl/layout.tmpl", "public/tmpl/partials/grid.tmpl"))

	turn = 0
	cells := make([][]int, 3)

	for i := 0; i < 3; i++ {
		for j := 1; j <= 3; j++ {
			cells[i] = append(cells[i], i*3+j)
		}
	}

	p := gridPage{"Tic Tac Goe", cells}
	err := t.ExecuteTemplate(w, "layout", p)

	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}

func play(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var cell dtoInCell
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cell)

	if err != nil {
		panic(err)
	}

	symbols := []string{"X", "O"}
	var key int
	if turn < 1 {
		key = rand.Intn(1)
	} else {
		key = turn % 2
	}
	turn++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtoOutSymbol{symbols[key]})
}

func main() {
	router := httprouter.New()

	router.GET("/", index)
	router.POST("/play", play)

	router.ServeFiles("/css/*filepath", http.Dir("public/assets/css"))
	router.ServeFiles("/js/*filepath", http.Dir("public/assets/js"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
