// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rijuCB/RAFT2/node (interfaces: Inode)

// Package mock_node is a generated GoMock package.
package mock_node

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInode is a mock of Inode interface.
type MockInode struct {
	ctrl     *gomock.Controller
	recorder *MockInodeMockRecorder
}

// MockInodeMockRecorder is the mock recorder for MockInode.
type MockInodeMockRecorder struct {
	mock *MockInode
}

// NewMockInode creates a new mock instance.
func NewMockInode(ctrl *gomock.Controller) *MockInode {
	mock := &MockInode{ctrl: ctrl}
	mock.recorder = &MockInodeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInode) EXPECT() *MockInodeMockRecorder {
	return m.recorder
}

// candidateAction mocks base method.
func (m *MockInode) candidateAction() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "candidateAction")
}

// candidateAction indicates an expected call of candidateAction.
func (mr *MockInodeMockRecorder) candidateAction() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "candidateAction", reflect.TypeOf((*MockInode)(nil).candidateAction))
}

// followerAction mocks base method.
func (m *MockInode) followerAction() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "followerAction")
}

// followerAction indicates an expected call of followerAction.
func (mr *MockInodeMockRecorder) followerAction() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "followerAction", reflect.TypeOf((*MockInode)(nil).followerAction))
}

// leaderAction mocks base method.
func (m *MockInode) leaderAction() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "leaderAction")
}

// leaderAction indicates an expected call of leaderAction.
func (mr *MockInodeMockRecorder) leaderAction() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "leaderAction", reflect.TypeOf((*MockInode)(nil).leaderAction))
}

// performRankAction mocks base method.
func (m *MockInode) performRankAction() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "performRankAction")
}

// performRankAction indicates an expected call of performRankAction.
func (mr *MockInodeMockRecorder) performRankAction() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "performRankAction", reflect.TypeOf((*MockInode)(nil).performRankAction))
}
