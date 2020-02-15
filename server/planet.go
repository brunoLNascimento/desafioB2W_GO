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
	sid := strings.TrimPrefix(r.URL.Path, "/planet/")
	id, _ := strconv.Atoi(sid)

	switch {
	case r.Method == "GET" && id > 0:
		planetPorID(w, r, id)
	case r.Method == "GET":
		planetTodos(w, r)
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

func planetTodos(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@/desafioGO")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, _ := db.Query("SELECT ID, PLANET_NAME, PLANET_FILMS, PLANET_TERRAIN FROM planets")

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
