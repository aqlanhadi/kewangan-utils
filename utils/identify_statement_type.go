package utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/ledongthuc/pdf"
)

func IdentifyStatementAccount(cfg *Config, fileReader *pdf.Reader) (string, error) {

	p := fileReader.Page(1)
	if p.V.IsNull() {
		return "", errors.New("page is empty")
	}

	content, err := p.GetPlainText(nil)
	if err != nil {
		return "", err
	}

	v := reflect.ValueOf(cfg.AccountTypeIdentifiers)
	
	for i := 0; i < v.NumField(); i++ {

		key := v.Type().Field(i).Name
		
		if value, ok := v.Field(i).Interface().(string); ok {

			found := strings.Contains(content, value)

			if found {
				return key, nil
			}

		} else {
			return "", errors.New("unable to cast configuration value to string")
		}
	}

	return "", nil

}