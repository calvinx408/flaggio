package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
	"github.com/vmihailenco/msgpack/v4"
)

var _ repository.Evaluation = (*EvaluationRepository)(nil)

// EvaluationRepository implements repository.Evaluation interface using redis.
type EvaluationRepository struct {
	redis *redis.Client
	store repository.Evaluation
	ttl   time.Duration
}

// FindAllByUserID returns all previous flag evaluations for a given user ID.
func (r *EvaluationRepository) FindAllByUserID(ctx context.Context, userID string, search *string, offset, limit *int64) (*flaggio.EvaluationResults, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.FindAllByUserID")
	defer span.Finish()

	// this is called by the admin api, don't cache it
	return r.store.FindAllByUserID(ctx, userID, search, offset, limit)
}

// FindByReqHashAndFlagKey returns a previous flag evaluation for a given request hash and flag key.
func (r *EvaluationRepository) FindByReqHashAndFlagKey(ctx context.Context, reqHash, flagKey string) (*flaggio.Evaluation, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.FindByReqHashAndFlagKey")
	defer span.Finish()

	cacheKey := flaggio.EvalCacheKey(reqHash, flagKey)

	// fetch evaluation results from cache
	cached, err := r.redis.WithContext(ctx).Get(cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		// an unexpected error occurred, return it
		return nil, err
	}
	if cached != "" {
		// cache hit, unmarshall and return result
		var e flaggio.Evaluation
		if err := msgpack.Unmarshal([]byte(cached), &e); err == nil {
			// return if no errors, otherwise defer to the store
			return &e, nil
		}
	}

	// cache miss or disabled, fetch from store
	eval, err := r.store.FindByReqHashAndFlagKey(ctx, reqHash, flagKey)
	if err != nil {
		return nil, err
	}

	// marshal and save result
	b, err := msgpack.Marshal(eval)
	if err != nil {
		return nil, err
	}
	if err := r.redis.Set(cacheKey, b, r.ttl).Err(); err != nil {
		return nil, err
	}

	return eval, nil
}

// FindAllByReqHash returns all previous flag evaluations for a given request hash.
func (r *EvaluationRepository) FindAllByReqHash(ctx context.Context, reqHash string) (flaggio.EvaluationList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.FindAllByReqHash")
	defer span.Finish()

	cacheKey := flaggio.EvalCacheKey(reqHash)

	// fetch evaluation results from cache
	cached, err := r.redis.WithContext(ctx).Get(cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		// an unexpected error occurred, return it
		return nil, err
	}
	if cached != "" {
		// cache hit, unmarshall and return result
		var el flaggio.EvaluationList
		if err := msgpack.Unmarshal([]byte(cached), &el); err == nil {
			// return if no errors, otherwise defer to the store
			return el, nil
		}
	}

	// cache miss or disabled, fetch from store
	evals, err := r.store.FindAllByReqHash(ctx, reqHash)
	if err != nil {
		return nil, err
	}

	// marshal and save result
	b, err := msgpack.Marshal(evals)
	if err != nil {
		return nil, err
	}
	if err := r.redis.Set(cacheKey, b, r.ttl).Err(); err != nil {
		return nil, err
	}

	return evals, nil
}

// ReplaceOne creates or replaces one evaluation for a combination of user ID, request hash and flag key.
func (r *EvaluationRepository) ReplaceOne(ctx context.Context, userID string, eval *flaggio.Evaluation) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.ReplaceOne")
	defer span.Finish()

	// call underlying store
	err := r.store.ReplaceOne(ctx, userID, eval)
	if err != nil {
		return err
	}

	// marshall and save result
	b, err := msgpack.Marshal(eval)
	if err != nil {
		return err
	}
	cacheKey := flaggio.EvalCacheKey(eval.RequestHash, eval.FlagKey)
	if err := r.redis.Set(cacheKey, b, r.ttl).Err(); err != nil {
		return err
	}

	// invalidate all relevant keys
	return nil
}

// ReplaceAll creates or replaces evaluations for a combination of user and request hash.
func (r *EvaluationRepository) ReplaceAll(ctx context.Context, userID, reqHash string, evals flaggio.EvaluationList) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.ReplaceAll")
	defer span.Finish()

	// call underlying store
	err := r.store.ReplaceAll(ctx, userID, reqHash, evals)
	if err != nil {
		return err
	}

	// marshall and save result
	b, err := msgpack.Marshal(evals)
	if err != nil {
		return err
	}
	cacheKey := flaggio.EvalCacheKey(reqHash)
	if err := r.redis.Set(cacheKey, b, r.ttl).Err(); err != nil {
		return err
	}

	// invalidate all relevant keys
	return nil
}

// DeleteAllByUserID deletes evaluations for a user.
func (r *EvaluationRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.DeleteAllByUserID")
	defer span.Finish()

	if err := r.store.DeleteAllByUserID(ctx, userID); err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx)
}

// DeleteByID deletes an evaluation by its ID.
func (r *EvaluationRepository) DeleteByID(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisEvaluationRepository.DeleteByID")
	defer span.Finish()

	// delete the flag
	if err := r.store.DeleteByID(ctx, id); err != nil {
		return err
	}

	// invalidate all relevant keys
	return r.invalidateRelevantCacheKeys(ctx)
}

func (r *EvaluationRepository) invalidateRelevantCacheKeys(ctx context.Context) error {
	redisCtx := r.redis.WithContext(ctx)

	// invalidate all relevant keys
	keysToInvalidate, err := redisCtx.Keys(flaggio.EvalCacheKey("*")).Result()
	if err != nil {
		return err
	}

	return redisCtx.Del(keysToInvalidate...).Err()
}

// NewEvaluationRepository returns a new evaluation repository that uses redis
// as underlying storage.
func NewEvaluationRepository(redisClient *redis.Client, store repository.Evaluation) repository.Evaluation {
	return &EvaluationRepository{
		redis: redisClient,
		store: store,
		ttl:   1 * time.Hour,
	}
}
