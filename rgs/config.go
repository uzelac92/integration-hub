package main

import (
	"os"
	"rgs/observability"
)

type Config struct {
	DbUrl        string
	WalletUrl    string
	WalletSecret string
}

func LoadConfig() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		observability.Logger.Fatal("DATABASE_URL must be set")
	}

	walletUrl := os.Getenv("WALLET_URL")
	walletSecret := os.Getenv("WALLET_SECRET")
	if walletUrl == "" || walletSecret == "" {
		observability.Logger.Fatal("WALLET_URL and WALLET_SECRET must be set")
	}

	return Config{
		DbUrl:        dbURL,
		WalletUrl:    walletUrl,
		WalletSecret: walletSecret,
	}
}
