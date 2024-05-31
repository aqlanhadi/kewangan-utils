package loader

import (
	"database/sql"
	"encoding/json"
	"fmt"
	d "mysimpan/statements/extractor"
)

func Load(db *sql.DB, d *d.Data) {

	json, err := json.MarshalIndent(*d, "", "  ")
	if err != nil {
		panic("error marshalling struct")
	}


	fmt.Println(string(json))
}