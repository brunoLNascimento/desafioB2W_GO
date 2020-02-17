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
	url := strings.TrimPrefix(r.URL.Path, "/planet/")
	id, _ := strconv.Atoi(url)

	switch {
	case r.Method == "GET":

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
	case r.Method == "DELETE":
		if id <= 0 {
			fmt.Fprint(w, string("O ID planeta é obrigatório!"))
			return
		}
		delete(w, r, id)
	case r.Method == "POST":
		savePlanet(w, r)
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

func delete(w http.ResponseWriter, r *http.Request, id int) {

	db, err := sql.Open("mysql", "root:@/desafioGO")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	res, err := db.Exec("DELETE FROM planets WHERE ID = ?", id)

	if err == nil {
		count, err := res.RowsAffected()
		if err != nil {
			fmt.Fprint(w, string("Erro ao tentar deletar planeta!"))
		} else if count > 0 {
			fmt.Fprint(w, string("Planeta deletado com sucesso!"))
		} else {
			fmt.Fprint(w, string("Nenhum planeta encontrado!"))
		}
	}
}

func savePlanet(w http.ResponseWriter, r *http.Request) {

	var p Planeta

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		fmt.Println("err => ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if p.PLANET_NAME == "" {
		fmt.Fprint(w, string("Nome planeta é obrigatório!"))
	}

	db, err := sql.Open("mysql", "root:@/desafioGO")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	stmt, err := db.Prepare("INSERT FROM planets (PLANET_NAME, PLANET_TERRAIN,PLANET_FILMS) VALUES (?,?,?)")
	stmt.Exec(p.PLANET_NAME)
	stmt.Exec(p.PLANET_TERRAIN)
	stmt.Exec(p.PLANET_FILMS)

	// stmt, _ := db.Prepare("insert into usuarios(nome) values(?)")
	// stmt.Exec("Maria")
	// stmt.Exec("João")

	res, _ := stmt.Exec(p)

	id, _ := res.LastInsertId()
	fmt.Println(id)

	linhas, _ := res.RowsAffected()
	fmt.Println(linhas)

	/*func Insert(w http.ResponseWriter, r *http.Request) {
		db := dbConn()
		if r.Method == "POST" {
			name := r.FormValue("name")
			city := r.FormValue("city")
			insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)")
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(name, city)
			log.Println("INSERT: Name: " + name + " | City: " + city)
		}
		defer db.Close()
		http.Redirect(w, r, "/", 301)
	}

	if err == nil {
		count, err := res.RowsAffected()
		if err != nil {
			fmt.Fprint(w, string("Erro ao tentar deletar planeta!"))
		} else if count > 0 {
			fmt.Fprint(w, string("Planeta deletado com sucesso!"))
		} else {
			fmt.Fprint(w, string("Nenhum planeta encontrado!"))
		}
	}*/
}
