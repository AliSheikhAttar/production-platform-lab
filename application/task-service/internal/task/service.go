package task

import (
	"context"
	"time"
)

func (r *MemoryRepository) Create(
	_ context.Context,
	title string,
) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()

	t := Task{
		ID:        r.nextID,
		Title:     title,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	r.tasks[t.ID] = t
	r.nextID++

	return t, nil
}

func (r *MemoryRepository) Get(
	_ context.Context,
	id int64,
) (Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}

	return t, nil
}

func (r *MemoryRepository) List(
	_ context.Context,
) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]Task, 0, len(r.tasks))

	for _, t := range r.tasks {
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *MemoryRepository) UpdateStatus(
	_ context.Context,
	id int64,
	status Status,
) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}

	t.Status = status
	t.UpdatedAt = time.Now().UTC()

	r.tasks[id] = t

	return t, nil
}
