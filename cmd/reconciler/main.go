package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"integration-hub/config"
	"integration-hub/internal/storage/db"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OperatorTx struct {
	RefID        string `json:"refId"`
	PlayerID     string `json:"playerId"`
	Type         string `json:"type"`
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	BalanceAfter int64  `json:"balanceAfter"`
}

func main() {
	cfg := config.LoadConfig()

	pool, err := pgxpool.New(context.Background(), cfg.DbUrl)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}

	q := db.New(pool)
	ctx := context.Background()

	hubTxs, err := q.ListHubTransactions(ctx)
	if err != nil {
		log.Fatalf("query failed: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(cfg.WalletUrl + "/v2/reconciliation")
	if err != nil {
		log.Fatalf("fetch operator reconciliation failed: %v", err)
	}
	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			log.Println("failed to close body", errClose)
		}
	}(resp.Body)

	var opResp struct {
		Transactions []OperatorTx `json:"transactions"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &opResp); err != nil {
		log.Fatalf("decode operator response failed: %v", err)
	}

	opMap := map[string]OperatorTx{}
	for _, t := range opResp.Transactions {
		opMap[t.RefID] = t
	}

	mismatches := [][]string{
		{"refId", "playerId", "type", "hubAmount", "opAmount", "hubStatus", "opBalance"},
	}

	for _, htx := range hubTxs {
		otx, ok := opMap[htx.RefID]
		if !ok {
			mismatches = append(mismatches, []string{
				htx.RefID, htx.PlayerID, htx.Type,
				fmt.Sprint(htx.AmountCents),
				"MISSING",
				htx.OperatorStatus,
				"MISSING",
			})
			continue
		}

		if htx.AmountCents != otx.Amount || htx.OperatorBalance != otx.BalanceAfter {
			mismatches = append(mismatches, []string{
				htx.RefID, htx.PlayerID, htx.Type,
				fmt.Sprint(htx.AmountCents),
				fmt.Sprint(otx.Amount),
				htx.OperatorStatus,
				fmt.Sprint(otx.BalanceAfter),
			})
		}
	}

	f, err := os.Create("reconciliation_result.csv")
	if err != nil {
		log.Fatalf("cannot create csv: %v", err)
	}
	defer func(f *os.File) {
		errClose := f.Close()
		if errClose != nil {
			log.Fatalf("cannot close file: %v", errClose)
		}
	}(f)

	w := csv.NewWriter(f)
	_ = w.WriteAll(mismatches)
	w.Flush()

	if len(mismatches) > 1 {
		fmt.Println("MISMATCHES FOUND")
		os.Exit(1)
	}

	fmt.Println("NO MISMATCHES")
}
