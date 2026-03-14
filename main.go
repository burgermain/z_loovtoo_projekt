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
	case "plsporn":
		easterEgg()
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
			harjumusTehtud INTEGER DEFAULT 0,
			harjumuseStriik INTEGER DEFAULT 0,
			viimatiTehtud DATE,
			viimatiResetitud DATE
		);`
	_, err = db.Exec(tabel)
	if err != nil {
		log.Fatal(err)
	}

	statement, err := db.Prepare("INSERT INTO userdata (harjumuseNimi) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(harjumuseNimi)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	fmt.Println()
	fmt.Println("******************************")
	fmt.Println("Harjumus lisatud.")
	fmt.Println("******************************") // siin on 30 "*"
	fmt.Println()

	time.Sleep(500 * time.Millisecond)

	tervitaja()
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
	tervitaja()
}

func margiTehtuks() {
	var valik string
	var striik int64
	tana := time.Now().Truncate(24 * time.Hour)
	tanaString := tana.Format("2006-01-02")

	fmt.Print("Sisestage harjumuse ID: ")
	fmt.Scanln(&valik)

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT harjumuseStriik, viimatiTehtud FROM userdata")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	rows.Scan(&striik)
	striik += 1

	query := ("UPDATE userdata SET harjumusTehtud = ?, harjumuseStriik = ?, viimatiTehtud = ? WHERE id = ?;")

	db.Exec(query, 1, striik, tanaString, valik)

	time.Sleep(1 * time.Second)

	fmt.Println()
	fmt.Println("******************************")
	fmt.Println("Harjumus uuendatud.")
	fmt.Println("******************************")
	fmt.Println()
	time.Sleep(500 * time.Millisecond)

	tervitaja()
}

func harjumusteResetija(tana time.Time) error {
	tanaString := tana.Format("2006-01-02")

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE userdata
		SET harjumusTehtud = 0
			viimatiResetitud = ?
		WHERE viimatiResetitud IS NULL
			OR viimatiResetitud != ?
		`, tanaString, tanaString)

	return err
}

func striikideResetija(tana time.Time) error {
	kaksPaevaTagasi := tana.AddDate(0, 0, -2).Format("2006-01-02")

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
		UPDATE userdata
		SET harjumuseStriik = 0
		WHERE harjumusTehtud = 0
			AND (viimatiTehtud IS NULL OR viimatiTehtud < ?)
		`, kaksPaevaTagasi)

	return err
}

func dailyMaintenance() error {
	tana := time.Now().Truncate(24 * time.Hour)

	if err := harjumusteResetija(tana); err != nil {
		return fmt.Errorf("harjumuste resetija: %w", err)
	}
	if err := striikideResetija(tana); err != nil {
		return fmt.Errorf("striikide resetija: %w", err)
	}

	return nil
}

func easterEgg() {
	fmt.Println(`⠀⠀⠀⣴⣾⣿⣿⣶⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⢸⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠈⢿⣿⣿⣿⣿⠏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠈⣉⣩⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣼⣿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⢀⣼⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⢀⣾⣿⣿⣿⣿⣿⣿⣷⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⢠⣾⣿⣿⠉⣿⣿⣿⣿⣿⡄⠀⢀⣠⣤⣤⣀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠙⣿⣿⣧⣿⣿⣿⣿⣿⡇⢠⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠈⠻⣿⣿⣿⣿⣿⣿⣷⠸⣿⣿⣿⣿⣿⡿⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠘⠿⢿⣿⣿⣿⣿⡄⠙⠻⠿⠿⠛⠁⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⡟⣩⣝⢿⠀⠀⣠⣶⣶⣦⡀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣷⡝⣿⣦⣠⣾⣿⣿⣿⣿⣷⡀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣿⣿⣮⢻⣿⠟⣿⣿⣿⣿⣿⣷⡀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⡇⠀⠀⠻⠿⠻⣿⣿⣿⣿⣦⡀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⠇⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⣿⡆⠀⠀
⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⠀⠀⠀⠀⠀⠀⣠⣾⣿⣿⣿⣿⠇⠀⠀
⠀⠀⠀⠀⠀⠀⢸⣿⣿⡿⠀⠀⠀⢀⣴⣿⣿⣿⣿⣟⣋⣁⣀⣀⠀
⠀⠀⠀⠀⠀⠀⠹⣿⣿⠇⠀⠀⠀⠸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠇`)
	fmt.Println("see on Kari süü!!")
}

func main() {
	dailyMaintenance()
	tervitaja()
}
