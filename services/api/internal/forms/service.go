package forms

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidName      = errors.New("form name is required")
	ErrInvalidDomain    = errors.New("form domain is required")
	ErrInvalidAmount    = errors.New("amount must be between 1 and 10000 sats")
	ErrInvalidMinAmount = errors.New("minimum amount must be between 1 and 10000 sats")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, ownerID string) ([]Form, error) {
	return s.repository.ListByOwner(ctx, ownerID)
}

func (s *Service) Create(ctx context.Context, ownerID string, request CreateFormRequest) (*Form, error) {
	form := Form{
		ID:            uuid.NewString(),
		Name:          strings.TrimSpace(request.Name),
		Domain:        strings.TrimSpace(request.Domain),
		AmountSats:    request.AmountSats,
		MinAmountSats: request.MinAmountSats,
		CreatedAt:     time.Now().UTC(),
	}

	if form.MinAmountSats == 0 {
		form.MinAmountSats = 5
	}

	if err := validateForm(form.Name, form.Domain, form.AmountSats, form.MinAmountSats); err != nil {
		return nil, err
	}

	if err := s.repository.Create(ctx, ownerID, form); err != nil {
		return nil, err
	}

	return &form, nil
}

func (s *Service) Update(ctx context.Context, ownerID, id string, request UpdateFormRequest) (*Form, error) {
	form := Form{
		ID:            id,
		Name:          strings.TrimSpace(request.Name),
		Domain:        strings.TrimSpace(request.Domain),
		AmountSats:    request.AmountSats,
		MinAmountSats: request.MinAmountSats,
	}

	if form.MinAmountSats == 0 {
		form.MinAmountSats = 5
	}

	if err := validateForm(form.Name, form.Domain, form.AmountSats, form.MinAmountSats); err != nil {
		return nil, err
	}

	if err := s.repository.Update(ctx, ownerID, form); err != nil {
		return nil, err
	}

	return s.repository.GetByIDAndOwner(ctx, id, ownerID)
}

func (s *Service) Delete(ctx context.Context, ownerID, id string) error {
	return s.repository.Delete(ctx, id, ownerID)
}

func (s *Service) GetPublic(ctx context.Context, id string) (*Form, error) {
	return s.repository.GetByID(ctx, id)
}

func validateForm(name, domain string, amountSats, minAmountSats int64) error {
	if name == "" {
		return ErrInvalidName
	}
	if domain == "" {
		return ErrInvalidDomain
	}
	if amountSats < 1 || amountSats > 10000 {
		return ErrInvalidAmount
	}
	if minAmountSats < 1 || minAmountSats > 10000 {
		return ErrInvalidMinAmount
	}
	return nil
}
