package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestReservationService(t *testing.T) {
	dbSrv := setupTestDB(t)
	svc := NewReservationService(dbSrv)
	ctx := context.Background()

	t.Run("CreateReservation - Success", func(t *testing.T) {
		input := CreateReservationInput{
			Name:    "Alice Smith",
			Contact: "alice@example.com",
			Service: "Maintenance Valet",
			Date:    "2026-09-01T10:00:00Z",
			Message: "Bring extra detailing towels",
		}

		res, err := svc.CreateReservation(ctx, input)
		if err != nil {
			t.Fatalf("CreateReservation returned error: %v", err)
		}

		if res.ID == 0 {
			t.Error("expected reservation ID to be non-zero")
		}
		if res.Name != input.Name {
			t.Errorf("expected Name %s, got %s", input.Name, res.Name)
		}
		if res.Contact != input.Contact {
			t.Errorf("expected Contact %s, got %s", input.Contact, res.Contact)
		}
		if res.Service != input.Service {
			t.Errorf("expected Service %s, got %s", input.Service, res.Service)
		}
		if !res.Date.Equal(time.Date(2026, 9, 1, 10, 0, 0, 0, time.UTC)) {
			t.Errorf("expected Date 2026-09-01T10:00:00Z, got %v", res.Date)
		}
		if res.Message != input.Message {
			t.Errorf("expected Message %s, got %s", input.Message, res.Message)
		}
	})

	t.Run("CreateReservation - Invalid Date Format", func(t *testing.T) {
		input := CreateReservationInput{
			Name:    "Bob Smith",
			Contact: "bob@example.com",
			Service: "Deep Clean Valet",
			Date:    "invalid-date",
		}

		_, err := svc.CreateReservation(ctx, input)
		if err == nil {
			t.Error("expected error for invalid date format, got nil")
		}
	})

	t.Run("CreateReservation - Duplicate and Slot Taken Conflicts", func(t *testing.T) {
		dbSrv3 := setupTestDB(t)
		svc3 := NewReservationService(dbSrv3)

		input1 := CreateReservationInput{
			Name:    "John Doe",
			Contact: "john@example.com",
			Service: "Deep Clean Valet",
			Date:    "2026-09-05T14:00:00Z",
		}
		_, err := svc3.CreateReservation(ctx, input1)
		if err != nil {
			t.Fatalf("failed to create first reservation: %v", err)
		}

		// Try duplicate (same customer, same time)
		_, err = svc3.CreateReservation(ctx, input1)
		if !errors.Is(err, ErrDuplicateReservation) {
			t.Errorf("expected ErrDuplicateReservation, got %v", err)
		}

		// Try slot conflict (different customer, same time)
		input2 := CreateReservationInput{
			Name:    "Jane Doe",
			Contact: "jane@example.com",
			Service: "Maintenance Valet",
			Date:    "2026-09-05T14:00:00Z",
		}
		_, err = svc3.CreateReservation(ctx, input2)
		if !errors.Is(err, ErrSlotTaken) {
			t.Errorf("expected ErrSlotTaken, got %v", err)
		}
	})

	t.Run("ListReservationsPublic & Admin", func(t *testing.T) {
		// Clean / reset DB
		dbSrv2 := setupTestDB(t)
		svc2 := NewReservationService(dbSrv2)

		inputs := []CreateReservationInput{
			{Name: "User 1", Contact: "111", Service: "S1", Date: "2026-09-01T09:00:00Z"},
			{Name: "User 2", Contact: "222", Service: "S2", Date: "2026-09-02T10:00:00Z"},
		}

		for _, inp := range inputs {
			_, err := svc2.CreateReservation(ctx, inp)
			if err != nil {
				t.Fatalf("failed to create reservation: %v", err)
			}
		}

		// Test Public List
		publicList, err := svc2.ListReservationsPublic(ctx)
		if err != nil {
			t.Fatalf("ListReservationsPublic failed: %v", err)
		}
		if len(publicList) != 2 {
			t.Errorf("expected 2 public reservations, got %d", len(publicList))
		}
		// Confirm only Date is exposed
		if !publicList[0].Date.Equal(time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC)) {
			t.Errorf("expected public date 2026-09-01T09:00:00Z, got %v", publicList[0].Date)
		}

		// Test Admin List
		adminList, err := svc2.ListReservationsAdmin(ctx)
		if err != nil {
			t.Fatalf("ListReservationsAdmin failed: %v", err)
		}
		if len(adminList) != 2 {
			t.Errorf("expected 2 admin reservations, got %d", len(adminList))
		}
		if adminList[0].Name != "User 1" || adminList[0].Contact != "111" {
			t.Errorf("expected admin list to contain full details, got Name: %s, Contact: %s", adminList[0].Name, adminList[0].Contact)
		}
	})
}
