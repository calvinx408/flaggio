package repository

//go:generate mockgen -destination=./mocks/evaluation_mock.go -package=repository_mock github.com/victorkt/flaggio/internal/repository Evaluation

import (
	"context"

	"github.com/victorkt/flaggio/internal/flaggio"
)

// Flag represents a set of operations available to list and manage evaluations.
type Evaluation interface {
	// FindAllByUserID returns all previous flag evaluations for a given user ID.
	FindAllByUserID(ctx context.Context, userID string) (flaggio.EvaluationList, error)
	// FindByUserIDAndFlagID returns a previous flag evaluation for a given user ID and flag ID.
	FindByUserIDAndFlagID(ctx context.Context, userID, flagID string) (*flaggio.Evaluation, error)
	// Replace creates or replaces evaluations for a user.
	Replace(ctx context.Context, userID string, evals flaggio.EvaluationList) error
	// DeleteAllByUserID deletes evaluations for a user.
	DeleteAllByUserID(ctx context.Context, userID string) error
}
