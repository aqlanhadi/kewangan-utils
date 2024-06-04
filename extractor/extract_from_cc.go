package extractor

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/shopspring/decimal"

	v "github.com/spf13/viper"
)

type CCRegex struct {
	MasterStatementStartRegex string
	AmexStatementStartRegex string
	StatementEndRegex string
	StatementDebitRegex string
}

var Regexes CCRegex

func getRegexDefinitions() {
	v.SetConfigName("config")
	wd, _ := os.Getwd()
	v.AddConfigPath(wd)

	Regexes.AmexStatementStartRegex = v.GetString("helper_regex.mbb_2_cc_amex_start_pattern")
	Regexes.MasterStatementStartRegex = v.GetString("helper_regex.mbb_2_cc_mastercard_start_pattern")
	Regexes.StatementEndRegex = v.GetString("helper_regex.mbb_2_cc_statement_credit_pattern")
	Regexes.StatementDebitRegex = v.GetString("helper_regex.mbb_2_cc_statement_debit_pattern")

}

var TotalParsedCredit decimal.Decimal
var TotalParsedDebit decimal.Decimal

func ExtractFromMBBCC(file *os.File, fileReader *pdf.Reader) {

	getRegexDefinitions()
	// cc needs to do 2 validations transactions, one amex, one master
	// initialize structs
	Transactions = []Transaction{}

	statement_content := ""

	var count_since_main_record int
	var beginning_balance decimal.Decimal
	
	// var total_parsed_debit decimal.Decimal
	var card string
	var write bool

	transaction := new(Transaction)

	// loop pages
	for pageIndex := 1; pageIndex <= fileReader.NumPage(); pageIndex++ {
		p := fileReader.Page(pageIndex)

		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()

		// loop rows
		for _, row := range rows {

			var rowContent string
			for _, word := range row.Content {
				rowContent = rowContent + " " + word.S
			}

			if pageIndex == 1 {
				found, value := mbb_cc_getBeginningBalanceFromStatement(&rowContent)

				if found {
					beginning_balance = value
				}
			}

			// if found start write
			mbb_cc_flipWrite(&write, &card, &rowContent)

			if write {
				mbb_cc_parseMainRecord(&count_since_main_record, transaction, &card, &rowContent)
			}

			statement_content = statement_content + "\n" + rowContent
		}
	}

	// fmt.Println(statement_content)

	// loop and do calculation to compare with parsed (show correct extraction)
	var calculated_debit, calculated_credit decimal.Decimal
	current_balance := beginning_balance

	for _, item := range Transactions {

		fmt.Println(item)
		current_balance = current_balance.Add(item.Amount)

		if item.Amount.Cmp(decimal.Zero) > 0 {
			calculated_credit = calculated_credit.Add(item.Amount)
		} else {
			calculated_debit = calculated_debit.Add(item.Amount)
		}
		
	}

	fmt.Printf("TOTAL CREDIT\t%s\n", calculated_credit.StringFixed(2))
	fmt.Printf("PARSED CREDIT\t%s\n", TotalParsedCredit.StringFixed(2))

	fmt.Printf("TOTAL DEBIT\t%s\n", calculated_debit.StringFixed(2))
	fmt.Printf("PARSED DEBIT\t%s\n", TotalParsedDebit.StringFixed(2))


	ParsedData.AddTransactions(Transactions)

	ParsedData.SetTotalCredit(calculated_credit)
	ParsedData.SetTotalDebit(calculated_debit)
	// ParsedData.SetEndingBalance(current_balance)

	// ParsedData.SetParsedEndingBalance(EndingBalance)

	// tests

	// fmt.Println("beginning balance", beginning_balance)

	// all_okay := testEquality(&calculated_debit, &total_parsed_debit) && testEquality(&calculated_credit, &TotalParsedCredit) && testEquality(&current_balance, &TotalParsedCredit)

	// if all_okay {
	// 	fmt.Println("all equality checks passed.")
	// 	fmt.Println("generating json.")
	// } else {
	// 	fmt.Println("some checks failed.")
	// }

}

