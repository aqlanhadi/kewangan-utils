package loader

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	d "mysimpan/statements/extractor"
	"strings"
	"time"
)

var AffectedRows int

func Load(db *sql.DB, d *d.Data) {
	AffectedRows = 0
	// json, jerr := json.MarshalIndent(*d, "", "  ")
	// if jerr != nil {
	// 	panic("error marshalling struct")
	// }


	// // fmt.Println(string(json))
	// panic("zzz")
	// delete all from similar source
	err := deleteTransactionsOriginatingFromSameSource(db, d.Source)
	if err != nil {
		panic(err)
	}

	// insert query
	for _, transaction := range d.Transactions {
		// fmt.Println(transaction)
		err := insertTransaction(db, d, transaction)
		if err != nil {
			panic(err)
		}
	}

	log.Printf("%d rows affected", AffectedRows)
}

func deleteTransactionsOriginatingFromSameSource(db *sql.DB, source string) error {
	query := `DELETE FROM myduit.transaction WHERE source = $1`

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancelFunc()

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil{
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}

	res, err := stmt.ExecContext(ctx, source)

	if err != nil {
		log.Printf("Error %s when deleting rows from transactions table", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}

	log.Printf("%d transactions from %s deleted", rows, source)
	return nil

}

func insertTransaction(db *sql.DB, d *d.Data, t d.Transaction) error {
	query := `
		INSERT INTO myduit.transaction(
			account,
			account_number,
			account_type,
			posting_date,
			date,
			action,
			beneficiary,
			method,
			amount,
			balance,
			source
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)`
	
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	stmt, err := db.PrepareContext(ctx, query)

	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()

	date_string := fmt.Sprintf("%s-%s-%s", d.Year, d.Month, strings.Split(t.Date, "/")[0])

	res, err := stmt.ExecContext(
		ctx, 
		d.Account,
		d.AccountNumber,
		d.AccountType,
		nil,
		date_string,
		t.Action,
		t.Beneficiary,
		t.Method,
		t.Amount,
		t.Balance,
		d.Source)
	
	if err != nil {
		log.Printf("Error %s when inserting row into transactions table", err)
		return err
	}

	

	rows, _ := res.RowsAffected()
	// if err != nil {
	// 	log.Printf("Error %s when finding rows affected", err)
	// 	return err
	// }

	AffectedRows += int(rows)

	// log.Printf("%d transactions created", rows)
	return nil
}