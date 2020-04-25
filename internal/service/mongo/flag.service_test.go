package mongo_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/victorkt/flaggio/internal/flaggio"
	repository_mock "github.com/victorkt/flaggio/internal/repository/mocks"
	"github.com/victorkt/flaggio/internal/service"
	service_mock "github.com/victorkt/flaggio/internal/service/mocks"
	mongo_svc "github.com/victorkt/flaggio/internal/service/mongo"
)

var (
	evalsList = flaggio.EvaluationList{
		{FlagKey: "f1", Value: "abc"},
		{FlagKey: "f2", Value: int64(1)},
	}
	ctxInterface = reflect.TypeOf((*context.Context)(nil)).Elem()
)

func TestFlagService_Evaluate(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T, repo *service_mock.MockFlag, evalsRepo *repository_mock.MockEvaluation, usersRepo *repository_mock.MockUser)
	}{
		{
			name: "replaces users and evaluations when not a debug request",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag, evalsRepo *repository_mock.MockEvaluation, usersRepo *repository_mock.MockUser) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalRes := &service.EvaluationResponse{Evaluation: evalsList[0], UserContext: nil}
				flagMongoSvc := mongo_svc.NewFlagService(evalsRepo, usersRepo, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
				}
				flagSvc.EXPECT().Evaluate(gomock.AssignableToTypeOf(ctxInterface), "f1", evalReq).
					Times(1).Return(evalRes, nil)
				usersRepo.EXPECT().Replace(gomock.AssignableToTypeOf(ctxInterface), "123", flaggio.UserContext{"test": "abc"}).
					Times(1).Return(nil)
				evalsRepo.EXPECT().Replace(gomock.AssignableToTypeOf(ctxInterface), "123", flaggio.EvaluationList{evalRes.Evaluation}).
					Times(1).Return(nil)

				res, err := flagMongoSvc.Evaluate(ctx, "f1", evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalRes, res)
			},
		},
		{
			name: "don't replace users and evaluations when debug request",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag, evalsRepo *repository_mock.MockEvaluation, usersRepo *repository_mock.MockUser) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalRes := &service.EvaluationResponse{Evaluation: evalsList[0], UserContext: nil}
				flagMongoSvc := mongo_svc.NewFlagService(evalsRepo, usersRepo, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
					Debug:       boolPtr(true),
				}
				flagSvc.EXPECT().Evaluate(gomock.AssignableToTypeOf(ctxInterface), "f1", evalReq).
					Times(1).Return(evalRes, nil)
				usersRepo.EXPECT().Replace(gomock.AssignableToTypeOf(ctxInterface), "123", flaggio.UserContext{"test": "abc"}).
					Times(0).Return(nil)
				evalsRepo.EXPECT().Replace(gomock.AssignableToTypeOf(ctxInterface), "123", flaggio.EvaluationList{evalRes.Evaluation}).
					Times(0).Return(nil)

				res, err := flagMongoSvc.Evaluate(ctx, "f1", evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalRes, res)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			evalsRepo := repository_mock.NewMockEvaluation(mockCtrl)
			usersRepo := repository_mock.NewMockUser(mockCtrl)
			flagSvc := service_mock.NewMockFlag(mockCtrl)

			tt.run(t, flagSvc, evalsRepo, usersRepo)
		})
	}
}

func TestFlagService_EvaluateAll(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T, repo *service_mock.MockFlag, evalsRepo *repository_mock.MockEvaluation, usersRepo *repository_mock.MockUser)
	}{
		{
			name: "replaces users and evaluations when not a debug request",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag, evalsRepo *repository_mock.MockEvaluation, usersRepo *repository_mock.MockUser) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalsRes := &service.EvaluationsResponse{Evaluations: evalsList, UserContext: nil}
				flagMongoSvc := mongo_svc.NewFlagService(evalsRepo, usersRepo, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
				}
				flagSvc.EXPECT().EvaluateAll(gomock.AssignableToTypeOf(ctxInterface), evalReq).
					Times(1).Return(evalsRes, nil)
				usersRepo.EXPECT().Replace(gomock.AssignableToTypeOf(ctxInterface), "123", flaggio.UserContext{"test": "abc"}).
					Times(1).Return(nil)
				evalsRepo.EXPECT().Replace(gomock.AssignableToTypeOf(ctxInterface), "123", evalsRes.Evaluations).
					Times(1).Return(nil)

				res, err := flagMongoSvc.EvaluateAll(ctx, evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalsRes, res)
			},
		},
		{
			name: "don't replace users and evaluations when debug request",
			run: func(t *testing.T, flagSvc *service_mock.MockFlag, evalsRepo *repository_mock.MockEvaluation, usersRepo *repository_mock.MockUser) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				evalsRes := &service.EvaluationsResponse{Evaluations: evalsList, UserContext: nil}
				flagMongoSvc := mongo_svc.NewFlagService(evalsRepo, usersRepo, flagSvc)
				evalReq := &service.EvaluationRequest{
					UserID:      "123",
					UserContext: flaggio.UserContext{"test": "abc"},
					Debug:       boolPtr(true),
				}
				flagSvc.EXPECT().EvaluateAll(gomock.AssignableToTypeOf(ctxInterface), evalReq).
					Times(1).Return(evalsRes, nil)
				usersRepo.EXPECT().Replace(gomock.AssignableToTypeOf(ctxInterface), "123", flaggio.UserContext{"test": "abc"}).
					Times(0).Return(nil)
				evalsRepo.EXPECT().Replace(gomock.AssignableToTypeOf(ctxInterface), "123", evalsRes.Evaluations).
					Times(0).Return(nil)

				res, err := flagMongoSvc.EvaluateAll(ctx, evalReq)
				assert.NoError(t, err)
				assert.Equal(t, evalsRes, res)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			evalsRepo := repository_mock.NewMockEvaluation(mockCtrl)
			usersRepo := repository_mock.NewMockUser(mockCtrl)
			flagSvc := service_mock.NewMockFlag(mockCtrl)

			tt.run(t, flagSvc, evalsRepo, usersRepo)
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}
