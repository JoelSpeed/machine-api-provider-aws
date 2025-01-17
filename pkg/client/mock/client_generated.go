// Code generated by MockGen. DO NOT EDIT.
// Source: ./client.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	ec2 "github.com/aws/aws-sdk-go/service/ec2"
	elb "github.com/aws/aws-sdk-go/service/elb"
	elbv2 "github.com/aws/aws-sdk-go/service/elbv2"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CreateTags mocks base method.
func (m *MockClient) CreateTags(arg0 *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTags", arg0)
	ret0, _ := ret[0].(*ec2.CreateTagsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTags indicates an expected call of CreateTags.
func (mr *MockClientMockRecorder) CreateTags(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTags", reflect.TypeOf((*MockClient)(nil).CreateTags), arg0)
}

// DescribeAvailabilityZones mocks base method.
func (m *MockClient) DescribeAvailabilityZones(arg0 *ec2.DescribeAvailabilityZonesInput) (*ec2.DescribeAvailabilityZonesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeAvailabilityZones", arg0)
	ret0, _ := ret[0].(*ec2.DescribeAvailabilityZonesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeAvailabilityZones indicates an expected call of DescribeAvailabilityZones.
func (mr *MockClientMockRecorder) DescribeAvailabilityZones(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeAvailabilityZones", reflect.TypeOf((*MockClient)(nil).DescribeAvailabilityZones), arg0)
}

// DescribeDHCPOptions mocks base method.
func (m *MockClient) DescribeDHCPOptions(input *ec2.DescribeDhcpOptionsInput) (*ec2.DescribeDhcpOptionsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeDHCPOptions", input)
	ret0, _ := ret[0].(*ec2.DescribeDhcpOptionsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeDHCPOptions indicates an expected call of DescribeDHCPOptions.
func (mr *MockClientMockRecorder) DescribeDHCPOptions(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeDHCPOptions", reflect.TypeOf((*MockClient)(nil).DescribeDHCPOptions), input)
}

// DescribeImages mocks base method.
func (m *MockClient) DescribeImages(arg0 *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeImages", arg0)
	ret0, _ := ret[0].(*ec2.DescribeImagesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeImages indicates an expected call of DescribeImages.
func (mr *MockClientMockRecorder) DescribeImages(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeImages", reflect.TypeOf((*MockClient)(nil).DescribeImages), arg0)
}

// DescribeInstances mocks base method.
func (m *MockClient) DescribeInstances(arg0 *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeInstances", arg0)
	ret0, _ := ret[0].(*ec2.DescribeInstancesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeInstances indicates an expected call of DescribeInstances.
func (mr *MockClientMockRecorder) DescribeInstances(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeInstances", reflect.TypeOf((*MockClient)(nil).DescribeInstances), arg0)
}

// DescribeSecurityGroups mocks base method.
func (m *MockClient) DescribeSecurityGroups(arg0 *ec2.DescribeSecurityGroupsInput) (*ec2.DescribeSecurityGroupsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeSecurityGroups", arg0)
	ret0, _ := ret[0].(*ec2.DescribeSecurityGroupsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeSecurityGroups indicates an expected call of DescribeSecurityGroups.
func (mr *MockClientMockRecorder) DescribeSecurityGroups(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeSecurityGroups", reflect.TypeOf((*MockClient)(nil).DescribeSecurityGroups), arg0)
}

// DescribeSubnets mocks base method.
func (m *MockClient) DescribeSubnets(arg0 *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeSubnets", arg0)
	ret0, _ := ret[0].(*ec2.DescribeSubnetsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeSubnets indicates an expected call of DescribeSubnets.
func (mr *MockClientMockRecorder) DescribeSubnets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeSubnets", reflect.TypeOf((*MockClient)(nil).DescribeSubnets), arg0)
}

// DescribeVolumes mocks base method.
func (m *MockClient) DescribeVolumes(arg0 *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeVolumes", arg0)
	ret0, _ := ret[0].(*ec2.DescribeVolumesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeVolumes indicates an expected call of DescribeVolumes.
func (mr *MockClientMockRecorder) DescribeVolumes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeVolumes", reflect.TypeOf((*MockClient)(nil).DescribeVolumes), arg0)
}

// DescribeVpcs mocks base method.
func (m *MockClient) DescribeVpcs(arg0 *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeVpcs", arg0)
	ret0, _ := ret[0].(*ec2.DescribeVpcsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeVpcs indicates an expected call of DescribeVpcs.
func (mr *MockClientMockRecorder) DescribeVpcs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeVpcs", reflect.TypeOf((*MockClient)(nil).DescribeVpcs), arg0)
}

// ELBv2DeregisterTargets mocks base method.
func (m *MockClient) ELBv2DeregisterTargets(arg0 *elbv2.DeregisterTargetsInput) (*elbv2.DeregisterTargetsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ELBv2DeregisterTargets", arg0)
	ret0, _ := ret[0].(*elbv2.DeregisterTargetsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ELBv2DeregisterTargets indicates an expected call of ELBv2DeregisterTargets.
func (mr *MockClientMockRecorder) ELBv2DeregisterTargets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ELBv2DeregisterTargets", reflect.TypeOf((*MockClient)(nil).ELBv2DeregisterTargets), arg0)
}

// ELBv2DescribeLoadBalancers mocks base method.
func (m *MockClient) ELBv2DescribeLoadBalancers(arg0 *elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ELBv2DescribeLoadBalancers", arg0)
	ret0, _ := ret[0].(*elbv2.DescribeLoadBalancersOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ELBv2DescribeLoadBalancers indicates an expected call of ELBv2DescribeLoadBalancers.
func (mr *MockClientMockRecorder) ELBv2DescribeLoadBalancers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ELBv2DescribeLoadBalancers", reflect.TypeOf((*MockClient)(nil).ELBv2DescribeLoadBalancers), arg0)
}

// ELBv2DescribeTargetGroups mocks base method.
func (m *MockClient) ELBv2DescribeTargetGroups(arg0 *elbv2.DescribeTargetGroupsInput) (*elbv2.DescribeTargetGroupsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ELBv2DescribeTargetGroups", arg0)
	ret0, _ := ret[0].(*elbv2.DescribeTargetGroupsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ELBv2DescribeTargetGroups indicates an expected call of ELBv2DescribeTargetGroups.
func (mr *MockClientMockRecorder) ELBv2DescribeTargetGroups(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ELBv2DescribeTargetGroups", reflect.TypeOf((*MockClient)(nil).ELBv2DescribeTargetGroups), arg0)
}

// ELBv2DescribeTargetHealth mocks base method.
func (m *MockClient) ELBv2DescribeTargetHealth(arg0 *elbv2.DescribeTargetHealthInput) (*elbv2.DescribeTargetHealthOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ELBv2DescribeTargetHealth", arg0)
	ret0, _ := ret[0].(*elbv2.DescribeTargetHealthOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ELBv2DescribeTargetHealth indicates an expected call of ELBv2DescribeTargetHealth.
func (mr *MockClientMockRecorder) ELBv2DescribeTargetHealth(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ELBv2DescribeTargetHealth", reflect.TypeOf((*MockClient)(nil).ELBv2DescribeTargetHealth), arg0)
}

// ELBv2RegisterTargets mocks base method.
func (m *MockClient) ELBv2RegisterTargets(arg0 *elbv2.RegisterTargetsInput) (*elbv2.RegisterTargetsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ELBv2RegisterTargets", arg0)
	ret0, _ := ret[0].(*elbv2.RegisterTargetsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ELBv2RegisterTargets indicates an expected call of ELBv2RegisterTargets.
func (mr *MockClientMockRecorder) ELBv2RegisterTargets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ELBv2RegisterTargets", reflect.TypeOf((*MockClient)(nil).ELBv2RegisterTargets), arg0)
}

// RegisterInstancesWithLoadBalancer mocks base method.
func (m *MockClient) RegisterInstancesWithLoadBalancer(arg0 *elb.RegisterInstancesWithLoadBalancerInput) (*elb.RegisterInstancesWithLoadBalancerOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterInstancesWithLoadBalancer", arg0)
	ret0, _ := ret[0].(*elb.RegisterInstancesWithLoadBalancerOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterInstancesWithLoadBalancer indicates an expected call of RegisterInstancesWithLoadBalancer.
func (mr *MockClientMockRecorder) RegisterInstancesWithLoadBalancer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterInstancesWithLoadBalancer", reflect.TypeOf((*MockClient)(nil).RegisterInstancesWithLoadBalancer), arg0)
}

// RunInstances mocks base method.
func (m *MockClient) RunInstances(arg0 *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInstances", arg0)
	ret0, _ := ret[0].(*ec2.Reservation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunInstances indicates an expected call of RunInstances.
func (mr *MockClientMockRecorder) RunInstances(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInstances", reflect.TypeOf((*MockClient)(nil).RunInstances), arg0)
}

// TerminateInstances mocks base method.
func (m *MockClient) TerminateInstances(arg0 *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TerminateInstances", arg0)
	ret0, _ := ret[0].(*ec2.TerminateInstancesOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TerminateInstances indicates an expected call of TerminateInstances.
func (mr *MockClientMockRecorder) TerminateInstances(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TerminateInstances", reflect.TypeOf((*MockClient)(nil).TerminateInstances), arg0)
}
