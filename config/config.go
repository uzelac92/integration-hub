package config

import (
	"log"
	"os"
)

type Config struct {
	Port      string
	DbUrl     string
	WalletUrl string
	RgsUrl    string
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://rgs:rgs@localhost:5432/rgs?sslmode=disable"
	}

	walletUrl := os.Getenv("WALLET_URL")
	rgsUrl := os.Getenv("RGS_URL")
	if walletUrl == "" || rgsUrl == "" {
		log.Fatal("WALLET_URL and RGS_URL must be set")
	}

	return Config{
		DbUrl:     dbURL,
		WalletUrl: walletUrl,
		RgsUrl:    rgsUrl,
		Port:      port,
	}
}
