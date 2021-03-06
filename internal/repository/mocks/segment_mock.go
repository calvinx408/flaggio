// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/victorkt/flaggio/internal/repository (interfaces: Segment)

// Package repository_mock is a generated GoMock package.
package repository_mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	flaggio "github.com/victorkt/flaggio/internal/flaggio"
	reflect "reflect"
)

// MockSegment is a mock of Segment interface
type MockSegment struct {
	ctrl     *gomock.Controller
	recorder *MockSegmentMockRecorder
}

// MockSegmentMockRecorder is the mock recorder for MockSegment
type MockSegmentMockRecorder struct {
	mock *MockSegment
}

// NewMockSegment creates a new mock instance
func NewMockSegment(ctrl *gomock.Controller) *MockSegment {
	mock := &MockSegment{ctrl: ctrl}
	mock.recorder = &MockSegmentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSegment) EXPECT() *MockSegmentMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockSegment) Create(arg0 context.Context, arg1 flaggio.NewSegment) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockSegmentMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSegment)(nil).Create), arg0, arg1)
}

// Delete mocks base method
func (m *MockSegment) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockSegmentMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSegment)(nil).Delete), arg0, arg1)
}

// FindAll mocks base method
func (m *MockSegment) FindAll(arg0 context.Context, arg1, arg2 *int64) ([]*flaggio.Segment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*flaggio.Segment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll
func (mr *MockSegmentMockRecorder) FindAll(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockSegment)(nil).FindAll), arg0, arg1, arg2)
}

// FindByID mocks base method
func (m *MockSegment) FindByID(arg0 context.Context, arg1 string) (*flaggio.Segment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*flaggio.Segment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockSegmentMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockSegment)(nil).FindByID), arg0, arg1)
}

// Update mocks base method
func (m *MockSegment) Update(arg0 context.Context, arg1 string, arg2 flaggio.UpdateSegment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockSegmentMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSegment)(nil).Update), arg0, arg1, arg2)
}
