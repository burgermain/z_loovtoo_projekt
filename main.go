package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
)

func helloFunc() {
	var valik string
	praeguneAeg := time.Now().Format("15:04:05")

	fmt.Println("Teretulemast, minu harjumuste jälgiasse!")
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
		fmt.Println("Vigane valik, proovi uuesti.")
	}
}

func lisaHarjumus() {
	var harjumuseNimi string
	var harjumuseAeg string

	fmt.Print("Sisestage harjumuse nimi: ")
	fmt.Scanln(&harjumuseNimi)
	fmt.Print("Sisestage harjumuse aeg: ")
	fmt.Scanln(&harjumuseAeg)

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tabel := `CREATE TABLE IF NOT EXISTS userdata (
			id INTEGER PRIMATY KEY,
			harjumuseNimi TEXT NOT NULL,
			harjumuseAeg TEXT NOT NULL,
			harjumusTehtud BOOL NOT NULL
		);`
	_, err = db.Exec(tabel)
	if err != nil {
		log.Fatal(err)
	}

	statement, err := db.Prepare("INSERT INTO userdata (harjumuseNimi, harjumuseAeg, harjumusTehtud) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(harjumuseNimi, harjumuseAeg, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("******************************") //siin on 30 "*"

	time.Sleep(time.Duration(time.Duration.Seconds(3)))
	main()
}

func naitaHarjumusi() {

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, harjumuseNimi, harjumuseAeg, harjumusTehtud FROM userdata")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	tabel := tablewriter.NewWriter(os.Stdout)
	tabel.Header([]string{"ID", "Nimi", "Aeg", "Staatus"})

	for rows.Next() {
		var id, nimi, aeg, staatus string
		rows.Scan(&id, &nimi, &aeg, &staatus)
		tabel.Append([]string{
			id,
			nimi,
			aeg,
			staatus,
		})
	}

	tabel.Render()

	time.Sleep(time.Duration(time.Duration.Seconds(3)))
	main()
}

func margiTehtuks() {
	fmt.Println("3")
}

func main() {
	helloFunc()
}
