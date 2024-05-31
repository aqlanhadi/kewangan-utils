package extractor

import (
	"fmt"
	u "mysimpan/statements/utils"
	"os"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/shopspring/decimal"
)

func ExtractFromMAE(file *os.File, fileReader *pdf.Reader) {

	// initialize structs
	Transactions = []Transaction{}
	ParsedData = Data{}

	statement_content := ""

	var count_since_main_record int
	var beginning_balance decimal.Decimal
	var ending_balance decimal.Decimal
	var total_debit decimal.Decimal
	var total_credit decimal.Decimal
	var in_footer bool = false

	transaction := new(Transaction)

	// extract date
	extractDate(file.Name())

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
				found, value := mae_getBeginningBalanceFromStatement(&rowContent)

				if found {
					beginning_balance = value
					ParsedData.SetStartingBalance(value)
				}
			}
			mae_parseMainRecord(&count_since_main_record, transaction, &rowContent)
			mae_parseTransactionDescription(&count_since_main_record, transaction, &rowContent)
			ending_found := mae_parseStatementEnd(transaction, &rowContent)

			if (in_footer || ending_found) {
				in_footer = true
				mae_populateChecksVariable(&rowContent, &ending_balance, &total_debit, &total_credit)
			}

			statement_content = statement_content + "\n" + rowContent
		}
	}

	// return Transactions

	// loop and do calculation to compare with parsed (show correct extraction)
	var calculated_debit, calculated_credit decimal.Decimal
	current_balance := beginning_balance

	// fmt.Println("------------------------------------")
	// fmt.Println("Beginning Balance\t: ", beginning_balance)
	// fmt.Println("------------------------------------")

	for _, item := range Transactions {
		current_balance = current_balance.Add(item.Amount)

		if item.Amount.Cmp(decimal.Zero) > 0 {
			calculated_credit = calculated_credit.Add(item.Amount)
			// fmt.Print("+" + item.Amount.StringFixed(2) + "\t | ")
			// fmt.Print(current_balance.StringFixed(2) + "\t | ")
			// fmt.Print(calculated_credit.StringFixed(2) + "\t | \t\t")

			// fmt.Print(item)
			// fmt.Println()
		} else {
			calculated_debit = calculated_debit.Add(item.Amount)
			// fmt.Print(item.Amount.StringFixed(2) + "\t | ")
			// fmt.Print(current_balance.StringFixed(2) + "\t | ")
			// fmt.Print("\t\t | " + calculated_debit.StringFixed(2) + "\t")

			// fmt.Print(item)
			// fmt.Println()

		}
		
	}

	ParsedData.AddTransactions(Transactions)

	ParsedData.SetTotalCredit(calculated_credit)
	ParsedData.SetTotalDebit(calculated_debit)
	ParsedData.SetEndingBalance(current_balance)

	ParsedData.SetParsedEndingBalance(ending_balance)

	// tests
	// fmt.Println("------------------------------------")
	// fmt.Println("Beginning Balance\t: ", beginning_balance)
	// fmt.Println("------------------------------------")
	// fmt.Println("DEBIT")
	// fmt.Println("\tCalculated\t: ", calculated_debit.StringFixed(2))
	// fmt.Println("\tParsed\t\t: ", total_debit.StringFixed(2))
	// fmt.Println("------------------------------------")
	// fmt.Println("CREDIT")
	// fmt.Println("\tCalculated\t: ", calculated_credit.StringFixed(2))
	// fmt.Println("\tParsed\t\t: ", total_credit.StringFixed(2))
	// fmt.Println("------------------------------------")
	// fmt.Println("ENDING")
	// fmt.Println("\tCalculated\t: ", current_balance.StringFixed(2))
	// fmt.Println("\tParsed\t\t: ", ending_balance.StringFixed(2))
	// fmt.Println("------------------------------------")


}

func extractDate(fileName string) {
	date_pattern, _ := regexp.Compile(u.Cfg.AccountTypeRegex.MBB_MAE_REGEX)
	date_match := date_pattern.FindStringSubmatch(fileName)

	ParsedData.SetYearAndMonth(date_match[1], date_match[2])
}

func testEquality(calc *decimal.Decimal, extracted *decimal.Decimal) (bool) {

	calc_abs := calc.Abs()
	extr_abs := extracted.Abs()

	return calc_abs.Equal(extr_abs)
}

