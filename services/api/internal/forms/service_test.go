package forms

import (
	"context"
	"errors"
	"testing"
)

func newFormsTestService(t *testing.T) *Service {
	t.Helper()
	return NewService(newFormsTestRepository(t))
}

func TestServiceCreateTrimsInput(t *testing.T) {
	service := newFormsTestService(t)

	form, err := service.Create(context.Background(), "owner-a", CreateFormRequest{
		Name:       "  Contact form  ",
		Domain:     "  example.com  ",
		AmountSats: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	if form.Name != "Contact form" {
		t.Errorf("name = %q; want Contact form", form.Name)
	}
	if form.Domain != "example.com" {
		t.Errorf("domain = %q; want example.com", form.Domain)
	}
	if form.ID == "" {
		t.Error("expected generated form ID")
	}
	if form.CreatedAt.IsZero() {
		t.Error("expected creation timestamp")
	}
}

func TestServiceCreateRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		request CreateFormRequest
		wantErr error
	}{
		{
			name: "missing name",
			request: CreateFormRequest{
				Domain: "example.com", AmountSats: 100,
			},
			wantErr: ErrInvalidName,
		},
		{
			name: "missing domain",
			request: CreateFormRequest{
				Name: "Contact form", AmountSats: 100,
			},
			wantErr: ErrInvalidDomain,
		},
		{
			name: "amount below minimum",
			request: CreateFormRequest{
				Name: "Contact form", Domain: "example.com", AmountSats: 0,
			},
			wantErr: ErrInvalidAmount,
		},
		{
			name: "amount above maximum",
			request: CreateFormRequest{
				Name: "Contact form", Domain: "example.com", AmountSats: 10001,
			},
			wantErr: ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newFormsTestService(t)

			_, err := service.Create(context.Background(), "owner-a", tt.request)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("got error %v; want %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	service := newFormsTestService(t)
	ctx := context.Background()

	created, err := service.Create(ctx, "owner-a", CreateFormRequest{
		Name: "Original", Domain: "example.com", AmountSats: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := service.Update(ctx, "owner-a", created.ID, UpdateFormRequest{
		Name: "  Updated  ", Domain: "  updated.example.com  ", AmountSats: 250,
	})
	if err != nil {
		t.Fatal(err)
	}

	if updated.Name != "Updated" || updated.Domain != "updated.example.com" || updated.AmountSats != 250 {
		t.Fatalf("unexpected updated form: %+v", updated)
	}
}

func TestServiceUpdateCannotChangeAnotherOwnersForm(t *testing.T) {
	service := newFormsTestService(t)
	ctx := context.Background()

	created, err := service.Create(ctx, "owner-a", CreateFormRequest{
		Name: "Original", Domain: "example.com", AmountSats: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.Update(ctx, "owner-b", created.ID, UpdateFormRequest{
		Name: "Hijacked", Domain: "attacker.example.com", AmountSats: 200,
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got error %v; want ErrNotFound", err)
	}
}

func TestServiceDelete(t *testing.T) {
	service := newFormsTestService(t)
	ctx := context.Background()

	created, err := service.Create(ctx, "owner-a", CreateFormRequest{
		Name: "To delete", Domain: "example.com", AmountSats: 100,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := service.Delete(ctx, "owner-a", created.ID); err != nil {
		t.Fatal(err)
	}

	_, err = service.GetPublic(ctx, created.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("got error %v; want ErrNotFound", err)
	}
}
