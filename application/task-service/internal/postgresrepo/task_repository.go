package postgresrepo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/alisheikhattar/production-platform-lab/application/task-service/internal/task"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(
	pool *pgxpool.Pool,
) *TaskRepository {
	return &TaskRepository{
		pool: pool,
	}
}

func (r *TaskRepository) Create(
	ctx context.Context,
	title string,
) (task.Task, error) {

	const query = `
		INSERT INTO tasks (title)
		VALUES ($1)
		RETURNING
			id,
			title,
			status,
			created_at,
			updated_at
	`

	var t task.Task

	err := r.pool.QueryRow(
		ctx,
		query,
		title,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

func (r *TaskRepository) Get(
	ctx context.Context,
	id int64,
) (task.Task, error) {

	const query = `
		SELECT
			id,
			title,
			status,
			created_at,
			updated_at
		FROM tasks
		WHERE id = $1
	`

	var t task.Task

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return task.Task{}, task.ErrNotFound
	}

	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

func (r *TaskRepository) List(
	ctx context.Context,
) ([]task.Task, error) {

	const query = `
		SELECT
			id,
			title,
			status,
			created_at,
			updated_at
		FROM tasks
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]task.Task, 0)

	for rows.Next() {
		var t task.Task

		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateStatus(
	ctx context.Context,
	id int64,
	status task.Status,
) (task.Task, error) {

	const query = `
		UPDATE tasks
		SET
			status = $2,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id,
			title,
			status,
			created_at,
			updated_at
	`

	var t task.Task

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
		status,
	).Scan(
		&t.ID,
		&t.Title,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return task.Task{}, task.ErrNotFound
	}

	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}
