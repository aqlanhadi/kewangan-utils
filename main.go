package main

import (
	"encoding/json"
	"fmt"
	"mysimpan/statements/extractor"
	u "mysimpan/statements/utils"
	"os"

	"github.com/ledongthuc/pdf"
)

func main() {
	
	cfg, err := u.LoadConfig()
	if err != nil {
		panic(err)
	}

	dir := "./test/CASA-i/"

	files, _ := os.ReadDir(dir)
	

	for _, file := range files {
		fmt.Println(file.Name())

		f, r, err := pdf.Open(dir + file.Name())
		if err != nil {
			panic(err)
		}

		defer f.Close()

		acc_type, err := u.IdentifyStatementAccount(cfg, r)
		if err != nil {
			panic(err)
		}

		d := extractor.Extract(f, r, acc_type)

		if (d.EndingBalanceMatches()) {
			fmt.Println("Parsing successful for " + file.Name())
			fmt.Println("Adding to DB")
		}

		out_file_name := fmt.Sprintf("%s_%s_%s.json", acc_type, d.Year, d.Month)

		b, _ := json.Marshal(d)

		os.WriteFile("./test/MAE/out/" + out_file_name, b, 0644)
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

