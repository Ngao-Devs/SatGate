package forms

import "time"

type Form struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Domain     string    `json:"domain"`
	AmountSats int64     `json:"amount_sats"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateFormRequest struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	AmountSats int64  `json:"amount_sats"`
}

type UpdateFormRequest struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	AmountSats int64  `json:"amount_sats"`
}
