// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/service/service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/kotche/url-shortening-service/internal/app/model"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockStorage) Add(userID string, url *model.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", userID, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockStorageMockRecorder) Add(userID, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockStorage)(nil).Add), userID, url)
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// GetByID mocks base method.
func (m *MockStorage) GetByID(id string) (*model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(*model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockStorageMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockStorage)(nil).GetByID), id)
}

// GetUserURLs mocks base method.
func (m *MockStorage) GetUserURLs(userID string) ([]*model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", userID)
	ret0, _ := ret[0].([]*model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockStorageMockRecorder) GetUserURLs(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockStorage)(nil).GetUserURLs), userID)
}

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockDatabase) Add(userID string, url *model.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", userID, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockDatabaseMockRecorder) Add(userID, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockDatabase)(nil).Add), userID, url)
}

// Close mocks base method.
func (m *MockDatabase) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockDatabaseMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDatabase)(nil).Close))
}

// DeleteBatch mocks base method.
func (m *MockDatabase) DeleteBatch(ctx context.Context, toDelete []model.DeleteUserURLs) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBatch", ctx, toDelete)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBatch indicates an expected call of DeleteBatch.
func (mr *MockDatabaseMockRecorder) DeleteBatch(ctx, toDelete interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBatch", reflect.TypeOf((*MockDatabase)(nil).DeleteBatch), ctx, toDelete)
}

// GetByID mocks base method.
func (m *MockDatabase) GetByID(id string) (*model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(*model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockDatabaseMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockDatabase)(nil).GetByID), id)
}

// GetNumberOfURLs mocks base method.
func (m *MockDatabase) GetNumberOfURLs(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNumberOfURLs", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNumberOfURLs indicates an expected call of GetNumberOfURLs.
func (mr *MockDatabaseMockRecorder) GetNumberOfURLs(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNumberOfURLs", reflect.TypeOf((*MockDatabase)(nil).GetNumberOfURLs), ctx)
}

// GetNumberOfUsers mocks base method.
func (m *MockDatabase) GetNumberOfUsers(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNumberOfUsers", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNumberOfUsers indicates an expected call of GetNumberOfUsers.
func (mr *MockDatabaseMockRecorder) GetNumberOfUsers(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNumberOfUsers", reflect.TypeOf((*MockDatabase)(nil).GetNumberOfUsers), ctx)
}

// GetUserURLs mocks base method.
func (m *MockDatabase) GetUserURLs(userID string) ([]*model.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", userID)
	ret0, _ := ret[0].([]*model.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockDatabaseMockRecorder) GetUserURLs(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockDatabase)(nil).GetUserURLs), userID)
}

// Ping mocks base method.
func (m *MockDatabase) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockDatabaseMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockDatabase)(nil).Ping))
}

// WriteBatch mocks base method.
func (m *MockDatabase) WriteBatch(ctx context.Context, userID string, urls map[string]*model.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteBatch", ctx, userID, urls)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteBatch indicates an expected call of WriteBatch.
func (mr *MockDatabaseMockRecorder) WriteBatch(ctx, userID, urls interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteBatch", reflect.TypeOf((*MockDatabase)(nil).WriteBatch), ctx, userID, urls)
}

// MockIGenerator is a mock of IGenerator interface.
type MockIGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockIGeneratorMockRecorder
}

// MockIGeneratorMockRecorder is the mock recorder for MockIGenerator.
type MockIGeneratorMockRecorder struct {
	mock *MockIGenerator
}

// NewMockIGenerator creates a new mock instance.
func NewMockIGenerator(ctrl *gomock.Controller) *MockIGenerator {
	mock := &MockIGenerator{ctrl: ctrl}
	mock.recorder = &MockIGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIGenerator) EXPECT() *MockIGeneratorMockRecorder {
	return m.recorder
}

// MakeShortURL mocks base method.
func (m *MockIGenerator) MakeShortURL() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeShortURL")
	ret0, _ := ret[0].(string)
	return ret0
}

// MakeShortURL indicates an expected call of MakeShortURL.
func (mr *MockIGeneratorMockRecorder) MakeShortURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeShortURL", reflect.TypeOf((*MockIGenerator)(nil).MakeShortURL))
}
