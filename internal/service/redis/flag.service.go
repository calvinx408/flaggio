package redis

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/victorkt/flaggio/internal/repository"
	"github.com/victorkt/flaggio/internal/service"
)

var _ service.Flag = (*flagService)(nil)

// flagService implements service.Flag interface using redis.
type flagService struct {
	evalsRepo repository.Evaluation
	usersRepo repository.User
	svc       service.Flag
}

// Evaluate returns the result of an evaluation of a single flag.
func (s *flagService) Evaluate(ctx context.Context, flagKey string, req *service.EvaluationRequest) (*service.EvaluationResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisFlagService.Evaluate")
	defer span.Finish()

	shouldCache := shouldCacheEvaluation(req)

	if shouldCache {
		reqHash, err := req.Hash()
		if err != nil {
			return nil, err
		}
		eval, err := s.evalsRepo.FindByReqHashAndFlagKey(ctx, reqHash, flagKey)
		if err != nil {
			return nil, err
		}
		if eval != nil {
			return &service.EvaluationResponse{
				Evaluation: eval,
			}, nil
		}
	}

	// cache miss or debug request, call underlying service
	res, err := s.svc.Evaluate(ctx, flagKey, req)
	if err != nil {
		return nil, err
	}

	if shouldCache {
		// create or update the user
		if err := s.usersRepo.Replace(ctx, req.UserID, req.UserContext); err != nil {
			return nil, err
		}
		// replace the evaluation for the user and flag key
		if err := s.evalsRepo.ReplaceOne(ctx, req.UserID, res.Evaluation); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// EvaluateAll returns the results of the evaluation of all flags.
func (s *flagService) EvaluateAll(ctx context.Context, req *service.EvaluationRequest) (*service.EvaluationsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisFlagService.EvaluateAll")
	defer span.Finish()

	shouldCache := shouldCacheEvaluation(req)

	if shouldCache {
		reqHash, err := req.Hash()
		if err != nil {
			return nil, err
		}
		evals, err := s.evalsRepo.FindAllByReqHash(ctx, reqHash)
		if err != nil {
			return nil, err
		}
		if evals != nil {
			return &service.EvaluationsResponse{
				Evaluations: evals,
			}, nil
		}
	}

	// cache miss or debug request, call underlying service
	res, err := s.svc.EvaluateAll(ctx, req)
	if err != nil {
		return nil, err
	}

	if shouldCache {
		hash, err := req.Hash()
		if err != nil {
			return nil, err
		}
		// create or update the user
		if err := s.usersRepo.Replace(ctx, req.UserID, req.UserContext); err != nil {
			return nil, err
		}
		// replace the evaluations for the user
		if err := s.evalsRepo.ReplaceAll(ctx, req.UserID, hash, res.Evaluations); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func NewFlagService(
	evalsRepo repository.Evaluation,
	usersRepo repository.User,
	svc service.Flag,
) service.Flag {
	return &flagService{
		evalsRepo: evalsRepo,
		usersRepo: usersRepo,
		svc:       svc,
	}
}
