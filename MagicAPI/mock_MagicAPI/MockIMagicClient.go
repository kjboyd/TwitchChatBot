// Code generated by MockGen. DO NOT EDIT.
// Source: .\MagicAPI\MagicClient.go

// Package mock_MagicAPI is a generated GoMock package.
package mock_MagicAPI

import (
	MagicAPI "TwitchChatBot/MagicAPI"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIMagicClient is a mock of IMagicClient interface.
type MockIMagicClient struct {
	ctrl     *gomock.Controller
	recorder *MockIMagicClientMockRecorder
}

// MockIMagicClientMockRecorder is the mock recorder for MockIMagicClient.
type MockIMagicClientMockRecorder struct {
	mock *MockIMagicClient
}

// NewMockIMagicClient creates a new mock instance.
func NewMockIMagicClient(ctrl *gomock.Controller) *MockIMagicClient {
	mock := &MockIMagicClient{ctrl: ctrl}
	mock.recorder = &MockIMagicClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIMagicClient) EXPECT() *MockIMagicClientMockRecorder {
	return m.recorder
}

// LookupCardInformation mocks base method.
func (m *MockIMagicClient) LookupCardInformation(cardNameOrId string) (*MagicAPI.MagicCard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupCardInformation", cardNameOrId)
	ret0, _ := ret[0].(*MagicAPI.MagicCard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupCardInformation indicates an expected call of LookupCardInformation.
func (mr *MockIMagicClientMockRecorder) LookupCardInformation(cardNameOrId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupCardInformation", reflect.TypeOf((*MockIMagicClient)(nil).LookupCardInformation), cardNameOrId)
}