package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Continent struct {
	Code        string
	Name        string
	Description string
}

type Country struct {
	ID            string
	Code          string
	CodeLow       string
	Name          string
	OfficialName  string
	Iso3          string
	Number        string
	Currency      string
	Coord         []float32
	Capital       string
	Area          string
	ContinentCode string
}

var database *sql.DB

func main() {
	db, err := sql.Open("mysql", "root:password@/geoworld")
	if err != nil {
		log.Println(err)
	}

	database = db
	defer db.Close()

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.HandleFunc("/", getContinents)
	http.HandleFunc("/list", getCountriesByContinent)
	http.HandleFunc("/detail", getCountryById)

	fmt.Println("\nServer is listening on port 8181...")
	_ = http.ListenAndServe(":8181", nil)
}

func getContinents(writer http.ResponseWriter, request *http.Request) {
	rows, err := database.Query("SELECT code, name, description FROM continents")

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	continents := []Continent{}

	for rows.Next() {
		con := Continent{}
		err := rows.Scan(&con.Code, &con.Name, &con.Description)
		if err != nil {
			fmt.Println(err)
			continue
		}
		con.Code = strings.ToLower(con.Code)
		continents = append(continents, con)
	}

	tmpl, _ := template.ParseFiles("pages/index.html", "templates/header.html",
		"templates/navbar.html", "templates/footer.html")
	tmpl.Execute(writer, continents)
}

func getCountriesByContinent(writer http.ResponseWriter, request *http.Request) {
	code := request.URL.Query().Get("code")
	rows, err := database.Query("SELECT country_id, code, name, capital, area FROM countries " +
		"WHERE continent_code = '" + code + "' " + "ORDER BY display_order;")

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	countries := []Country{}

	for rows.Next() {
		cou := Country{}
		err := rows.Scan(&cou.ID, &cou.Code, &cou.Name, &cou.Capital, &cou.Area)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cou.CodeLow = strings.ToLower(cou.Code)
		countries = append(countries, cou)
	}

	tmpl, _ := template.ParseFiles("pages/countries.html", "templates/header.html",
		"templates/navbar.html", "templates/footer.html")
	tmpl.Execute(writer, countries)
}

func getCountryById(writer http.ResponseWriter, request *http.Request) {
	code := request.URL.Query().Get("id")
	rows, err := database.Query("SELECT code, name, official_name, iso3, number, currency, coords, capital, area, continent_code " +
		"FROM countries WHERE country_id = '" + code + "';")

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	country := Country{}

	var coords string
	for rows.Next() {
		err := rows.Scan(&country.Code, &country.Name, &country.OfficialName, &country.Iso3, &country.Number,
			&country.Currency, &coords, &country.Capital, &country.Area, &country.ContinentCode)
		if err != nil {
			fmt.Println(err)
			continue
		}
		country.CodeLow = strings.ToLower(country.Code)
		country.ContinentCode = strings.ToLower(country.ContinentCode)
		regCurrency := regexp.MustCompile("\\W")
		country.Currency = regCurrency.ReplaceAllString(country.Currency, "")

		err = json.Unmarshal([]byte(coords), &country.Coord)
		if err != nil {
			fmt.Println("Json parse error:", err)
		}
	}

	tmpl, _ := template.ParseFiles("pages/countryDetail.html", "templates/header.html",
		"templates/navbar.html", "templates/footer.html")
	tmpl.Execute(writer, country)
}