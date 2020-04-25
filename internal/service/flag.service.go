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

		hash, _ := req.Hash()
		eval = &flaggio.Evaluation{
			FlagID:      flg.ID,
			FlagVersion: flg.Version,
			FlagKey:     flg.Key,
			RequestHash: hash,
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
	prevEvals, err := s.evalsRepo.FindAllByUserID(ctx, req.UserID, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	// fetch all flags
	flgs, err := s.flagsRepo.FindAll(ctx, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	// fetch segments
	sgmts, err := segmentsAsIdentifiers(s.segmentsRepo.FindAll(ctx, nil, nil))
	if err != nil {
		return nil, err
	}

	// check for missing flag evaluations
	hash, _ := req.Hash()
	validEvals := validFlagEvals(hash, flgs.Flags, prevEvals.Evaluations)
	evals := make([]*flaggio.Evaluation, len(flgs.Flags))

	// evaluate flags
	evalSpan, _ := opentracing.StartSpanFromContext(ctx, "flaggio.Evaluate")
	for idx, flg := range flgs.Flags {
		if evltn, ok := validEvals[flg.ID]; ok {
			evals[idx] = evltn
			continue
		}
		flg.Populate(sgmts)

		evltn := &flaggio.Evaluation{
			FlagID:      flg.ID,
			FlagVersion: flg.Version,
			FlagKey:     flg.Key,
			RequestHash: hash,
		}
		res, err := flaggio.Evaluate(req.UserContext, flg)
		if err != nil {
			evltn.Error = err.Error()
		} else {
			evltn.Value = res.Answer
		}

		evals[idx] = evltn
	}
	evalSpan.Finish()

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

func validFlagEvals(reqHash string, flgs []*flaggio.Flag, evals flaggio.EvaluationList) map[string]*flaggio.Evaluation {
	validEvals := map[string]*flaggio.Evaluation{}
	for _, eval := range evals {
		validEvals[eval.FlagID] = eval
	}
	for _, flg := range flgs {
		eval, ok := validEvals[flg.ID]
		// eval is missing if no evaluation found, the evaluation was
		// for a previous flag version or the user context changed
		if !ok || flg.Version != eval.FlagVersion || reqHash != eval.RequestHash {
			delete(validEvals, flg.ID)
		}
	}
	return validEvals
}