func mbb_cc_flipWrite(write *bool, card *string, line *string) {
	// fmt.Println("checking line\t", *line)
	regex_card_ends, _ := regexp.Compile(Regexes.StatementEndRegex)
	regex_master_card_pattern, _ := regexp.Compile(Regexes.MasterStatementStartRegex)
	regex_amex_card_pattern, _ := regexp.Compile(Regexes.AmexStatementStartRegex)
	regex_debit_pattern, _ := regexp.Compile(Regexes.StatementDebitRegex)

	master_match := regex_master_card_pattern.FindStringSubmatch(*line)
	amex_match := regex_amex_card_pattern.FindStringSubmatch(*line)
	statement_ends := regex_card_ends.FindStringSubmatch(*line)
	debit_match := regex_debit_pattern.FindStringSubmatch(*line)

	if statement_ends != nil {
		*write = false
		// fmt.Println(`> `, *line)
		// fmt.Println(statement_ends)

		parsed_cred, err := decimal.NewFromString(strings.ReplaceAll(statement_ends[1], ",", ""))
		if err != nil {
			log.Printf("error parsing value %s to decimal: %s", statement_ends[1], err)
		}
		// fmt.Println(parsed_cred.StringFixed(2))
		TotalParsedCredit = TotalParsedCredit.Add(parsed_cred)
		// fmt.Println(EndingBalance.Add(parsed_cred))
		// ParsedData.SetParsedEndingBalance(parsed_cred)
		// fmt.Println(`Statement ends`)
		// fmt.Println(ParsedData.EndingBalance)
		
	}

	if debit_match != nil {
		parsed_deb, err := decimal.NewFromString(strings.ReplaceAll(debit_match[1], ",", ""))
		if err != nil {
			log.Printf("error parsing value %s to decimal: %s", statement_ends[1], err)
		}
		TotalParsedDebit = TotalParsedDebit.Add(parsed_deb)
	}

	if master_match != nil {
		*write = true
		*card = "MASTERCARD"
		// if *write {
		// 	fmt.Println(`In master mode`)
		// }
	}

	if amex_match != nil {
		*write = true
		*card = "AMEX"
		// if *write {
		// 	fmt.Println(`In amex mode`)
		// }
	}
}

func mbb_cc_parseMainRecord(csmr *int, transaction *Transaction, card *string, line *string) {

	regex_main_record_pattern, _ := regexp.Compile(`(\d{2}\/\d{2})\s+(\d{2}\/\d{2})\s+([A-Z0-9\S ]+?)\s+((\d{1,3}(,\d{3})*(\.\d+))(CR)?)\s*$`)
	main_record_match := regex_main_record_pattern.FindStringSubmatch(*line)

	// fmt.Println(*line)
	if main_record_match != nil {
		// fmt.Println("[MAIN] > ", main_record_match)
		// for i, item := range main_record_match {
		// 	fmt.Printf("%d = %s\n", i, item)
		// }
		// panic("h")

		// if *csmr > 0 {
		// 	// save into struct
		// 	Transactions = append(Transactions, *transaction)
		// }

		*csmr = 0

		var amount_string string = strings.ReplaceAll(main_record_match[4], ",", "")

		if main_record_match[8] == "CR" {
			amount_string = strings.ReplaceAll(amount_string, "CR", "")
		} else {
			amount_string = "-" + amount_string
		}

		amount, err := decimal.NewFromString(amount_string)
		if err != nil {
			fmt.Println(amount)
			fmt.Println(err)
			panic("conversion error")
		}
		 
		// fmt.Println("date> ", main_record_match[2])
		// if transaction does not fall in this month, set it to first day
		// parsed_date, err := time.Parse("02/01", main_record_match[2])
		if ParsedData.Month != strings.Split(main_record_match[2], "/")[1] {
			transaction.Date = fmt.Sprintf(`%s/%s`, "01", ParsedData.Month)
		} else {
			transaction.Date = main_record_match[2]
		}
		// fmt.Println(transaction.Date)


		transaction.Origin = "CREDIT CARD"
		transaction.Action = main_record_match[3]
		transaction.Amount = amount
		transaction.Method = *card

		*csmr++

		Transactions = append(Transactions, *transaction)

	}
}

func mbb_cc_getBeginningBalanceFromStatement(line *string) (bool, decimal.Decimal) {
	regex_beginning_balance_pattern, _ := regexp.Compile(`YOUR PREVIOUS STATEMENT BALANCE (\d{1,3}(,\d{3})*(\.\d{2})?)`)
	beginning_balance_match := regex_beginning_balance_pattern.FindStringSubmatch(*line)
	
	if beginning_balance_match != nil {
		// fmt.Println("[BEGINNING BALANCE] > ", beginning_balance_match)
		value, _ := decimal.NewFromString(strings.ReplaceAll(beginning_balance_match[1], ",", ""))

		return true, value
	}

	return false, decimal.Decimal{}
}