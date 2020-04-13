package service

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	apperrors "github.com/victorkt/flaggio/internal/errors"
	"github.com/victorkt/flaggio/internal/flaggio"
	"github.com/victorkt/flaggio/internal/repository"
)

var _ Flag = (*flagService)(nil)

// NewFlagService returns a new Flag service
func NewFlagService(
	flagsRepo repository.Flag,
	segmentsRepo repository.Segment,
	evalsRepo repository.Evaluation,
) Flag {
	return &flagService{
		flagsRepo:    flagsRepo,
		segmentsRepo: segmentsRepo,
		evalsRepo:    evalsRepo,
	}
}

type flagService struct {
	flagsRepo    repository.Flag
	segmentsRepo repository.Segment
	evalsRepo    repository.Evaluation
}

// Evaluate evaluates a flag by key, returning a value based on the user context
func (s *flagService) Evaluate(ctx context.Context, flagKey string, req *EvaluationRequest) (*EvaluationResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlagService.Evaluate")
	defer span.Finish()

	// fetch flag
	flg, err := s.flagsRepo.FindByKey(ctx, flagKey)
	if err != nil {
		return nil, err
	}
	// fetch previous evaluations for this flag
	eval, err := s.evalsRepo.FindByUserIDAndFlagID(ctx, req.UserID, flg.ID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}

	// if there are no previous evaluations, evaluate the flag
	if eval == nil {
		sgmts, err := segmentsAsIdentifiers(s.segmentsRepo.FindAll(ctx, nil, nil))
		if err != nil {
			return nil, err
		}

		flg.Populate(sgmts)

		evalSpan, _ := opentracing.StartSpanFromContext(ctx, "flaggio.Evaluate")
		res, err := flaggio.Evaluate(req.UserContext, flg)
		evalSpan.Finish()
		if err != nil {
			return nil, err
		}

		eval = &flaggio.Evaluation{
			FlagID:      flg.ID,
			FlagVersion: flg.Version,
			FlagKey:     flg.Key,
			Value:       res.Answer,
		}
		if req.IsDebug() {
			eval.StackTrace = res.Stack()
		}
	}

	// build the response
	evalRes := &EvaluationResponse{
		Evaluation: eval,
	}

	if req.IsDebug() {
		evalRes.UserContext = &req.UserContext
	}

	return evalRes, nil
}

// EvaluateAll evaluates all flags, returning a value or an error for each flag based on the user context
func (s *flagService) EvaluateAll(ctx context.Context, req *EvaluationRequest) (*EvaluationsResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlagService.EvaluateAll")
	defer span.Finish()

	// fetch previous evaluations
	evals, err := s.evalsRepo.FindAllByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	// fetch all flags
	flgs, err := s.flagsRepo.FindAll(ctx, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	// check for missing flag evaluations
	missingFlagEvals := missingFlagEvals(flgs.Flags, evals)
	if len(missingFlagEvals) > 0 {
		// fetch segments
		sgmts, err := segmentsAsIdentifiers(s.segmentsRepo.FindAll(ctx, nil, nil))
		if err != nil {
			return nil, err
		}

		evalSpan, _ := opentracing.StartSpanFromContext(ctx, "flaggio.Evaluate")
		for _, flg := range missingFlagEvals {
			flg.Populate(sgmts)

			evltn := &flaggio.Evaluation{
				FlagID:      flg.ID,
				FlagVersion: flg.Version,
				FlagKey:     flg.Key,
			}
			res, err := flaggio.Evaluate(req.UserContext, flg)
			if err != nil {
				evltn.Error = err.Error()
			} else {
				evltn.Value = res.Answer
			}

			evals = append(evals, evltn)
		}
		evalSpan.Finish()
	}

	// build the response
	evalRes := &EvaluationsResponse{
		Evaluations: evals,
	}

	if req.IsDebug() {
		evalRes.UserContext = &req.UserContext
	}

	return evalRes, nil
}

func segmentsAsIdentifiers(sgmts []*flaggio.Segment, err error) ([]flaggio.Identifier, error) {
	if err != nil {
		return nil, err
	}
	iders := make([]flaggio.Identifier, len(sgmts))
	for idx, sgmnt := range sgmts {
		iders[idx] = sgmnt
	}
	return iders, nil
}

func missingFlagEvals(flgs []*flaggio.Flag, evals flaggio.EvaluationList) []*flaggio.Flag {
	evalMap := map[string]struct{}{}
	for _, eval := range evals {
		evalMap[eval.FlagID] = struct{}{}
	}
	var missing []*flaggio.Flag
	for _, flg := range flgs {
		if _, ok := evalMap[flg.ID]; !ok {
			missing = append(missing, flg)
		}
	}
	return missing
}
