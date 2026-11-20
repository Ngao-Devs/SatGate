package messages

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) ListPaidByOwner(ctx context.Context, ownerID string) ([]Message, error) {
	return s.repository.ListPaidByOwner(ctx, ownerID)
}
