package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
)

func tervitaja() {
	var valik string
	praeguneAeg := time.Now().Format("15:04:05")

	fmt.Println("Tere tulemast, minu harjumuste jälgiasse!")
	fmt.Println("Kell on", praeguneAeg)
	fmt.Println("******************************")
	fmt.Println("1. Lisa uus harjumus")
	fmt.Println("******************************")
	fmt.Println("2. Näita kõiki harjumusi")
	fmt.Println("******************************")
	fmt.Println("3. Märgi harjumus tehtuks")
	fmt.Println("******************************")
	fmt.Println("4. Välju")
	fmt.Println("******************************")
	fmt.Print("Vali tegevus: ")
	fmt.Scanln(&valik)
	fmt.Println("******************************")

	switch valik {
	case "1":
		lisaHarjumus()
	case "2":
		naitaHarjumusi()
	case "3":
		margiTehtuks()
	case "4":
		fmt.Println("Head aega!")
		return
	default:
		time.Sleep(500 * time.Millisecond)
		fmt.Println("Vigane valik, proovi uuesti.")
		fmt.Println("******************************")
		time.Sleep(500 * time.Millisecond)
		tervitaja()
	}
}

func lisaHarjumus() {
	var harjumuseNimi string

	fmt.Print("Sisestage harjumuse nimi: ")
	fmt.Scanln(&harjumuseNimi)

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tabel := `CREATE TABLE IF NOT EXISTS userdata (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			harjumuseNimi TEXT NOT NULL,
			harjumusTehtud BOOL NOT NULL,
			harjumuseStriik INTEGER
		);`
	_, err = db.Exec(tabel)
	if err != nil {
		log.Fatal(err)
	}

	statement, err := db.Prepare("INSERT INTO userdata (harjumuseNimi, harjumusTehtud) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(harjumuseNimi, false)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	fmt.Println("Harjumus lisatud.")
	fmt.Println("******************************") //siin on 30 "*"
	fmt.Println()

	time.Sleep(500 * time.Millisecond)

	main()
}

func naitaHarjumusi() {

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, harjumuseNimi, harjumusTehtud, harjumuseStriik FROM userdata")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tabel := tablewriter.NewWriter(os.Stdout)
	tabel.Header([]string{"ID", "Nimi", "Striik", "Staatus"})

	for rows.Next() {
		var idString, nimi, staatusString, striikString string
		var id, striik int64
		var staatus bool
		rows.Scan(&id, &nimi, &striik, &staatus)
		idString = strconv.FormatInt(id, 10)
		striikString = strconv.FormatInt(striik, 10) + "🔥"
		if staatus {
			staatusString = "✅"
		}
		if !staatus {
			staatusString = "❎"
		}
		tabel.Append([]string{
			idString,
			nimi,
			striikString,
			staatusString,
		})
	}

	tabel.Render()

	fmt.Println("Vajutage 'Enter', et tabel sulgeda.")
	fmt.Scanln()
	main()
}

func margiTehtuks() {
	fmt.Println("3")
}

func main() {
	tervitaja()
}
