package task

import (
	"sync"
	"testing"
	"time"
)

func TestMemoryRepository_Create(t *testing.T) {
	repo := NewMemoryRepository()

	task := repo.Create("Learn Go")

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}

	if task.Title != "Learn Go" {
		t.Errorf("expected title %q, got %q", "Learn Go", task.Title)
	}

	if task.Status != StatusPending {
		t.Errorf("expected status %q, got %q", StatusPending, task.Status)
	}

	if task.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if task.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	if !task.CreatedAt.Equal(task.UpdatedAt) {
		t.Error("expected CreatedAt and UpdatedAt to initially be equal")
	}
}

func TestMemoryRepository_Create_AssignsIncrementingIDs(t *testing.T) {
	repo := NewMemoryRepository()

	task1 := repo.Create("Task 1")
	task2 := repo.Create("Task 2")
	task3 := repo.Create("Task 3")

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

	created := repo.Create("Learn PostgreSQL")

	got, err := repo.Get(created.ID)
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

	_, err := repo.Get(999)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryRepository_List(t *testing.T) {
	repo := NewMemoryRepository()

	repo.Create("Task 1")
	repo.Create("Task 2")
	repo.Create("Task 3")

	tasks := repo.List()

	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}
}

func TestMemoryRepository_List_Empty(t *testing.T) {
	repo := NewMemoryRepository()

	tasks := repo.List()

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}

func TestMemoryRepository_UpdateStatus(t *testing.T) {
	repo := NewMemoryRepository()

	created := repo.Create("Process report")

	// Small delay so UpdatedAt should definitely change.
	time.Sleep(time.Millisecond)

	updated, err := repo.UpdateStatus(
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
	got, err := repo.Get(created.ID)
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

	_, err := repo.UpdateStatus(999, StatusCompleted)

	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryRepository_ConcurrentCreate(t *testing.T) {
	repo := NewMemoryRepository()

	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()

			repo.Create("task")
		}()
	}

	wg.Wait()

	tasks := repo.List()

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
