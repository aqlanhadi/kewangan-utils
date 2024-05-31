package extractor

import (
	"os"
	"path"
	"strings"

	"github.com/ledongthuc/pdf"
)

var Transactions []Transaction
var ParsedData Data

func Extract(f *os.File, r *pdf.Reader, acc_supertype string, acc_subtype string) (Data) {

	ParsedData = Data{
		AccountType: acc_supertype,
		Account: acc_subtype,
		Source: path.Base(f.Name()),
		AccountNumber: strings.Split(path.Base(f.Name()), "_")[0],
	}

	ExtractDate(f.Name())

	if ParsedData.AccountType == "mbb_casa" {
		ExtractFromCASA(f, r)
	}

	// if acc_type == "MBB_SAVINGS_I" {
	// 	ExtractFromMAE(f, r)
	// }

	// if acc_type == "MBB_MAYBANK_2_CREDIT_CARDS" {
	// 	fmt.Println("statement is from CC Account type. extracting")
	// 	ExtractFromMBBCC((r))
	// }

	// fmt.Println(ParsedData.Month)
	// fmt.Println("\t" + ParsedData.StartingBalance.StringFixed(2))
	// if ParsedData.EndingBalanceMatches() {
	// 	fmt.Println("\t" + "CHECKS PASSED")
	// } else {
	// 	fmt.Println("\t" + "CHECKS FAILED")
	// }

	// fmt.Print(json.Marshal(ParsedData))
	// panic("z")

	return ParsedData
}	