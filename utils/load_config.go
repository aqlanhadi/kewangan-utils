package utils

import (
	"os"

	yml "gopkg.in/yaml.v2"
)

type Config struct {
	AccountTypeIdentifiers struct {
		MBB_MAE string `yaml:"mbb_mae"`
		MBB_SAVINGS_I string `yaml:"mbb_savings_i"`
		MBB_MAYBANK_2_CREDIT_CARDS string `yaml:"mbb_maybank_2_credit_cards"`
	} `yaml:"account_type_string_identifiers"`
	AccountTypeRegex struct {
		MBB_MAE_REGEX string `yaml:"mbb_mae"`
	} `yaml:"account_type_file_regex"`
}

var Cfg Config

func LoadConfig() (*Config, error) {
	f, err := os.Open("config.yml")

	if err != nil {
		return nil, err
	}

	defer f.Close()

	decoder := yml.NewDecoder(f)
	err = decoder.Decode(&Cfg)

	if err != nil {
		return nil, err
	}

	return &Cfg, nil
}