package webhooks

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"satgate/internal/invoices"

	"github.com/google/uuid"
)

type ErrorKind string

const (
	ErrSignature         ErrorKind = "signature"
	ErrInvalidBody       ErrorKind = "invalid_body"
	ErrInvalidStatus     ErrorKind = "invalid_status"
	ErrInternal          ErrorKind = "internal"
	ErrNotFound          ErrorKind = "not_found"
	ErrPaymentHash       ErrorKind = "payment_hash"
	ErrDuplicateNotFound ErrorKind = "duplicate_not_found"
)

type ProcessingError struct {
	Kind ErrorKind
	Err  error
}

func (e *ProcessingError) Error() string {
	return e.Err.Error()
}

func (e *ProcessingError) Unwrap() error {
	return e.Err
}

type Service struct {
	repository     *Repository
	invoiceService *invoices.Service
	secret         string
}

func NewService(repository *Repository, invoiceService *invoices.Service, secret string) *Service {
	return &Service{
		repository:     repository,
		invoiceService: invoiceService,
		secret:         secret,
	}
}

func (s *Service) Process(
	ctx context.Context,
	headers http.Header,
	body []byte,
) (*invoices.InvoiceStatusResponse, error) {
	if err := verifySignature(headers, string(body), s.secret); err != nil {
		return nil, &ProcessingError{Kind: ErrSignature, Err: err}
	}

	var payload LightningWebhook
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, &ProcessingError{Kind: ErrInvalidBody, Err: err}
	}

	if payload.Status != "paid" && payload.Status != "settled" {
		return nil, &ProcessingError{
			Kind: ErrInvalidStatus,
			Err:  errors.New("webhook status is not paid"),
		}
	}

	eventID := headers.Get("x-satgate-event-id")
	if eventID == "" {
		eventID = uuid.New().String()
	}

	recorded, err := s.repository.RecordEvent(
		ctx,
		eventID,
		payload.InvoiceID,
		string(body),
		time.Now().UTC(),
	)
	if err != nil {
		return nil, &ProcessingError{Kind: ErrInternal, Err: err}
	}

	if !recorded {
		response, err := s.invoiceService.GetStatus(ctx, payload.InvoiceID)
		if err != nil {
			return nil, &ProcessingError{
				Kind: ErrDuplicateNotFound,
				Err:  errors.New("not found"),
			}
		}
		return response, nil
	}

	if payload.PaymentHash != nil {
		existingHash, err := s.repository.GetPaymentHash(ctx, payload.InvoiceID)
		if err == nil && existingHash != *payload.PaymentHash {
			return nil, &ProcessingError{
				Kind: ErrPaymentHash,
				Err:  errors.New("payment hash does not match invoice"),
			}
		}
	}

	if err := s.invoiceService.MarkPaid(ctx, payload.InvoiceID, payload.Preimage); err != nil {
		return nil, &ProcessingError{Kind: ErrNotFound, Err: err}
	}

	response, err := s.invoiceService.GetStatus(ctx, payload.InvoiceID)
	if err != nil {
		return nil, &ProcessingError{Kind: ErrInternal, Err: err}
	}

	return response, nil
}

func verifySignature(headers http.Header, body, secret string) error {
	signature := headers.Get("x-satgate-signature")
	if signature == "" {
		return errors.New("missing signature")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	expected := hex.EncodeToString(mac.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(expected), []byte(signature)) != 1 {
		return errors.New("signature mismatch")
	}

	return nil
}
