package invoices

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetStatus(ctx context.Context, invoiceID string) (*InvoiceStatusResponse, error) {
	return s.repository.GetStatus(ctx, invoiceID)
}

func (s *Service) MarkPaid(ctx context.Context, invoiceID string, preimage *string) error {
	return s.repository.MarkPaid(ctx, invoiceID, preimage)
}
