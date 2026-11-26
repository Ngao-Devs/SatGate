package webhooks

type LightningWebhook struct {
	InvoiceID   string  `json:"invoice_id"`
	PaymentHash *string `json:"payment_hash"`
	Preimage    *string `json:"preimage"`
	Status      string  `json:"status"`
}
