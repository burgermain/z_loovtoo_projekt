package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func helloFunc() {
	var valik string
	praeguneAeg := time.Now().Format("15:04:05")
	fmt.Println("Teretulemast, minu harjumuste jälgiasse! (palun vajuta nr 5 kui see on su esimene kord!)")
	fmt.Println("Kell on", praeguneAeg)
	fmt.Println("1. Lisa uus harjumus")
	fmt.Println("2. Näita kõiki harjumusi")
	fmt.Println("3. Märgi harjumus tehtuks")
	fmt.Println("4. Välju")
	fmt.Print("Vali tegevus: ")
	fmt.Scanln(&valik)
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
	case "5":
		esimeneKord()
	default:
		fmt.Println("Vigane valik, proovi uuesti.")
	}

}

func esimeneKord() {
	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	sql := `CREATE TABLE userdata (
			id INTEGER PRIMATY KEY,
			harjumuseNimi TEXT NOT NULL,
			harjumuseAeg TEXT NOT NULL,
			harjumusTehtud BOOL NOT NULL,
		);`

	db.Exec(sql)
	main()
}

func lisaHarjumus() {
	var harjumuseNimi string
	var harjumuseAeg string
	fmt.Scanln("Sisestage harjumuse nimi: ", &harjumuseNimi)
	fmt.Scanln("Sisestage harjumuse aeg: ", &harjumuseAeg)

	db, err := sql.Open("sqlite3", "userdata.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO userdata (harjumuseNimi, harjumuseAeg, harjumusTehtud) VALUES (?, ?, ?)", harjumuseNimi, harjumuseAeg, false)
	if err != nil {
		log.Fatal(err)
	}

	main()
}

func naitaHarjumusi() {

}

func margiTehtuks() {
	fmt.Println("3")
}

func main() {
	helloFunc()
}
