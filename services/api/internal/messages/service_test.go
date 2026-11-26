package messages

import (
	"context"
	"testing"
)

func TestServiceListPaidByOwner(t *testing.T) {
	repository := newTestRepository(t)
	service := NewService(repository)

	messages, err := service.ListPaidByOwner(context.Background(), "owner-a")
	if err != nil {
		t.Fatal(err)
	}

	if messages == nil {
		t.Fatal("expected an empty slice, got nil")
	}

	if len(messages) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(messages))
	}
}
