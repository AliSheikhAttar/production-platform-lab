package task

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMemoryRepository_Create(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	created, err := repo.Create(
		ctx,
		"Learn Go",
	)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("expected ID 1, got %d", created.ID)
	}

	if created.Title != "Learn Go" {
		t.Errorf("expected title %q, got %q", "Learn Go", created.Title)
	}

	if created.Status != StatusPending {
		t.Errorf("expected status %q, got %q", StatusPending, created.Status)
	}

	if created.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if created.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	if !created.CreatedAt.Equal(created.UpdatedAt) {
		t.Error("expected CreatedAt and UpdatedAt to initially be equal")
	}
}

func TestMemoryRepository_Create_AssignsIncrementingIDs(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	task1, err := repo.Create(ctx, "Task 1")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	task2, err := repo.Create(ctx, "Task 2")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	task3, err := repo.Create(ctx, "Task 3")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if task1.ID != 1 {
		t.Errorf("expected first ID 1, got %d", task1.ID)
	}

	if task2.ID != 2 {
		t.Errorf("expected second ID 2, got %d", task2.ID)
	}

	if task3.ID != 3 {
		t.Errorf("expected third ID 3, got %d", task3.ID)
	}
}

func TestMemoryRepository_Get(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	created, err := repo.Create(ctx, "Learn PostgreSQL")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	got, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, got.ID)
	}

	if got.Title != created.Title {
		t.Errorf("expected title %q, got %q", created.Title, got.Title)
	}

	if got.Status != created.Status {
		t.Errorf("expected status %q, got %q", created.Status, got.Status)
	}
}

func TestMemoryRepository_Get_NotFound(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	_, err := repo.Get(ctx, 999)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryRepository_List(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	if _, err := repo.Create(ctx, "Task 1"); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if _, err := repo.Create(ctx, "Task 2"); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if _, err := repo.Create(ctx, "Task 3"); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestMemoryRepository_List_Empty(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestMemoryRepository_UpdateStatus(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	created, err := repo.Create(ctx, "Process report")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	// Small delay so UpdatedAt should definitely change.
	time.Sleep(time.Millisecond)

	updated, err := repo.UpdateStatus(
		ctx,
		created.ID,
		StatusProcessing,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.Status != StatusProcessing {
		t.Errorf(
			"expected status %q, got %q",
			StatusProcessing,
			updated.Status,
		)
	}

	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Error("expected UpdatedAt to be updated")
	}

	// Verify the change was actually persisted.
	got, err := repo.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Status != StatusProcessing {
		t.Errorf(
			"expected persisted status %q, got %q",
			StatusProcessing,
			got.Status,
		)
	}
}

func TestMemoryRepository_UpdateStatus_NotFound(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	_, err := repo.UpdateStatus(
		ctx,
		999,
		StatusCompleted,
	)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryRepository_ConcurrentCreate(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	errs := make(chan error, goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()

			_, err := repo.Create(ctx, "task")
			if err != nil {
				errs <- err
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("Create returned error: %v", err)
	}

	tasks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(tasks) != goroutines {
		t.Fatalf(
			"expected %d tasks, got %d",
			goroutines,
			len(tasks),
		)
	}

	// Verify every ID is unique.
	ids := make(map[int64]bool)

	for _, task := range tasks {
		if ids[task.ID] {
			t.Fatalf("duplicate task ID found: %d", task.ID)
		}

		ids[task.ID] = true
	}
}
