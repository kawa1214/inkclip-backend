// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/bookmark-manager/bookmark-manager/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/bookmark-manager/bookmark-manager/db/sqlc"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateNote mocks base method.
func (m *MockStore) CreateNote(arg0 context.Context, arg1 db.CreateNoteParams) (db.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNote", arg0, arg1)
	ret0, _ := ret[0].(db.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNote indicates an expected call of CreateNote.
func (mr *MockStoreMockRecorder) CreateNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNote", reflect.TypeOf((*MockStore)(nil).CreateNote), arg0, arg1)
}

// CreateNoteWeb mocks base method.
func (m *MockStore) CreateNoteWeb(arg0 context.Context, arg1 db.CreateNoteWebParams) (db.NoteWeb, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNoteWeb", arg0, arg1)
	ret0, _ := ret[0].(db.NoteWeb)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNoteWeb indicates an expected call of CreateNoteWeb.
func (mr *MockStoreMockRecorder) CreateNoteWeb(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNoteWeb", reflect.TypeOf((*MockStore)(nil).CreateNoteWeb), arg0, arg1)
}

// CreateSession mocks base method.
func (m *MockStore) CreateSession(arg0 context.Context, arg1 db.CreateSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockStoreMockRecorder) CreateSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockStore)(nil).CreateSession), arg0, arg1)
}

// CreateTemporaryUser mocks base method.
func (m *MockStore) CreateTemporaryUser(arg0 context.Context, arg1 db.CreateTemporaryUserParams) (db.TemporaryUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTemporaryUser", arg0, arg1)
	ret0, _ := ret[0].(db.TemporaryUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTemporaryUser indicates an expected call of CreateTemporaryUser.
func (mr *MockStoreMockRecorder) CreateTemporaryUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTemporaryUser", reflect.TypeOf((*MockStore)(nil).CreateTemporaryUser), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// CreateWeb mocks base method.
func (m *MockStore) CreateWeb(arg0 context.Context, arg1 db.CreateWebParams) (db.Web, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWeb", arg0, arg1)
	ret0, _ := ret[0].(db.Web)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateWeb indicates an expected call of CreateWeb.
func (mr *MockStoreMockRecorder) CreateWeb(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWeb", reflect.TypeOf((*MockStore)(nil).CreateWeb), arg0, arg1)
}

// DeleteNote mocks base method.
func (m *MockStore) DeleteNote(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNote", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNote indicates an expected call of DeleteNote.
func (mr *MockStoreMockRecorder) DeleteNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNote", reflect.TypeOf((*MockStore)(nil).DeleteNote), arg0, arg1)
}

// DeleteNoteWeb mocks base method.
func (m *MockStore) DeleteNoteWeb(arg0 context.Context, arg1 db.DeleteNoteWebParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNoteWeb", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNoteWeb indicates an expected call of DeleteNoteWeb.
func (mr *MockStoreMockRecorder) DeleteNoteWeb(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNoteWeb", reflect.TypeOf((*MockStore)(nil).DeleteNoteWeb), arg0, arg1)
}

// DeleteNoteWebsByNoteId mocks base method.
func (m *MockStore) DeleteNoteWebsByNoteId(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNoteWebsByNoteId", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNoteWebsByNoteId indicates an expected call of DeleteNoteWebsByNoteId.
func (mr *MockStoreMockRecorder) DeleteNoteWebsByNoteId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNoteWebsByNoteId", reflect.TypeOf((*MockStore)(nil).DeleteNoteWebsByNoteId), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), arg0, arg1)
}

// DeleteWeb mocks base method.
func (m *MockStore) DeleteWeb(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWeb", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteWeb indicates an expected call of DeleteWeb.
func (mr *MockStoreMockRecorder) DeleteWeb(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWeb", reflect.TypeOf((*MockStore)(nil).DeleteWeb), arg0, arg1)
}

// GetNote mocks base method.
func (m *MockStore) GetNote(arg0 context.Context, arg1 uuid.UUID) (db.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNote", arg0, arg1)
	ret0, _ := ret[0].(db.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNote indicates an expected call of GetNote.
func (mr *MockStoreMockRecorder) GetNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNote", reflect.TypeOf((*MockStore)(nil).GetNote), arg0, arg1)
}

// GetNoteWeb mocks base method.
func (m *MockStore) GetNoteWeb(arg0 context.Context, arg1 db.GetNoteWebParams) (db.NoteWeb, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNoteWeb", arg0, arg1)
	ret0, _ := ret[0].(db.NoteWeb)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNoteWeb indicates an expected call of GetNoteWeb.
func (mr *MockStoreMockRecorder) GetNoteWeb(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNoteWeb", reflect.TypeOf((*MockStore)(nil).GetNoteWeb), arg0, arg1)
}

// GetSession mocks base method.
func (m *MockStore) GetSession(arg0 context.Context, arg1 uuid.UUID) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0, arg1)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockStoreMockRecorder) GetSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockStore)(nil).GetSession), arg0, arg1)
}

// GetTemporaryUserByEmailAndToken mocks base method.
func (m *MockStore) GetTemporaryUserByEmailAndToken(arg0 context.Context, arg1 db.GetTemporaryUserByEmailAndTokenParams) (db.TemporaryUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemporaryUserByEmailAndToken", arg0, arg1)
	ret0, _ := ret[0].(db.TemporaryUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemporaryUserByEmailAndToken indicates an expected call of GetTemporaryUserByEmailAndToken.
func (mr *MockStoreMockRecorder) GetTemporaryUserByEmailAndToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemporaryUserByEmailAndToken", reflect.TypeOf((*MockStore)(nil).GetTemporaryUserByEmailAndToken), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 uuid.UUID) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// GetUserByEmail mocks base method.
func (m *MockStore) GetUserByEmail(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockStoreMockRecorder) GetUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockStore)(nil).GetUserByEmail), arg0, arg1)
}

