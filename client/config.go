package client

import "errors"

type TwitterAcc struct {
	AdAccountName string
	AdAccountID   string
}

type Config struct {
	APIKey       string
	APISecret    string
	AccessToken  string
	AccessSecret string
	AdAccounts   []TwitterAcc
}

var ErrNotValidConfig = errors.New("config is not valid")

func (c Config) Validate() error {
	if c.APIKey == "" || c.APISecret == "" || c.AccessToken == "" || c.AccessSecret == "" {
		return ErrNotValidConfig
	}

	return nil
}
