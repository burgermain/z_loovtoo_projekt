package main

import (
	"fmt"
)

func helloFunc() {
	fmt.Println("Teretulemast, minu harjumuste jälgiasse!")
	fmt.Println("1. Lisa uus harjumus")
	fmt.Println("2. Näita kõiki harjumusi")
	fmt.Println("3. Märgi harjumus tehtuks")
	fmt.Println("4. Välju")
	fmt.Print("Vali tegevus: ")
	var valik string
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
	default:
		fmt.Println("Vigane valik, proovi uuesti.")
	}

}

func lisaHarjumus() {
	fmt.Println("1")
}

func naitaHarjumusi() {
	fmt.Println("2")
}

func margiTehtuks() {
	fmt.Println("3")
}

func main() {
	helloFunc()
}
