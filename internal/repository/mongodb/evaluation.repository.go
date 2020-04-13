package mongodb

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/victorkt/flaggio/internal/errors"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ repository.Evaluation = (*EvaluationRepository)(nil)

// EvaluationRepository implements repository.Flag interface using mongodb.
type EvaluationRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

// FindAllByUserID returns all previous flag evaluations for a given user ID.
func (r *EvaluationRepository) FindAllByUserID(ctx context.Context, userID string) (flaggio.EvaluationList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.FindAllByUserID")
	defer span.Finish()

	cursor, err := r.col.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}

	var evals flaggio.EvaluationList
	for cursor.Next(ctx) {
		var e evaluationModel
		// decode the document
		if err := cursor.Decode(&e); err != nil {
			return nil, err
		}
		evals = append(evals, e.asEvaluation())
	}

	// check if the cursor encountered any errors while iterating
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return evals, nil
}

// FindByUserIDAndFlagKey returns a previous flag evaluation for a given user ID and flag ID.
func (r *EvaluationRepository) FindByUserIDAndFlagID(ctx context.Context, userID, flagID string) (*flaggio.Evaluation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.FindByUserIDAndFlagKey")
	defer span.Finish()

	flgID, err := primitive.ObjectIDFromHex(flagID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"userId": userID, "flagId": flgID}

	var e evaluationModel
	if err := r.col.FindOne(ctx, filter).Decode(&e); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("evaluation")
		}
		return nil, err
	}
	return e.asEvaluation(), nil
}

// Replace creates or replaces evaluations for a user.
func (r *EvaluationRepository) Replace(ctx context.Context, userID string, evals flaggio.EvaluationList) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.Replace")
	defer span.Finish()

	// get list of evaluated flag ids
	evalsToDelete := make([]primitive.ObjectID, len(evals))
	for idx, eval := range evals {
		flgID, err := primitive.ObjectIDFromHex(eval.FlagID)
		if err != nil {
			return err
		}
		evalsToDelete[idx] = flgID
	}
	// delete current
	_, err := r.col.DeleteMany(ctx, bson.M{"userId": userID, "flagId": bson.M{"$in": evalsToDelete}})
	if err != nil {
		return err
	}

	// prepare list of evaluations to insert
	evalsToInsert := make([]interface{}, len(evals))
	for idx, eval := range evals {
		flgID, err := primitive.ObjectIDFromHex(eval.FlagID)
		if err != nil {
			return err
		}
		evalsToInsert[idx] = &evaluationModel{
			ID:          primitive.NewObjectID(),
			FlagID:      flgID,
			FlagKey:     eval.FlagKey,
			FlagVersion: eval.FlagVersion,
			UserID:      userID,
			Value:       eval.Value,
			CreatedAt:   time.Now(),
		}
	}
	_, err = r.col.InsertMany(ctx, evalsToInsert)
	return err
}

// Delete deletes evaluations for a user.
func (r *EvaluationRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoEvaluationRepository.DeleteAllByUserID")
	defer span.Finish()

	_, err := r.col.DeleteMany(ctx, bson.M{"userId": userID})
	return err
}

// NewEvaluationRepository returns a new evaluation repository that uses mongodb as underlying storage.
// It also creates all needed indexes, if they don't yet exist.
func NewEvaluationRepository(ctx context.Context, db *mongo.Database) (repository.Evaluation, error) {
	col := db.Collection("evaluations")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "flagId", Value: 1}},
			Options: options.Index().SetBackground(false),
		},
	})
	if err != nil {
		return nil, err
	}
	return &EvaluationRepository{
		db:  db,
		col: col,
	}, nil
}
