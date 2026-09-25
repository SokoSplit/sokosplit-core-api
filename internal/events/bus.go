package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// ReleaseFundsEvent is emitted when a split list needs its escrowed funds
// released. sokosplit-wallet-service listens for this (via a webhook, in
// this simple implementation) and executes the corresponding on-chain
// contract call.
type ReleaseFundsEvent struct {
	SplitID string `json:"split_id"`
	// ConfirmerSecret is only present for manually-confirmed splits where
	// core-api holds a service-managed signing key. Splits confirmed via a
	// user's own wallet (e.g. Freighter) are released client-side instead.
	ConfirmerSecret string `json:"confirmer_secret,omitempty"`
}

// Bus is a minimal internal event bus. In this skeleton it's a thin HTTP
// webhook caller; swap the Publish implementation for a message queue
// (e.g. NATS, SQS) as the project grows past a single wallet-service
// instance.
type Bus struct {
	WalletServiceURL string
	WebhookSecret    string
}

func NewBus() *Bus {
	return &Bus{
		WalletServiceURL: os.Getenv("WALLET_SERVICE_URL"),
		WebhookSecret:    os.Getenv("WALLET_SERVICE_WEBHOOK_SECRET"),
	}
}

// PublishReleaseFunds notifies the wallet service that a split list is
// ready to be released.
func (b *Bus) PublishReleaseFunds(event ReleaseFundsEvent) error {
	if b.WalletServiceURL == "" {
		return fmt.Errorf("WALLET_SERVICE_URL is not set")
	}
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event: %w", err)
	}
	req, err := http.NewRequest(
		http.MethodPost,
		b.WalletServiceURL+"/events/release",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Secret", b.WebhookSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling wallet service: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}
	return nil
}
