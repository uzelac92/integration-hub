package webhook

import (
	"bytes"
	"context"
	"integration-hub/internal/storage/db"
	"io"
	"log"
	"net/http"
	"time"
)

type Dispatcher struct {
	queries *db.Queries
	rgsURL  string
	client  *http.Client
}

func NewDispatcher(q *db.Queries, rgsURL string) *Dispatcher {
	return &Dispatcher{
		queries: q,
		rgsURL:  rgsURL,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (d *Dispatcher) Start() {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			d.processPending()
		}
	}()
}

func (d *Dispatcher) processPending() {
	ctx := context.Background()

	events, err := d.queries.GetDueWebhooks(ctx)
	if err != nil {
		log.Println("dispatcher: error getting due webhooks: ", err)
		return
	}

	for _, event := range events {
		d.deliver(ctx, event)
	}
}

func (d *Dispatcher) deliver(ctx context.Context, event db.GetDueWebhooksRow) {
	log.Println("dispatcher: delivering webhook: ", event.ID)

	request, err := http.NewRequest("POST", d.rgsURL, bytes.NewReader(event.Payload))
	if err != nil {
		log.Println("dispatcher: error creating request: ", err)
		_ = d.queries.MarkWebhookFailed(ctx, event.ID)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := d.client.Do(request)
	if err != nil {
		log.Println("dispatcher: error sending webhook: ", err)
		_ = d.queries.MarkWebhookFailed(ctx, event.ID)
		return
	}
	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			log.Println("dispatcher: error closing body: ", errClose)
		}
	}(response.Body)

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		_ = d.queries.MarkWebhookFailed(ctx, event.ID)
		log.Println("dispatcher: sent", event.ID)
	} else {
		log.Println("dispatcher: RGS status: ", response.Status)
		_ = d.queries.MarkWebhookFailed(ctx, event.ID)
	}
}
