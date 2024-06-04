package main

import (
	"fmt"
	"log"
	"mysimpan/statements/extractor"
	"mysimpan/statements/loader"
	u "mysimpan/statements/utils"
	"os"

	"database/sql"

	"github.com/ledongthuc/pdf"
	_ "github.com/lib/pq"
)

func main() {


	// Connect to database
	connStr := "postgresql://postgres:password@localhost:5432/myduit_test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dir := "./test/PSA-i/"

	files, _ := os.ReadDir(dir)

	for _, file := range files {

		accSupertype, accSubtype, err := u.IdentifyAccountTypeFromFileName(file.Name())
		if err != nil {
			fmt.Println("file type is unknown, skipping...")
			continue
		}

		fmt.Printf("extracting %s (%s) as %s\n", file.Name(), accSubtype, accSupertype)

		f, r, err := pdf.Open(dir + file.Name())
		if err != nil {
			panic(err)
		}

		defer f.Close()

		d := extractor.Extract(f, r, accSupertype, accSubtype)

		if (d.EndingBalanceMatches()) {
			fmt.Println("Parsing successful for " + file.Name())
			fmt.Println("Adding to DB")

			loader.Load(db, &d)
			// panic("stop")
			
		}

		// out_file_name := fmt.Sprintf("%s_%s_%s.json", accType, d.Year, d.Month)

		// b, _ := json.Marshal(d)

		// os.WriteFile("./test/MAE/out/" + out_file_name, b, 0644)
	}

	// statement_file_path := "0398121207523300_20240428.pdf" // cc
	// statement_file_path := "514169996465_20240229.pdf" // mae
	// statement_file_path := "514169996465_20240131.pdf"
	// statement_file_path := "114013-315457_20240430.pdf" // savings

	

	// fmt.Println(string(b))

	// // fmt.Println(content)
	// fmt.Println()
	// fmt.Println(acc_type)

	
}

