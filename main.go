package main

import (
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/julienschmidt/httprouter"
)

// Config parameters
type Config struct {
	FileSystem struct {
		TemplatesDir string `yaml:"templates_dir"`
	} `yaml:"fs"`
	Application struct {
		NbCells int      `yaml:"nb_cells"`
		Tokens  []string `yaml:"tokens"`
	} `yaml:"application"`
}

// GridPage structure
type GridPage struct {
	Title string
	Cells [][]int
}

// Token for users
type Token struct {
	Token string `json:"token"`
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := template.New("grid")
	t = template.Must(t.ParseFiles(cfg.FileSystem.TemplatesDir+"layout.tmpl", cfg.FileSystem.TemplatesDir+"/partials/grid.tmpl"))

	turn = 0
	nbCells := cfg.Application.NbCells
	cells := make([][]int, nbCells)

	for i := 0; i < nbCells; i++ {
		for j := 1; j <= nbCells; j++ {
			cells[i] = append(cells[i], i*nbCells+j)
		}
	}

	err := t.ExecuteTemplate(w, "layout", GridPage{"Tic Tac Goe", cells})

	if err != nil {
		log.Fatalf("Template execution: %s", err)
	}
}

func play(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var i int
	if turn < 1 {
		i = rand.Intn(2)
		turn = i
	} else {
		i = turn % 2
	}
	turn++

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Token{cfg.Application.Tokens[i]})
}

func readConfig(cfg *Config) {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal("Unable to read config", err)
	}

	defer f.Close()

	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		log.Fatal("Unable to parse config", err)
	}
}

var cfg Config
var turn int

func main() {
	readConfig(&cfg)

	router := httprouter.New()

	router.GET("/", index)
	router.POST("/play", play)

	router.ServeFiles("/css/*filepath", http.Dir("public/css"))
	router.ServeFiles("/js/*filepath", http.Dir("public/js"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
