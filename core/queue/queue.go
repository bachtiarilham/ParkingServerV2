package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

var ErrJobNotReady = errors.New("job not completed yet")

type Job struct {
	ID            string          `json:"job_id"`
	Topic         string          `json:"topic"`
	Status        string          `json:"status"`
	Attempts      int             `json:"attempts"`
	MaxAttempts   int             `json:"max_attempts"`
	Payload       json.RawMessage `json:"payload,omitempty"`
	ResultPayload json.RawMessage `json:"result,omitempty"`
	LastError     string          `json:"last_error,omitempty"`
	CreatedAt     string          `json:"created_at,omitempty"`
	UpdatedAt     string          `json:"updated_at,omitempty"`
	CompletedAt   string          `json:"completed_at,omitempty"`
}

type Queue struct {
	db           *sql.DB
	leaseTimeout time.Duration
	retryDelay   time.Duration
}

func Open(path string, leaseTimeout, retryDelay time.Duration, maxOpenConns int) (*Queue, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	if maxOpenConns <= 0 {
		maxOpenConns = 1
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxOpenConns)

	queue := &Queue{
		db:           db,
		leaseTimeout: leaseTimeout,
		retryDelay:   retryDelay,
	}
	if err := queue.ensureSchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return queue, nil
}

func (q *Queue) Close() error {
	return q.db.Close()
}

