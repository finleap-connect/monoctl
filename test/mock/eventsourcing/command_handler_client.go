// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing (interfaces: CommandHandlerClient)

// Package eventsourcing is a generated GoMock package.
package eventsourcing

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	eventsourcing "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	commands "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/commands"
	grpc "google.golang.org/grpc"
)

// MockCommandHandlerClient is a mock of CommandHandlerClient interface.
type MockCommandHandlerClient struct {
	ctrl     *gomock.Controller
	recorder *MockCommandHandlerClientMockRecorder
}

// MockCommandHandlerClientMockRecorder is the mock recorder for MockCommandHandlerClient.
type MockCommandHandlerClientMockRecorder struct {
	mock *MockCommandHandlerClient
}

// NewMockCommandHandlerClient creates a new mock instance.
func NewMockCommandHandlerClient(ctrl *gomock.Controller) *MockCommandHandlerClient {
	mock := &MockCommandHandlerClient{ctrl: ctrl}
	mock.recorder = &MockCommandHandlerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandHandlerClient) EXPECT() *MockCommandHandlerClientMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockCommandHandlerClient) Execute(arg0 context.Context, arg1 *commands.Command, arg2 ...grpc.CallOption) (*eventsourcing.CommandReply, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(*eventsourcing.CommandReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockCommandHandlerClientMockRecorder) Execute(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCommandHandlerClient)(nil).Execute), varargs...)
}
