// Automatically generated by MockGen. DO NOT EDIT!
// Source: client/autoscaling/client.go

package autoscaling

import (
	autoscaling "github.com/aws/aws-sdk-go/service/autoscaling"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *_MockClientRecorder
}

// Recorder for MockClient (not exported)
type _MockClientRecorder struct {
	mock *MockClient
}

func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &_MockClientRecorder{mock}
	return mock
}

func (_m *MockClient) EXPECT() *_MockClientRecorder {
	return _m.recorder
}

func (_m *MockClient) DescribeAutoScalingGroups(groups []string) (map[string]*autoscaling.Group, error) {
	ret := _m.ctrl.Call(_m, "DescribeAutoScalingGroups", groups)
	ret0, _ := ret[0].(map[string]*autoscaling.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) DescribeAutoScalingGroups(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeAutoScalingGroups", arg0)
}

func (_m *MockClient) DescribeLoadBalancerState(group string) (map[string]*autoscaling.LoadBalancerState, error) {
	ret := _m.ctrl.Call(_m, "DescribeLoadBalancerState", group)
	ret0, _ := ret[0].(map[string]*autoscaling.LoadBalancerState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) DescribeLoadBalancerState(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeLoadBalancerState", arg0)
}

func (_m *MockClient) AttachLoadBalancers(group string, lb []string) error {
	ret := _m.ctrl.Call(_m, "AttachLoadBalancers", group, lb)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) AttachLoadBalancers(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AttachLoadBalancers", arg0, arg1)
}

func (_m *MockClient) DetachLoadBalancers(group string, lb []string) error {
	ret := _m.ctrl.Call(_m, "DetachLoadBalancers", group, lb)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) DetachLoadBalancers(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DetachLoadBalancers", arg0, arg1)
}

func (_m *MockClient) AttachLoadBalancerTargetGroups(group string, targetGroupARNs []*string) error {
	ret := _m.ctrl.Call(_m, "AttachLoadBalancerTargetGroups", group, targetGroupARNs)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) AttachLoadBalancerTargetGroups(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AttachLoadBalancerTargetGroups", arg0, arg1)
}

func (_m *MockClient) DetachLoadBalancerTargetGroups(group string, targetGroupARNs []*string) error {
	ret := _m.ctrl.Call(_m, "DetachLoadBalancerTargetGroups", group, targetGroupARNs)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockClientRecorder) DetachLoadBalancerTargetGroups(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DetachLoadBalancerTargetGroups", arg0, arg1)
}

func (_m *MockClient) DescribeLoadBalancerTargetGroups(group string) ([]*autoscaling.LoadBalancerTargetGroupState, error) {
	ret := _m.ctrl.Call(_m, "DescribeLoadBalancerTargetGroups", group)
	ret0, _ := ret[0].([]*autoscaling.LoadBalancerTargetGroupState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockClientRecorder) DescribeLoadBalancerTargetGroups(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DescribeLoadBalancerTargetGroups", arg0)
}
