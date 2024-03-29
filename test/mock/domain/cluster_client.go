// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/finleap-connect/monoskope/pkg/api/domain (interfaces: ClusterClient,Cluster_GetAllClient,ClusterAccessClient,ClusterAccess_GetClusterAccessV2Client)

// Package domain is a generated GoMock package.
package domain

import (
	context "context"
	reflect "reflect"

	domain "github.com/finleap-connect/monoskope/pkg/api/domain"
	projections "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// MockClusterClient is a mock of ClusterClient interface.
type MockClusterClient struct {
	ctrl     *gomock.Controller
	recorder *MockClusterClientMockRecorder
}

// MockClusterClientMockRecorder is the mock recorder for MockClusterClient.
type MockClusterClientMockRecorder struct {
	mock *MockClusterClient
}

// NewMockClusterClient creates a new mock instance.
func NewMockClusterClient(ctrl *gomock.Controller) *MockClusterClient {
	mock := &MockClusterClient{ctrl: ctrl}
	mock.recorder = &MockClusterClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterClient) EXPECT() *MockClusterClientMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockClusterClient) GetAll(arg0 context.Context, arg1 *domain.GetAllRequest, arg2 ...grpc.CallOption) (domain.Cluster_GetAllClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAll", varargs...)
	ret0, _ := ret[0].(domain.Cluster_GetAllClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockClusterClientMockRecorder) GetAll(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockClusterClient)(nil).GetAll), varargs...)
}

// GetById mocks base method.
func (m *MockClusterClient) GetById(arg0 context.Context, arg1 *wrapperspb.StringValue, arg2 ...grpc.CallOption) (*projections.Cluster, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetById", varargs...)
	ret0, _ := ret[0].(*projections.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockClusterClientMockRecorder) GetById(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockClusterClient)(nil).GetById), varargs...)
}

// GetByName mocks base method.
func (m *MockClusterClient) GetByName(arg0 context.Context, arg1 *wrapperspb.StringValue, arg2 ...grpc.CallOption) (*projections.Cluster, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetByName", varargs...)
	ret0, _ := ret[0].(*projections.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockClusterClientMockRecorder) GetByName(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockClusterClient)(nil).GetByName), varargs...)
}

// MockCluster_GetAllClient is a mock of Cluster_GetAllClient interface.
type MockCluster_GetAllClient struct {
	ctrl     *gomock.Controller
	recorder *MockCluster_GetAllClientMockRecorder
}

// MockCluster_GetAllClientMockRecorder is the mock recorder for MockCluster_GetAllClient.
type MockCluster_GetAllClientMockRecorder struct {
	mock *MockCluster_GetAllClient
}

// NewMockCluster_GetAllClient creates a new mock instance.
func NewMockCluster_GetAllClient(ctrl *gomock.Controller) *MockCluster_GetAllClient {
	mock := &MockCluster_GetAllClient{ctrl: ctrl}
	mock.recorder = &MockCluster_GetAllClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCluster_GetAllClient) EXPECT() *MockCluster_GetAllClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockCluster_GetAllClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockCluster_GetAllClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockCluster_GetAllClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockCluster_GetAllClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockCluster_GetAllClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockCluster_GetAllClient)(nil).Context))
}

// Header mocks base method.
func (m *MockCluster_GetAllClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockCluster_GetAllClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockCluster_GetAllClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockCluster_GetAllClient) Recv() (*projections.Cluster, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*projections.Cluster)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockCluster_GetAllClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockCluster_GetAllClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockCluster_GetAllClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockCluster_GetAllClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockCluster_GetAllClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method.
func (m *MockCluster_GetAllClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockCluster_GetAllClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockCluster_GetAllClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockCluster_GetAllClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockCluster_GetAllClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockCluster_GetAllClient)(nil).Trailer))
}

// MockClusterAccessClient is a mock of ClusterAccessClient interface.
type MockClusterAccessClient struct {
	ctrl     *gomock.Controller
	recorder *MockClusterAccessClientMockRecorder
}

// MockClusterAccessClientMockRecorder is the mock recorder for MockClusterAccessClient.
type MockClusterAccessClientMockRecorder struct {
	mock *MockClusterAccessClient
}

// NewMockClusterAccessClient creates a new mock instance.
func NewMockClusterAccessClient(ctrl *gomock.Controller) *MockClusterAccessClient {
	mock := &MockClusterAccessClient{ctrl: ctrl}
	mock.recorder = &MockClusterAccessClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterAccessClient) EXPECT() *MockClusterAccessClientMockRecorder {
	return m.recorder
}