// GetWeb mocks base method.
func (m *MockStore) GetWeb(arg0 context.Context, arg1 uuid.UUID) (db.Web, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWeb", arg0, arg1)
	ret0, _ := ret[0].(db.Web)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWeb indicates an expected call of GetWeb.
func (mr *MockStoreMockRecorder) GetWeb(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWeb", reflect.TypeOf((*MockStore)(nil).GetWeb), arg0, arg1)
}

// ListNoteWebsByNoteId mocks base method.
func (m *MockStore) ListNoteWebsByNoteId(arg0 context.Context, arg1 uuid.UUID) ([]db.NoteWeb, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNoteWebsByNoteId", arg0, arg1)
	ret0, _ := ret[0].([]db.NoteWeb)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNoteWebsByNoteId indicates an expected call of ListNoteWebsByNoteId.
func (mr *MockStoreMockRecorder) ListNoteWebsByNoteId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNoteWebsByNoteId", reflect.TypeOf((*MockStore)(nil).ListNoteWebsByNoteId), arg0, arg1)
}

// ListNotesByUserId mocks base method.
func (m *MockStore) ListNotesByUserId(arg0 context.Context, arg1 db.ListNotesByUserIdParams) ([]db.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNotesByUserId", arg0, arg1)
	ret0, _ := ret[0].([]db.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNotesByUserId indicates an expected call of ListNotesByUserId.
func (mr *MockStoreMockRecorder) ListNotesByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNotesByUserId", reflect.TypeOf((*MockStore)(nil).ListNotesByUserId), arg0, arg1)
}

// ListUsers mocks base method.
func (m *MockStore) ListUsers(arg0 context.Context, arg1 db.ListUsersParams) ([]db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", arg0, arg1)
	ret0, _ := ret[0].([]db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockStoreMockRecorder) ListUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockStore)(nil).ListUsers), arg0, arg1)
}

// ListWebByNoteId mocks base method.
func (m *MockStore) ListWebByNoteId(arg0 context.Context, arg1 uuid.UUID) ([]db.Web, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWebByNoteId", arg0, arg1)
	ret0, _ := ret[0].([]db.Web)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWebByNoteId indicates an expected call of ListWebByNoteId.
func (mr *MockStoreMockRecorder) ListWebByNoteId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWebByNoteId", reflect.TypeOf((*MockStore)(nil).ListWebByNoteId), arg0, arg1)
}

// ListWebByNoteIds mocks base method.
func (m *MockStore) ListWebByNoteIds(arg0 context.Context, arg1 []uuid.UUID) ([]db.ListWebByNoteIdsRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWebByNoteIds", arg0, arg1)
	ret0, _ := ret[0].([]db.ListWebByNoteIdsRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWebByNoteIds indicates an expected call of ListWebByNoteIds.
func (mr *MockStoreMockRecorder) ListWebByNoteIds(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWebByNoteIds", reflect.TypeOf((*MockStore)(nil).ListWebByNoteIds), arg0, arg1)
}

// ListWebsByUserId mocks base method.
func (m *MockStore) ListWebsByUserId(arg0 context.Context, arg1 db.ListWebsByUserIdParams) ([]db.Web, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListWebsByUserId", arg0, arg1)
	ret0, _ := ret[0].([]db.Web)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListWebsByUserId indicates an expected call of ListWebsByUserId.
func (mr *MockStoreMockRecorder) ListWebsByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListWebsByUserId", reflect.TypeOf((*MockStore)(nil).ListWebsByUserId), arg0, arg1)
}

// TxCreateNote mocks base method.
func (m *MockStore) TxCreateNote(arg0 context.Context, arg1 db.TxCreateNoteParams) (db.TxCreateNoteResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxCreateNote", arg0, arg1)
	ret0, _ := ret[0].(db.TxCreateNoteResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TxCreateNote indicates an expected call of TxCreateNote.
func (mr *MockStoreMockRecorder) TxCreateNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxCreateNote", reflect.TypeOf((*MockStore)(nil).TxCreateNote), arg0, arg1)
}

// TxDeleteNote mocks base method.
func (m *MockStore) TxDeleteNote(arg0 context.Context, arg1 db.TxDeleteNoteParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxDeleteNote", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// TxDeleteNote indicates an expected call of TxDeleteNote.
func (mr *MockStoreMockRecorder) TxDeleteNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxDeleteNote", reflect.TypeOf((*MockStore)(nil).TxDeleteNote), arg0, arg1)
}

// TxUpdateNote mocks base method.
func (m *MockStore) TxUpdateNote(arg0 context.Context, arg1 db.TxUpdateNoteParams) (db.TxUpdateNoteResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TxUpdateNote", arg0, arg1)
	ret0, _ := ret[0].(db.TxUpdateNoteResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TxUpdateNote indicates an expected call of TxUpdateNote.
func (mr *MockStoreMockRecorder) TxUpdateNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxUpdateNote", reflect.TypeOf((*MockStore)(nil).TxUpdateNote), arg0, arg1)
}

// UpdateNote mocks base method.
func (m *MockStore) UpdateNote(arg0 context.Context, arg1 db.UpdateNoteParams) (db.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNote", arg0, arg1)
	ret0, _ := ret[0].(db.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNote indicates an expected call of UpdateNote.
func (mr *MockStoreMockRecorder) UpdateNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNote", reflect.TypeOf((*MockStore)(nil).UpdateNote), arg0, arg1)
}

// UpdateUser mocks base method.
func (m *MockStore) UpdateUser(arg0 context.Context, arg1 db.UpdateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStoreMockRecorder) UpdateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStore)(nil).UpdateUser), arg0, arg1)
}
