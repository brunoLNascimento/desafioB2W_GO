package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Planeta struct {
	ID             int    `json:"id"`
	PLANET_NAME    string `json:"planeta"`
	PLANET_TERRAIN string `json:"tipoTerreno"`
	PLANET_FILMS   int    `json:"qtsFilms"`
}

// PlanetHandler analisa o request e delega para função adequada
func PlanetHandler(w http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method == "GET":

		url := strings.TrimPrefix(r.URL.Path, "/planet/")
		id, _ := strconv.Atoi(url)
		pageNumber := strings.TrimPrefix(url, "page/")
		var _, err = strconv.Atoi(pageNumber)
		page, _ := strconv.Atoi(pageNumber)

		if id > 0 {
			planetPorID(w, r, id)
		} else if err != nil {
			planetPorNome(w, r, url)
		} else {
			planetTodos(w, r, page)
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Desculpa... :(")
	}
}

func planetPorID(w http.ResponseWriter, r *http.Request, id int) {
	db, err := sql.Open("mysql", "root:@/desafioGO")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	var p Planeta

	db.QueryRow("SELECT ID, PLANET_NAME, PLANET_TERRAIN, PLANET_FILMS FROM planets WHERE ID = ?", id).Scan(&p.ID, &p.PLANET_NAME, &p.PLANET_TERRAIN, &p.PLANET_FILMS)

	json, _ := json.Marshal(p)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func planetPorNome(w http.ResponseWriter, r *http.Request, id string) {
	db, err := sql.Open("mysql", "root:@/desafioGO")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	var p Planeta

	db.QueryRow("SELECT ID, PLANET_NAME, PLANET_TERRAIN, PLANET_FILMS FROM planets WHERE PLANET_NAME = ?", id).Scan(&p.ID, &p.PLANET_NAME, &p.PLANET_TERRAIN, &p.PLANET_FILMS)

	json, _ := json.Marshal(p)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}

func planetTodos(w http.ResponseWriter, r *http.Request, page int) {
	db, err := sql.Open("mysql", "root:@/desafioGO")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var limit = 10
	var OFFSET = page * limit

	rows, _ := db.Query("SELECT ID, PLANET_NAME, PLANET_FILMS, PLANET_TERRAIN FROM planets limit ? OFFSET ?", limit, OFFSET)
	defer rows.Close()

	var planetas []Planeta
	for rows.Next() {
		var planeta Planeta
		rows.Scan(&planeta.ID, &planeta.PLANET_NAME, &planeta.PLANET_FILMS, &planeta.PLANET_TERRAIN)
		planetas = append(planetas, planeta)
	}

	json, _ := json.Marshal(planetas)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))
}
