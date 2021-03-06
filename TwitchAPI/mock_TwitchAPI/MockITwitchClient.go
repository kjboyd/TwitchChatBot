// Code generated by MockGen. DO NOT EDIT.
// Source: .\TwitchAPI\TwitchClient.go

// Package mock_TwitchAPI is a generated GoMock package.
package mock_TwitchAPI

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockITwitchClient is a mock of ITwitchClient interface.
type MockITwitchClient struct {
	ctrl     *gomock.Controller
	recorder *MockITwitchClientMockRecorder
}

// MockITwitchClientMockRecorder is the mock recorder for MockITwitchClient.
type MockITwitchClientMockRecorder struct {
	mock *MockITwitchClient
}

// NewMockITwitchClient creates a new mock instance.
func NewMockITwitchClient(ctrl *gomock.Controller) *MockITwitchClient {
	mock := &MockITwitchClient{ctrl: ctrl}
	mock.recorder = &MockITwitchClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITwitchClient) EXPECT() *MockITwitchClientMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockITwitchClient) Authenticate(userName, oauthToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", userName, oauthToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockITwitchClientMockRecorder) Authenticate(userName, oauthToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockITwitchClient)(nil).Authenticate), userName, oauthToken)
}

// ConnectToIrcServer mocks base method.
func (m *MockITwitchClient) ConnectToIrcServer() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectToIrcServer")
	ret0, _ := ret[0].(error)
	return ret0
}

// ConnectToIrcServer indicates an expected call of ConnectToIrcServer.
func (mr *MockITwitchClientMockRecorder) ConnectToIrcServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectToIrcServer", reflect.TypeOf((*MockITwitchClient)(nil).ConnectToIrcServer))
}

// Disconnect mocks base method.
func (m *MockITwitchClient) Disconnect() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Disconnect")
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockITwitchClientMockRecorder) Disconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockITwitchClient)(nil).Disconnect))
}

// JoinChannel mocks base method.
func (m *MockITwitchClient) JoinChannel(channel string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JoinChannel", channel)
	ret0, _ := ret[0].(error)
	return ret0
}

// JoinChannel indicates an expected call of JoinChannel.
func (mr *MockITwitchClientMockRecorder) JoinChannel(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinChannel", reflect.TypeOf((*MockITwitchClient)(nil).JoinChannel), channel)
}

// LeaveChannel mocks base method.
func (m *MockITwitchClient) LeaveChannel(channel string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LeaveChannel", channel)
	ret0, _ := ret[0].(error)
	return ret0
}

// LeaveChannel indicates an expected call of LeaveChannel.
func (mr *MockITwitchClientMockRecorder) LeaveChannel(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LeaveChannel", reflect.TypeOf((*MockITwitchClient)(nil).LeaveChannel), channel)
}

// ReadLine mocks base method.
func (m *MockITwitchClient) ReadLine() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadLine")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadLine indicates an expected call of ReadLine.
func (mr *MockITwitchClientMockRecorder) ReadLine() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadLine", reflect.TypeOf((*MockITwitchClient)(nil).ReadLine))
}

// SendPong mocks base method.
func (m *MockITwitchClient) SendPong() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendPong")
	ret0, _ := ret[0].(error)
	return ret0
}

// SendPong indicates an expected call of SendPong.
func (mr *MockITwitchClientMockRecorder) SendPong() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendPong", reflect.TypeOf((*MockITwitchClient)(nil).SendPong))
}

// WriteMessage mocks base method.
func (m *MockITwitchClient) WriteMessage(message, channel, messageType, user string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteMessage", message, channel, messageType, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WriteMessage indicates an expected call of WriteMessage.
func (mr *MockITwitchClientMockRecorder) WriteMessage(message, channel, messageType, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMessage", reflect.TypeOf((*MockITwitchClient)(nil).WriteMessage), message, channel, messageType, user)
}
