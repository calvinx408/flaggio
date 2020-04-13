package mongo

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
	"github.com/victorkt/flaggio/internal/service"
)

var _ service.Flag = (*flagService)(nil)

// flagService implements service.Flag interface using mongo.
type flagService struct {
	evalsRepo repository.Evaluation
	usersRepo repository.User
	svc       service.Flag
}

func (f *flagService) Evaluate(ctx context.Context, flagKey string, req *service.EvaluationRequest) (*service.EvaluationResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoFlagService.Evaluate")
	defer span.Finish()

	// evaluate all flags
	res, err := f.svc.Evaluate(ctx, flagKey, req)
	if err != nil {
		return nil, err
	}

	if !req.IsDebug() {
		// create or update the user
		if err := f.usersRepo.Replace(ctx, req.UserID, req.UserContext); err != nil {
			return nil, err
		}

		// replace the evaluations for the user
		if err := f.evalsRepo.Replace(ctx, req.UserID, flaggio.EvaluationList{res.Evaluation}); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (f *flagService) EvaluateAll(ctx context.Context, req *service.EvaluationRequest) (*service.EvaluationsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MongoFlagService.EvaluateAll")
	defer span.Finish()

	// evaluate all flags
	res, err := f.svc.EvaluateAll(ctx, req)
	if err != nil {
		return nil, err
	}

	if !req.IsDebug() {
		// create or update the user
		if err := f.usersRepo.Replace(ctx, req.UserID, req.UserContext); err != nil {
			return nil, err
		}

		// replace the evaluations for the user
		if err := f.evalsRepo.Replace(ctx, req.UserID, res.Evaluations); err != nil {
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
