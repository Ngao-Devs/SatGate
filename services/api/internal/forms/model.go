package forms

import "time"

type Form struct {
	OwnerID       string    `json:"-"`
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Domain        string    `json:"domain"`
	AmountSats    int64     `json:"amount_sats"`
	SatBalance    int64     `json:"sat_balance"`
	MinAmountSats int64     `json:"min_amount_sats"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateFormRequest struct {
	Name          string `json:"name"`
	Domain        string `json:"domain"`
	AmountSats    int64  `json:"amount_sats"`
	MinAmountSats int64  `json:"min_amount_sats"`
}

type UpdateFormRequest struct {
	Name          string `json:"name"`
	Domain        string `json:"domain"`
	AmountSats    int64  `json:"amount_sats"`
	MinAmountSats int64  `json:"min_amount_sats"`
}