func mae_parseMainRecord(csmr *int, transaction *Transaction, line *string) {

	regex_main_record_pattern, _ := regexp.Compile(`(\d{1,2}\/\d{1,2})\s+([A-Z0-9a-z \/.\-\*]+)\s+(\.\d{2}|\d{1,3}(,\d{3})*(\.\d{2})?)([-+])\s+(\.\d{2}|\d{1,3}(,\d{3})*(\.\d{2})?)(DR)?`)
	main_record_match := regex_main_record_pattern.FindStringSubmatch(*line)

	// fmt.Println(*line)
	if main_record_match != nil {
		// fmt.Println("[MAIN] > ", main_record_match)
		// for i, item := range main_record_match {
		// 	fmt.Printf("%d = %s\n", i, item)
		// }

		if *csmr > 0 {
			// save into struct
			Transactions = append(Transactions, *transaction)
		}

		*csmr = 0

		var amount_string string = strings.ReplaceAll(main_record_match[3], ",", "")
		var balance_string string = strings.ReplaceAll(main_record_match[7], ",", "")

		if main_record_match[6] == "-" {
			amount_string = "-" + amount_string
		}

		if main_record_match[10] == "DR" {
			balance_string = "-" + balance_string
		}
		
		amount, err := decimal.NewFromString(amount_string)
		if err != nil {
			fmt.Println(err)
		}

		balance, err := decimal.NewFromString(balance_string)
		if err != nil {
			panic("unable to parse float from " + balance_string)
		}
		 
		// fmt.Println(main_record_match)
		transaction.Origin = "MAE"
		transaction.Date = main_record_match[1]
		transaction.Action = main_record_match[2]
		transaction.Amount = amount
		transaction.Balance = balance


		*csmr++
	}
}

func mae_parseTransactionDescription(csmr *int, transaction *Transaction, line *string) {

	regex_transaction_description_pattern, _ := regexp.Compile(`[ ]{3,}([\S ]+)`)
	transaction_description_match := regex_transaction_description_pattern.FindStringSubmatch(strings.TrimRight(*line, " "))

	if transaction_description_match != nil {
		// fmt.Println("[DESC] > ", transaction_description_match)

		// for i, item := range transaction_description_match {
		// 	fmt.Printf("%d = %s\n", i, item)
		// }
		if *csmr == 1 {
			transaction.Beneficiary = transaction_description_match[1]
		}

		if *csmr == 2 {
			transaction.Description = transaction_description_match[1]
		}

		if *csmr == 3 {
			transaction.Method = transaction_description_match[1]
		}

		*csmr++
	}
}

func mae_parseStatementEnd(transaction *Transaction, line *string) (bool) {

	regex_statement_end_pattern, _ := regexp.Compile(`ENDING BALANCE : (\d+.\d+)`)
	statement_end_match := regex_statement_end_pattern.FindStringSubmatch(*line)

	if statement_end_match != nil {
		// fmt.Println("[END] > ", statement_end_match)
		// for i, item := range statement_end_match {
		// 	fmt.Printf("%d = %s\n", i, item)
		// }

		// if *csmr > 0 {
			Transactions = append(Transactions, *transaction)
		// }
		return true
	}

	return false

}

func mae_getBeginningBalanceFromStatement(line *string) (bool, decimal.Decimal) {
	regex_beginning_balance_pattern, _ := regexp.Compile(`BEGINNING BALANCE (\d{1,3}(,\d{3})*(\.\d{2})?)`)
	beginning_balance_match := regex_beginning_balance_pattern.FindStringSubmatch(*line)
	
	if beginning_balance_match != nil {
		value, _ := decimal.NewFromString(strings.ReplaceAll(beginning_balance_match[1], ",", ""))

		return true, value
	}

	return false, decimal.Decimal{}
}

func mae_populateChecksVariable(line *string, end *decimal.Decimal, debit *decimal.Decimal, credit *decimal.Decimal) {

	regex_ending_balance_pattern, _ := regexp.Compile(`ENDING BALANCE : (\d{1,3}(,\d{3})*(\.\d{2})?)`)
	regex_total_debit_pattern, _ := regexp.Compile(`TOTAL DEBIT : (\d{1,3}(,\d{3})*(\.\d{2})?)`)
	regex_total_credit_pattern, _ := regexp.Compile(`TOTAL CREDIT : (\d{1,3}(,\d{3})*(\.\d{2})?)`)

	ending_balance_match := regex_ending_balance_pattern.FindStringSubmatch(*line)
	total_debit_match := regex_total_debit_pattern.FindStringSubmatch(*line)
	total_credit_match := regex_total_credit_pattern.FindStringSubmatch(*line)

	if ending_balance_match != nil {
		value, err := decimal.NewFromString(strings.ReplaceAll(ending_balance_match[1], ",", ""))
		if err != nil {
			panic("unable to parse ending decimal from " + strings.ReplaceAll(ending_balance_match[1], ",", ""))
		}
		*end = value
	}

	if total_debit_match != nil {
		// fmt.Println(tot)
		value, err := decimal.NewFromString(strings.ReplaceAll(total_debit_match[1], ",", ""))
		if err != nil {
			panic("unable to parse debit decimal from " + strings.ReplaceAll(total_debit_match[1], ",", ""))
		}
		*debit = value
	}

	if total_credit_match != nil {
		value, err := decimal.NewFromString(strings.ReplaceAll(total_credit_match[1], ",", ""))
		if err != nil {
			panic("unable to parse credit decimal from " + strings.ReplaceAll(total_credit_match[1], ",", ""))
		}
		*credit = value
	}

}