func (q *Queue) ensureSchema(ctx context.Context) error {
	statements := []string{
		`PRAGMA journal_mode = WAL;`,
		`PRAGMA busy_timeout = 5000;`,
		`CREATE TABLE IF NOT EXISTS job_queue (
			id TEXT PRIMARY KEY,
			topic TEXT NOT NULL,
			payload TEXT NOT NULL,
			status TEXT NOT NULL,
			attempts INTEGER NOT NULL DEFAULT 0,
			max_attempts INTEGER NOT NULL DEFAULT 20,
			available_at DATETIME NOT NULL,
			lease_until DATETIME NULL,
			last_error TEXT NULL,
			result_payload TEXT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			completed_at DATETIME NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_job_queue_available ON job_queue (status, available_at, created_at);`,
	}

	for _, statement := range statements {
		if _, err := q.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

func newJobID() string {
	return fmt.Sprintf("job-%s", time.Now().Format("200601021504050000000"))
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func (q *Queue) Enqueue(ctx context.Context, topic string, payload any) (Job, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return Job{}, err
	}

	now := time.Now().UTC()
	job := Job{
		ID:          newJobID(),
		Topic:       topic,
		Status:      "pending",
		Attempts:    0,
		MaxAttempts: 20,
		Payload:     body,
		CreatedAt:   formatTime(now),
		UpdatedAt:   formatTime(now),
	}

	const query = `
		INSERT INTO job_queue
		(id, topic, payload, status, attempts, max_attempts, available_at, created_at, updated_at)
		VALUES (?, ?, ?, 'pending', 0, 20, ?, ?, ?)`
	if _, err := q.db.ExecContext(ctx, query, job.ID, job.Topic, string(job.Payload), now, now, now); err != nil {
		return Job{}, err
	}

	return job, nil
}

func (q *Queue) Get(ctx context.Context, jobID string) (Job, error) {
	const query = `
		SELECT id, topic, payload, status, attempts, max_attempts, COALESCE(last_error, ''), COALESCE(result_payload, ''), created_at, updated_at, completed_at
		FROM job_queue
		WHERE id = ?
		LIMIT 1`

	var job Job
	var payloadText string
	var resultText string
	var createdAt time.Time
	var updatedAt time.Time
	var completedAt sql.NullTime
	if err := q.db.QueryRowContext(ctx, query, jobID).Scan(
		&job.ID,
		&job.Topic,
		&payloadText,
		&job.Status,
		&job.Attempts,
		&job.MaxAttempts,
		&job.LastError,
		&resultText,
		&createdAt,
		&updatedAt,
		&completedAt,
	); err != nil {
		return Job{}, err
	}

	job.Payload = json.RawMessage(payloadText)
	if resultText != "" {
		job.ResultPayload = json.RawMessage(resultText)
	}
	job.CreatedAt = formatTime(createdAt)
	job.UpdatedAt = formatTime(updatedAt)
	if completedAt.Valid {
		job.CompletedAt = formatTime(completedAt.Time)
	}
	return job, nil
}

func (q *Queue) Wait(ctx context.Context, jobID string, interval time.Duration) (Job, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		job, err := q.Get(ctx, jobID)
		if err != nil {
			return Job{}, err
		}
		switch job.Status {
		case "completed", "failed":
			return job, nil
		}

		select {
		case <-ctx.Done():
			return Job{}, ctx.Err()
		case <-ticker.C:
		}
	}
}

func (q *Queue) ClaimNext(ctx context.Context) (Job, error) {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return Job{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now().UTC()
	const query = `
		SELECT id, topic, payload, status, attempts, max_attempts, COALESCE(last_error, ''), COALESCE(result_payload, ''), created_at, updated_at, completed_at
		FROM job_queue
		WHERE status IN ('pending', 'retry')
		  AND available_at <= ?
		  AND (lease_until IS NULL OR lease_until <= ?)
		ORDER BY created_at ASC
		LIMIT 1`

	var job Job
	var payloadText string
	var resultText string
	var createdAt time.Time
	var updatedAt time.Time
	var completedAt sql.NullTime
	err = tx.QueryRowContext(ctx, query, now, now).Scan(
		&job.ID,
		&job.Topic,
		&payloadText,
		&job.Status,
		&job.Attempts,
		&job.MaxAttempts,
		&job.LastError,
		&resultText,
		&createdAt,
		&updatedAt,
		&completedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Job{}, sql.ErrNoRows
		}
		return Job{}, err
	}

	job.Payload = json.RawMessage(payloadText)
	if resultText != "" {
		job.ResultPayload = json.RawMessage(resultText)
	}
	job.CreatedAt = formatTime(createdAt)
	job.UpdatedAt = formatTime(updatedAt)
	if completedAt.Valid {
		job.CompletedAt = formatTime(completedAt.Time)
	}

	leaseUntil := now.Add(q.leaseTimeout)
	if _, err := tx.ExecContext(ctx,
		`UPDATE job_queue SET status = 'processing', attempts = attempts + 1, lease_until = ?, updated_at = ? WHERE id = ?`,
		leaseUntil, now, job.ID,
	); err != nil {
		return Job{}, err
	}

	if err := tx.Commit(); err != nil {
		return Job{}, err
	}

	job.Status = "processing"
	job.Attempts++
	job.UpdatedAt = formatTime(now)
	return job, nil
}

func (q *Queue) MarkCompleted(ctx context.Context, jobID string, result any) error {
	body, err := json.Marshal(result)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	_, err = q.db.ExecContext(ctx,
		`UPDATE job_queue SET status = 'completed', result_payload = ?, lease_until = NULL, last_error = NULL, completed_at = ?, updated_at = ? WHERE id = ?`,
		string(body), now, now, jobID,
	)
	return err
}

func (q *Queue) MarkRetry(ctx context.Context, job Job, lastErr error) error {
	now := time.Now().UTC()
	nextTime := now.Add(q.retryDelay)
	status := "retry"
	if job.Attempts >= job.MaxAttempts {
		status = "failed"
	}

	var completedAt any
	if status == "failed" {
		completedAt = now
	} else {
		completedAt = nil
	}

	_, err := q.db.ExecContext(ctx,
		`UPDATE job_queue SET status = ?, last_error = ?, lease_until = NULL, available_at = ?, updated_at = ?, completed_at = ? WHERE id = ?`,
		status, lastErr.Error(), nextTime, now, completedAt, job.ID,
	)
	return err
}
