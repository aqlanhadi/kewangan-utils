package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"

	v "github.com/spf13/viper"
)

func IdentifyAccountTypeFromFileName(fileName string) (string, string, error) {

	v.SetConfigName("config")
	wd, _ := os.Getwd()
	v.AddConfigPath(wd)

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	type_map := v.GetStringMapStringSlice("account_type_map")

	for fileKey, fileRegex := range v.GetStringMapString("account_type_file_regex") {
		file_pattern, _ := regexp.Compile(fileRegex)
		match := file_pattern.FindStringSubmatch(fileName)

		if match != nil {
			for supertype, vals := range type_map {
				// fmt.Println(supertype, vals)
				if (slices.Contains(vals, fileKey)) {
					return supertype, fileKey, nil
				}
			}
		}
	}

	return "", "", errors.New("unknown account type")
}