// GetClusterAccess mocks base method.
func (m *MockClusterAccessClient) GetClusterAccess(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (domain.ClusterAccess_GetClusterAccessClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetClusterAccess", varargs...)
	ret0, _ := ret[0].(domain.ClusterAccess_GetClusterAccessClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusterAccess indicates an expected call of GetClusterAccess.
func (mr *MockClusterAccessClientMockRecorder) GetClusterAccess(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterAccess", reflect.TypeOf((*MockClusterAccessClient)(nil).GetClusterAccess), varargs...)
}

// GetClusterAccessV2 mocks base method.
func (m *MockClusterAccessClient) GetClusterAccessV2(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (domain.ClusterAccess_GetClusterAccessV2Client, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetClusterAccessV2", varargs...)
	ret0, _ := ret[0].(domain.ClusterAccess_GetClusterAccessV2Client)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusterAccessV2 indicates an expected call of GetClusterAccessV2.
func (mr *MockClusterAccessClientMockRecorder) GetClusterAccessV2(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusterAccessV2", reflect.TypeOf((*MockClusterAccessClient)(nil).GetClusterAccessV2), varargs...)
}

// GetTenantClusterMappingByTenantAndClusterId mocks base method.
func (m *MockClusterAccessClient) GetTenantClusterMappingByTenantAndClusterId(arg0 context.Context, arg1 *domain.GetClusterMappingRequest, arg2 ...grpc.CallOption) (*projections.TenantClusterBinding, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTenantClusterMappingByTenantAndClusterId", varargs...)
	ret0, _ := ret[0].(*projections.TenantClusterBinding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTenantClusterMappingByTenantAndClusterId indicates an expected call of GetTenantClusterMappingByTenantAndClusterId.
func (mr *MockClusterAccessClientMockRecorder) GetTenantClusterMappingByTenantAndClusterId(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTenantClusterMappingByTenantAndClusterId", reflect.TypeOf((*MockClusterAccessClient)(nil).GetTenantClusterMappingByTenantAndClusterId), varargs...)
}

// GetTenantClusterMappingsByClusterId mocks base method.
func (m *MockClusterAccessClient) GetTenantClusterMappingsByClusterId(arg0 context.Context, arg1 *wrapperspb.StringValue, arg2 ...grpc.CallOption) (domain.ClusterAccess_GetTenantClusterMappingsByClusterIdClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTenantClusterMappingsByClusterId", varargs...)
	ret0, _ := ret[0].(domain.ClusterAccess_GetTenantClusterMappingsByClusterIdClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTenantClusterMappingsByClusterId indicates an expected call of GetTenantClusterMappingsByClusterId.
func (mr *MockClusterAccessClientMockRecorder) GetTenantClusterMappingsByClusterId(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTenantClusterMappingsByClusterId", reflect.TypeOf((*MockClusterAccessClient)(nil).GetTenantClusterMappingsByClusterId), varargs...)
}

// GetTenantClusterMappingsByTenantId mocks base method.
func (m *MockClusterAccessClient) GetTenantClusterMappingsByTenantId(arg0 context.Context, arg1 *wrapperspb.StringValue, arg2 ...grpc.CallOption) (domain.ClusterAccess_GetTenantClusterMappingsByTenantIdClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTenantClusterMappingsByTenantId", varargs...)
	ret0, _ := ret[0].(domain.ClusterAccess_GetTenantClusterMappingsByTenantIdClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTenantClusterMappingsByTenantId indicates an expected call of GetTenantClusterMappingsByTenantId.
func (mr *MockClusterAccessClientMockRecorder) GetTenantClusterMappingsByTenantId(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTenantClusterMappingsByTenantId", reflect.TypeOf((*MockClusterAccessClient)(nil).GetTenantClusterMappingsByTenantId), varargs...)
}

// MockClusterAccess_GetClusterAccessV2Client is a mock of ClusterAccess_GetClusterAccessV2Client interface.
type MockClusterAccess_GetClusterAccessV2Client struct {
	ctrl     *gomock.Controller
	recorder *MockClusterAccess_GetClusterAccessV2ClientMockRecorder
}

// MockClusterAccess_GetClusterAccessV2ClientMockRecorder is the mock recorder for MockClusterAccess_GetClusterAccessV2Client.
type MockClusterAccess_GetClusterAccessV2ClientMockRecorder struct {
	mock *MockClusterAccess_GetClusterAccessV2Client
}

// NewMockClusterAccess_GetClusterAccessV2Client creates a new mock instance.
func NewMockClusterAccess_GetClusterAccessV2Client(ctrl *gomock.Controller) *MockClusterAccess_GetClusterAccessV2Client {
	mock := &MockClusterAccess_GetClusterAccessV2Client{ctrl: ctrl}
	mock.recorder = &MockClusterAccess_GetClusterAccessV2ClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterAccess_GetClusterAccessV2Client) EXPECT() *MockClusterAccess_GetClusterAccessV2ClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockClusterAccess_GetClusterAccessV2Client) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockClusterAccess_GetClusterAccessV2ClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockClusterAccess_GetClusterAccessV2Client)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockClusterAccess_GetClusterAccessV2Client) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockClusterAccess_GetClusterAccessV2ClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockClusterAccess_GetClusterAccessV2Client)(nil).Context))
}

// Header mocks base method.
func (m *MockClusterAccess_GetClusterAccessV2Client) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockClusterAccess_GetClusterAccessV2ClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockClusterAccess_GetClusterAccessV2Client)(nil).Header))
}

// Recv mocks base method.
func (m *MockClusterAccess_GetClusterAccessV2Client) Recv() (*projections.ClusterAccessV2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*projections.ClusterAccessV2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockClusterAccess_GetClusterAccessV2ClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockClusterAccess_GetClusterAccessV2Client)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockClusterAccess_GetClusterAccessV2Client) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockClusterAccess_GetClusterAccessV2ClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockClusterAccess_GetClusterAccessV2Client)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method.
func (m *MockClusterAccess_GetClusterAccessV2Client) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockClusterAccess_GetClusterAccessV2ClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockClusterAccess_GetClusterAccessV2Client)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockClusterAccess_GetClusterAccessV2Client) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockClusterAccess_GetClusterAccessV2ClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockClusterAccess_GetClusterAccessV2Client)(nil).Trailer))
}
