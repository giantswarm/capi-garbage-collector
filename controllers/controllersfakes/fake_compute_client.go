// Code generated by counterfeiter. DO NOT EDIT.
package controllersfakes

import (
	"context"
	"sync"

	"github.com/giantswarm/capi-garbage-collector/controllers"
	"sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
)

type FakeComputeClient struct {
	DeleteRoutesStub        func(context.Context, *v1beta1.GCPCluster) error
	deleteRoutesMutex       sync.RWMutex
	deleteRoutesArgsForCall []struct {
		arg1 context.Context
		arg2 *v1beta1.GCPCluster
	}
	deleteRoutesReturns struct {
		result1 error
	}
	deleteRoutesReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeComputeClient) DeleteRoutes(arg1 context.Context, arg2 *v1beta1.GCPCluster) error {
	fake.deleteRoutesMutex.Lock()
	ret, specificReturn := fake.deleteRoutesReturnsOnCall[len(fake.deleteRoutesArgsForCall)]
	fake.deleteRoutesArgsForCall = append(fake.deleteRoutesArgsForCall, struct {
		arg1 context.Context
		arg2 *v1beta1.GCPCluster
	}{arg1, arg2})
	stub := fake.DeleteRoutesStub
	fakeReturns := fake.deleteRoutesReturns
	fake.recordInvocation("DeleteRoutes", []interface{}{arg1, arg2})
	fake.deleteRoutesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeComputeClient) DeleteRoutesCallCount() int {
	fake.deleteRoutesMutex.RLock()
	defer fake.deleteRoutesMutex.RUnlock()
	return len(fake.deleteRoutesArgsForCall)
}

func (fake *FakeComputeClient) DeleteRoutesCalls(stub func(context.Context, *v1beta1.GCPCluster) error) {
	fake.deleteRoutesMutex.Lock()
	defer fake.deleteRoutesMutex.Unlock()
	fake.DeleteRoutesStub = stub
}

func (fake *FakeComputeClient) DeleteRoutesArgsForCall(i int) (context.Context, *v1beta1.GCPCluster) {
	fake.deleteRoutesMutex.RLock()
	defer fake.deleteRoutesMutex.RUnlock()
	argsForCall := fake.deleteRoutesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeComputeClient) DeleteRoutesReturns(result1 error) {
	fake.deleteRoutesMutex.Lock()
	defer fake.deleteRoutesMutex.Unlock()
	fake.DeleteRoutesStub = nil
	fake.deleteRoutesReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeComputeClient) DeleteRoutesReturnsOnCall(i int, result1 error) {
	fake.deleteRoutesMutex.Lock()
	defer fake.deleteRoutesMutex.Unlock()
	fake.DeleteRoutesStub = nil
	if fake.deleteRoutesReturnsOnCall == nil {
		fake.deleteRoutesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteRoutesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeComputeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteRoutesMutex.RLock()
	defer fake.deleteRoutesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeComputeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ controllers.ComputeClient = new(FakeComputeClient)