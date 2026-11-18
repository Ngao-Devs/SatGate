package messages

import "time"

type Message struct {
	ID          string     `json:"id"`
	FormID      string     `json:"form_id"`
	SenderName  string     `json:"sender_name"`
	SenderEmail string     `json:"sender_email"`
	Body        string     `json:"body"`
	Status      string     `json:"status"`
	AmountSats  int64      `json:"amount_sats"`
	CreatedAt   time.Time  `json:"created_at"`
	PaidAt      *time.Time `json:"paid_at"`
}
