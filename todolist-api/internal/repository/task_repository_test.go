package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/ShoAnn/go/todolist-api/internal/db/sqlc"
	"github.com/ShoAnn/go/todolist-api/internal/domain"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// repository_test.go
var testPool *pgxpool.Pool

func setupSchema(connStr string) {
	// Points to your actual migration files on disk
	m, err := migrate.New("file://../db/migrations", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Runs all "up" migrations to build the whole DB
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}
func TestMain(m *testing.M) {
	ctx := context.Background()
	dbName := "todolist"
	dbUsername := "user"
	dbPassword := "password"

	// 1. Spin up the Postgres container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUsername),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("Failed to start container: %s", err)
	}

	// 2. Get the connection string from the container
	connStr, _ := postgresContainer.ConnectionString(ctx, "sslmode=disable")

	// 3. Initialize pgxpool
	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	// 4. Setup Schema (Migrations)
	setupSchema(connStr)

	// 5. Run Tests and Cleanup
	code := m.Run()
	testPool.Close()
	postgresContainer.Terminate(ctx)
	os.Exit(code)
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()
	tx, _ := testPool.Begin(ctx)
	defer tx.Rollback(ctx)

	tests := []struct {
		name   string
		seed   func(q *db.Queries)
		assert func(t *testing.T, tasks []*domain.Task)
	}{
		{
			name: "empty table",
			seed: func(q *db.Queries) {
				// no-operation
			},
			assert: func(t *testing.T, tasks []*domain.Task) {
				require.Empty(t, tasks)
			},
		},
		{
			name: "seeded tasks table",
			seed: func(q *db.Queries) {
				_, err := q.CreateTask(ctx, db.CreateTaskParams{
					Title: "task one",
				})
				require.NoError(t, err)

				_, err = q.CreateTask(ctx, db.CreateTaskParams{
					Title:     "task two",
					Completed: true,
				})
				require.NoError(t, err)
			},
			assert: func(t *testing.T, tasks []*domain.Task) {
				require.Len(t, tasks, 2)
				require.Equal(t, "task one", tasks[0].Title)
				require.Equal(t, false, tasks[0].Completed)
				require.Equal(t, "task two", tasks[1].Title)
				require.Equal(t, true, tasks[1].Completed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testPool.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			q := db.New(tx)
			repo := NewTaskRepository(tx)

			tt.seed(q)

			tasks, err := repo.GetAll(ctx)
			require.NoError(t, err)

			tt.assert(t, tasks)
		})
	}
}

func TestGetById(t *testing.T) {
	ctx := context.Background()
	tx, _ := testPool.Begin(ctx)
	defer tx.Rollback(ctx)

	tests := []struct {
		name      string
		seed      func(q *db.Queries) int32
		assert    func(t *testing.T, task *domain.Task)
		expectErr bool
	}{
		{
			name: "Happy path",
			seed: func(q *db.Queries) int32 {
				task, err := q.CreateTask(ctx, db.CreateTaskParams{
					Title: "task one",
				})
				require.NoError(t, err)
				return task.ID
			},
			assert: func(t *testing.T, task *domain.Task) {
				require.Equal(t, "task one", task.Title)
				require.Equal(t, false, task.Completed)
			},
		},
		{
			name: "id not found",
			seed: func(q *db.Queries) int32 {
				return -1
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testPool.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			q := db.New(tx)
			repo := NewTaskRepository(tx)

			taskId := tt.seed(q)

			task, err := repo.GetById(ctx, int(taskId))

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err, pgx.ErrNoRows)
			tt.assert(t, task)
		})
	}
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	tx, _ := testPool.Begin(ctx)
	defer tx.Rollback(ctx)

	tests := []struct {
		name      string
		params    *domain.CreateTaskParams
		assert    func(t *testing.T, q *db.Queries, task *domain.Task)
		expectErr bool
	}{
		{
			name: "Happy path",
			params: &domain.CreateTaskParams{
				Title: "Shopping",
			},
			assert: func(t *testing.T, q *db.Queries, task *domain.Task) {
				require.NotZero(t, task.ID)
				require.Equal(t, "Shopping", task.Title)
				require.False(t, task.Completed)
				require.Zero(t, task.Version)

				createdTask, err := q.GetTask(ctx, int32(task.ID))
				require.NoError(t, err)
				require.Equal(t, int32(task.ID), createdTask.ID)
				require.Equal(t, task.Completed, createdTask.Completed)
				require.EqualValues(t, task.Version, createdTask.Version)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testPool.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			repo := NewTaskRepository(tx)

			task, err := repo.Create(ctx, tt.params)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			q := db.New(tx)
			tt.assert(t, q, task)
		})
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()
	tx, _ := testPool.Begin(ctx)
	defer tx.Rollback(ctx)

	newTitle := "Shopping Grocery"
	newCompleted := true

	tests := []struct {
		name      string
		params    *domain.UpdateTaskParams
		seed      func(q *db.Queries) int32
		assert    func(t *testing.T, q *db.Queries, task *domain.Task)
		expectErr bool
	}{
		{
			name: "Happy path",
			params: &domain.UpdateTaskParams{
				Title:     &newTitle,
				Completed: &newCompleted,
			},
			seed: func(q *db.Queries) int32 {
				task, err := q.CreateTask(ctx, db.CreateTaskParams{
					Title: "Shopping",
				})
				require.NoError(t, err)
				return task.ID
			},
			assert: func(t *testing.T, q *db.Queries, task *domain.Task) {
				updatedTask, err := q.GetTask(ctx, int32(task.ID))
				require.NoError(t, err)
				require.Equal(t, int32(task.ID), updatedTask.ID)
				require.Equal(t, task.Completed, updatedTask.Completed)
				require.EqualValues(t, task.Version+1, updatedTask.Version)
				require.NotEqualValues(t, task.UpdatedAt, updatedTask.UpdatedAt)
			},
			expectErr: false,
		},
		{
			name: "id not found",
			params: &domain.UpdateTaskParams{
				Title:     &newTitle,
				Completed: &newCompleted,
			},
			seed: func(q *db.Queries) int32 {
				return -1
			},
			assert: func(t *testing.T, q *db.Queries, task *domain.Task) {
				_, err := q.GetTask(ctx, int32(task.ID))
				require.Error(t, err, pgx.ErrNoRows)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testPool.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			repo := NewTaskRepository(tx)

			q := db.New(tx)
			taskId := tt.seed(q)
			task, err := repo.Update(ctx, int(taskId), tt.params)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			tt.assert(t, q, task)
		})
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	tx, _ := testPool.Begin(ctx)
	defer tx.Rollback(ctx)

	tests := []struct {
		name      string
		seed      func(q *db.Queries) int32
		assert    func(t *testing.T, q *db.Queries, id int32)
		expectErr bool
	}{
		{
			name: "Happy path",
			seed: func(q *db.Queries) int32 {
				task, err := q.CreateTask(ctx, db.CreateTaskParams{
					Title: "Shopping",
				})
				require.NoError(t, err)
				return task.ID
			},
			assert: func(t *testing.T, q *db.Queries, id int32) {
				_, err := q.GetTask(ctx, id)
				require.Error(t, err, pgx.ErrNoRows)
			},
			expectErr: false,
		},
		{
			name: "id not found",
			seed: func(q *db.Queries) int32 {
				return -1
			},
			assert: func(t *testing.T, q *db.Queries, id int32) {
				_, err := q.GetTask(ctx, id)
				require.Error(t, err, pgx.ErrNoRows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testPool.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			repo := NewTaskRepository(tx)

			q := db.New(tx)
			taskId := tt.seed(q)
			err = repo.Delete(ctx, int(taskId))

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			tt.assert(t, q, taskId)
		})
	}
}

func TestMarkCompleted(t *testing.T) {
	ctx := context.Background()
	tx, _ := testPool.Begin(ctx)
	defer tx.Rollback(ctx)

	tests := []struct {
		name      string
		seed      func(q *db.Queries) int32
		assert    func(t *testing.T, q *db.Queries, id int32)
		expectErr bool
	}{
		{
			name: "Happy path",
			seed: func(q *db.Queries) int32 {
				task, err := q.CreateTask(ctx, db.CreateTaskParams{
					Title: "Shopping",
				})
				require.NoError(t, err)
				return task.ID
			},
			assert: func(t *testing.T, q *db.Queries, id int32) {
				updatedTask, err := q.GetTask(ctx, id)
				require.NoError(t, err)
				require.True(t, updatedTask.Completed)
			},
			expectErr: false,
		},
		{
			name: "id not found",
			seed: func(q *db.Queries) int32 {
				return -1
			},
			assert: func(t *testing.T, q *db.Queries, id int32) {
				_, err := q.GetTask(ctx, id)
				require.Error(t, err, pgx.ErrNoRows)
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testPool.Begin(ctx)
			require.NoError(t, err)
			defer tx.Rollback(ctx)

			repo := NewTaskRepository(tx)

			q := db.New(tx)
			taskId := tt.seed(q)
			err = repo.MarkCompleted(ctx, int(taskId))

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			tt.assert(t, q, taskId)
		})
	}
}
