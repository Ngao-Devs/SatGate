package invoices

import "time"

type CreateInvoiceRequest struct {
	FormID      string `json:"form_id"`
	SenderName  string `json:"sender_name"`
	SenderEmail string `json:"sender_email"`
	Body        string `json:"body"`
}

type InvoiceResponse struct {
	InvoiceID      string    `json:"invoice_id"`
	MessageID      string    `json:"message_id"`
	PaymentRequest string    `json:"payment_request"`
	AmountSats     int64     `json:"amount_sats"`
	Status         string    `json:"status"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type InvoiceStatusResponse struct {
	InvoiceID string     `json:"invoice_id"`
	MessageID string     `json:"message_id"`
	Status    string     `json:"status"`
	PaidAt    *time.Time `json:"paid_at"`
}
