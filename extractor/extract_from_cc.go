package extractor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/shopspring/decimal"
)

func ExtractFromMBBCC(fileReader *pdf.Reader) {

	// cc needs to do 2 validations transactions, one amex, one master

	statement_content := ""

	var count_since_main_record int
	var beginning_balance decimal.Decimal
	var ending_balance decimal.Decimal
	var total_debit decimal.Decimal
	var total_credit decimal.Decimal
	var card string
	var in_footer bool = false
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
			
			ending_found := mbb_cc_parseStatementEnd(transaction, &rowContent)

			// mbb_cc_parseTransactionDescription(&count_since_main_record, transaction, &rowContent)
			// end write if found ending

			if (in_footer || ending_found) {
				in_footer = true
				mbb_cc_populateChecksVariable(&rowContent, &ending_balance, &total_debit, &total_credit)
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

	// tests

	fmt.Println("beginning balance", beginning_balance)

	all_okay := testEquality(&calculated_debit, &total_debit) && testEquality(&calculated_credit, &total_credit) && testEquality(&current_balance, &ending_balance)

	if all_okay {
		fmt.Println("all equality checks passed.")
		fmt.Println("generating json.")
	} else {
		fmt.Println("some checks failed.")
	}

}

func mbb_cc_flipWrite(write *bool, card *string, line *string) {
	regex_master_card_pattern, _ := regexp.Compile(`MAYBANK 2 PLAT MASTERCARD    :    5239 4503 2205 3215`)
	regex_card_ends, _ := regexp.Compile(`TOTAL CREDIT THIS MONTH`)
	regex_amex_card_pattern, _ := regexp.Compile(`MAYBANK 2 PLATINUM AMEX    :    3791 861275 63716`)

	master_match := regex_master_card_pattern.FindStringSubmatch(*line)
	amex_match := regex_amex_card_pattern.FindStringSubmatch(*line)
	statement_ends := regex_card_ends.FindStringSubmatch(*line)

	if statement_ends != nil {
		*write = false
		fmt.Println(`> `, *line)
		fmt.Println(`Statement ends`)

	}

	if master_match != nil {
		*write = true
		*card = "MASTERCARD"
		if *write {
			fmt.Println(`In master mode`)
		}
	}

	if amex_match != nil {
		*write = true
		*card = "AMEX"
		if *write {
			fmt.Println(`In amex mode`)
		}
	}
}

func mbb_cc_parseMainRecord(csmr *int, transaction *Transaction, card *string, line *string) {

	regex_main_record_pattern, _ := regexp.Compile(`(\d{2}\/\d{2})\s+(\d{2}\/\d{2})\s+([A-Z0-9\S ]+?)\s+((\d{1,3}(,\d{3})*(\.\d+))(CR)?)\s*$`)
	main_record_match := regex_main_record_pattern.FindStringSubmatch(*line)

	fmt.Println(*line)
	if main_record_match != nil {
		// fmt.Println("[MAIN] > ", main_record_match)
		// for i, item := range main_record_match {
		// 	fmt.Printf("%d = %s\n", i, item)
		// }
		// panic("h")

		if *csmr > 0 {
			// save into struct
			Transactions = append(Transactions, *transaction)
		}

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
		 
		// fmt.Println(main_record_match)
		transaction.Origin = "CREDIT CARD"
		transaction.Date = main_record_match[2]
		transaction.Action = main_record_match[3]
		transaction.Amount = amount
		transaction.Method = *card

		*csmr++
	}
}

func mbb_cc_parseTransactionDescription(csmr *int, transaction *Transaction, line *string) {

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

func mbb_cc_parseStatementEnd(transaction *Transaction, line *string) (bool) {

	regex_statement_end_pattern, _ := regexp.Compile(`ENDING BALANCE : (\d+.\d+)`)
	statement_end_match := regex_statement_end_pattern.FindStringSubmatch(*line)

	if statement_end_match != nil {
		// fmt.Println("[END] > ", statement_end_match)
		// for i, item := range statement_end_match {
		// 	fmt.Printf("%d = %s\n", i, item)
		// }


		Transactions = append(Transactions, *transaction)

		return true
	}

	return false

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

func mbb_cc_populateChecksVariable(line *string, end *decimal.Decimal, debit *decimal.Decimal, credit *decimal.Decimal) {

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

