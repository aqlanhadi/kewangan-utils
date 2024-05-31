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
		mbb_mae string `yaml:"mbb_mae"`
		mbb_casa string `yaml:"mbb_casa_i"`
		mbb_psa_i string `yaml:"mbb_psa_i"`
	} `yaml:"account_type_file_regex"`
	AccountTypeMap struct {
		CASA []string `yaml:"casa"`
	} `yaml:"account_type_map"`
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