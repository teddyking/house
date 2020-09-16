// Code generated by counterfeiter. DO NOT EDIT.
package controllersfakes

import (
	"context"
	"sync"

	"github.com/teddyking/house/controllers"
	"github.com/teddyking/house/pkg/types"
)

type FakeHouseRepo struct {
	UpsertStub        func(context.Context, types.House) error
	upsertMutex       sync.RWMutex
	upsertArgsForCall []struct {
		arg1 context.Context
		arg2 types.House
	}
	upsertReturns struct {
		result1 error
	}
	upsertReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHouseRepo) Upsert(arg1 context.Context, arg2 types.House) error {
	fake.upsertMutex.Lock()
	ret, specificReturn := fake.upsertReturnsOnCall[len(fake.upsertArgsForCall)]
	fake.upsertArgsForCall = append(fake.upsertArgsForCall, struct {
		arg1 context.Context
		arg2 types.House
	}{arg1, arg2})
	fake.recordInvocation("Upsert", []interface{}{arg1, arg2})
	fake.upsertMutex.Unlock()
	if fake.UpsertStub != nil {
		return fake.UpsertStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.upsertReturns
	return fakeReturns.result1
}

func (fake *FakeHouseRepo) UpsertCallCount() int {
	fake.upsertMutex.RLock()
	defer fake.upsertMutex.RUnlock()
	return len(fake.upsertArgsForCall)
}

func (fake *FakeHouseRepo) UpsertCalls(stub func(context.Context, types.House) error) {
	fake.upsertMutex.Lock()
	defer fake.upsertMutex.Unlock()
	fake.UpsertStub = stub
}

func (fake *FakeHouseRepo) UpsertArgsForCall(i int) (context.Context, types.House) {
	fake.upsertMutex.RLock()
	defer fake.upsertMutex.RUnlock()
	argsForCall := fake.upsertArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeHouseRepo) UpsertReturns(result1 error) {
	fake.upsertMutex.Lock()
	defer fake.upsertMutex.Unlock()
	fake.UpsertStub = nil
	fake.upsertReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeHouseRepo) UpsertReturnsOnCall(i int, result1 error) {
	fake.upsertMutex.Lock()
	defer fake.upsertMutex.Unlock()
	fake.UpsertStub = nil
	if fake.upsertReturnsOnCall == nil {
		fake.upsertReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.upsertReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeHouseRepo) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.upsertMutex.RLock()
	defer fake.upsertMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeHouseRepo) recordInvocation(key string, args []interface{}) {
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

var _ controllers.HouseRepo = new(FakeHouseRepo)
