package extractor

import (
	"fmt"
	"os"

	"github.com/ledongthuc/pdf"
)

var Transactions []Transaction
var ParsedData Data

func Extract(f *os.File, r *pdf.Reader, acc_type string) (Data) {
	
	if acc_type == "MBB_MAE" {
		fmt.Println("statement is from MAE account type. extracting")
		ExtractFromMAE(f, r)
	}

	if acc_type == "MBB_SAVINGS_I" {
		fmt.Println("statement is from CASA account type. extracting")
		ExtractFromCASA(r)
	}

	if acc_type == "MBB_MAYBANK_2_CREDIT_CARDS" {
		fmt.Println("statement is from CC Account type. extracting")
		ExtractFromMBBCC((r))
	}

	return ParsedData
}	